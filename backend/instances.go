package backend

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

var defaultPipedInstances = []string{
	"https://api.piped.private.coffee",
	"https://pipedapi.leptons.xyz",
	"https://pipedapi.drgns.space",
}

var defaultInvidiousInstances = []string{
	"https://yt.chocolatemoo53.com",
}

// InstancePool manages rotating instance endpoints for failover.
type InstancePool struct {
	mu        sync.RWMutex
	piped     []string
	invidious []string
	pipedIdx  int
	invidIdx  int
}

// NewInstancePool initializes a pool with optional custom or default instances.
func NewInstancePool(customPiped, customInvidious []string) *InstancePool {
	p := customPiped
	if customPiped == nil {
		p = append([]string{}, defaultPipedInstances...)
	}
	inv := customInvidious
	if customInvidious == nil {
		inv = append([]string{}, defaultInvidiousInstances...)
	}
	pool := &InstancePool{
		piped:     p,
		invidious: inv,
	}
	if customInvidious == nil {
		pool.RefreshInvidiousInstances()
	}
	return pool
}

// RefreshInvidiousInstances asynchronously fetches verified live instances from the official registry.
func (ip *InstancePool) RefreshInvidiousInstances() {
	go func() {
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get("https://api.invidious.io/instances.json")
		if err != nil {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return
		}

		var raw [][]json.RawMessage
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			return
		}

		var newInstances []string
		for _, item := range raw {
			if len(item) < 2 {
				continue
			}
			var domain string
			if err := json.Unmarshal(item[0], &domain); err != nil {
				continue
			}
			var meta struct {
				Type string `json:"type"`
				API  bool   `json:"api"`
			}
			if err := json.Unmarshal(item[1], &meta); err != nil {
				continue
			}
			if meta.Type == "https" && meta.API {
				d := strings.ToLower(domain)
				if strings.HasSuffix(d, ".onion") || strings.HasSuffix(d, ".i2p") || strings.HasSuffix(d, ".ygg") || strings.Contains(d, "ygg") || strings.HasSuffix(d, ".local") {
					continue
				}
				newInstances = append(newInstances, "https://"+domain)
			}
		}

		if len(newInstances) > 0 {
			ip.mu.Lock()
			seen := make(map[string]bool)
			merged := make([]string, 0, len(newInstances)+len(ip.invidious))
			// Always keep existing known-good instances at the front of the pool
			for _, inst := range ip.invidious {
				if !seen[inst] {
					seen[inst] = true
					merged = append(merged, inst)
				}
			}
			for _, inst := range newInstances {
				if !seen[inst] {
					seen[inst] = true
					merged = append(merged, inst)
				}
			}
			ip.invidious = merged
			ip.mu.Unlock()
		}
	}()
}

// GetPipedInstances returns the list of Piped instances in failover order.
func (ip *InstancePool) GetPipedInstances() []string {
	ip.mu.RLock()
	defer ip.mu.RUnlock()
	res := make([]string, len(ip.piped))
	for i := range ip.piped {
		idx := (ip.pipedIdx + i) % len(ip.piped)
		res[i] = ip.piped[idx]
	}
	return res
}

// RotatePiped shifts the current primary Piped instance.
func (ip *InstancePool) RotatePiped() {
	ip.mu.Lock()
	defer ip.mu.Unlock()
	if len(ip.piped) > 0 {
		ip.pipedIdx = (ip.pipedIdx + 1) % len(ip.piped)
	}
}

// GetInvidiousInstances returns the list of Invidious instances in failover order.
func (ip *InstancePool) GetInvidiousInstances() []string {
	ip.mu.RLock()
	defer ip.mu.RUnlock()
	res := make([]string, len(ip.invidious))
	for i := range ip.invidious {
		idx := (ip.invidIdx + i) % len(ip.invidious)
		res[i] = ip.invidious[idx]
	}
	return res
}

// RotateInvidious shifts the current primary Invidious instance.
func (ip *InstancePool) RotateInvidious() {
	ip.mu.Lock()
	defer ip.mu.Unlock()
	if len(ip.invidious) > 0 {
		ip.invidIdx = (ip.invidIdx + 1) % len(ip.invidious)
	}
}
