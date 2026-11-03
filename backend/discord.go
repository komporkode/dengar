package backend

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Opcode constants for Discord IPC.
const (
	DiscordOpHandshake = 0
	DiscordOpFrame     = 1
	DiscordOpClose     = 2
	DiscordOpPing      = 3
	DiscordOpPong      = 4
)

// EncodeDiscordMessage packs a Discord IPC opcode and payload into bytes.
func EncodeDiscordMessage(opcode uint32, payload []byte) []byte {
	buf := make([]byte, 8+len(payload))
	binary.LittleEndian.PutUint32(buf[0:4], opcode)
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(payload)))
	copy(buf[8:], payload)
	return buf
}

// BuildActivityPayload creates the JSON payload for SET_ACTIVITY command.
func BuildActivityPayload(pid int, title, artist string, duration, elapsed int64, isPlaying bool) []byte {
	activity := map[string]interface{}{
		"details": title,
		"state":   artist,
		"assets": map[string]interface{}{
			"large_image": "app_logo",
			"large_text":  "Dengerin Desktop",
			"small_image": map[bool]string{true: "play", false: "pause"}[isPlaying],
			"small_text":  map[bool]string{true: "Memutar", false: "Dijeda"}[isPlaying],
		},
	}

	if isPlaying && duration > 0 {
		now := time.Now().Unix()
		startTime := now - elapsed
		endTime := startTime + duration
		activity["timestamps"] = map[string]interface{}{
			"start": startTime,
			"end":   endTime,
		}
	}

	payload := map[string]interface{}{
		"cmd": "SET_ACTIVITY",
		"args": map[string]interface{}{
			"pid":      pid,
			"activity": activity,
		},
		"nonce": fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	data, _ := json.Marshal(payload)
	return data
}

// DiscordService handles communication with local Discord desktop client.
type DiscordService struct {
	clientID string
	conn     io.ReadWriteCloser
	enabled  bool
	mu       sync.Mutex
	pid      int
}

// NewDiscordService initializes a new Discord RPC service instance.
func NewDiscordService(clientID string) *DiscordService {
	if clientID == "" {
		clientID = "1343829023485759539" // Dengerin Desktop Application ID
	}
	return &DiscordService{
		clientID: clientID,
		enabled:  true,
		pid:      os.Getpid(),
	}
}

// SetEnabled toggles whether Discord RPC status is broadcast.
func (d *DiscordService) SetEnabled(enabled bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.enabled = enabled
	if !enabled && d.conn != nil {
		_ = d.conn.Close()
		d.conn = nil
	}
}

// IsEnabled returns the current active state of Discord RPC.
func (d *DiscordService) IsEnabled() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.enabled
}

// ensureConnection connects to Discord IPC pipe if not already connected.
func (d *DiscordService) ensureConnection() error {
	if d.conn != nil {
		return nil
	}

	for i := 0; i < 10; i++ {
		pipePath := fmt.Sprintf(`\\.\pipe\discord-ipc-%d`, i)
		file, err := os.OpenFile(pipePath, os.O_RDWR, 0)
		if err == nil {
			d.conn = file
			// Send Handshake
			handshake, _ := json.Marshal(map[string]interface{}{
				"v":         1,
				"client_id": d.clientID,
			})
			msg := EncodeDiscordMessage(DiscordOpHandshake, handshake)
			if _, writeErr := d.conn.Write(msg); writeErr != nil {
				_ = d.conn.Close()
				d.conn = nil
				return writeErr
			}

			// Read handshake response non-blocking (discard 1 message)
			go func(c io.Reader) {
				header := make([]byte, 8)
				if _, err := io.ReadFull(c, header); err == nil {
					lenBytes := binary.LittleEndian.Uint32(header[4:8])
					buf := make([]byte, lenBytes)
					_, _ = io.ReadFull(c, buf)
				}
			}(file)

			return nil
		}
	}
	return fmt.Errorf("discord client not running")
}

// UpdateActivity sends the currently playing track to Discord.
func (d *DiscordService) UpdateActivity(title, artist string, duration, elapsed int64, isPlaying bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.enabled {
		return nil
	}

	if err := d.ensureConnection(); err != nil {
		return err
	}

	payload := BuildActivityPayload(d.pid, title, artist, duration, elapsed, isPlaying)
	msg := EncodeDiscordMessage(DiscordOpFrame, payload)

	if _, err := d.conn.Write(msg); err != nil {
		_ = d.conn.Close()
		d.conn = nil
		return err
	}
	return nil
}

// ClearActivity clears the activity on Discord.
func (d *DiscordService) ClearActivity() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.conn == nil {
		return nil
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"cmd": "SET_ACTIVITY",
		"args": map[string]interface{}{
			"pid":      d.pid,
			"activity": nil,
		},
		"nonce": fmt.Sprintf("%d", time.Now().UnixNano()),
	})
	msg := EncodeDiscordMessage(DiscordOpFrame, payload)
	_, err := d.conn.Write(msg)
	return err
}

// Close closes the Discord pipe connection.
func (d *DiscordService) Close() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.conn != nil {
		_ = d.conn.Close()
		d.conn = nil
	}
}
