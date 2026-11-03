package backend

import (
	"encoding/binary"
	"encoding/json"
	"testing"
)

func TestEncodeDiscordMessage(t *testing.T) {
	payload := []byte(`{"test":"hello"}`)
	msg := EncodeDiscordMessage(1, payload)

	if len(msg) != 8+len(payload) {
		t.Fatalf("Expected length %d, got %d", 8+len(payload), len(msg))
	}

	opcode := binary.LittleEndian.Uint32(msg[0:4])
	if opcode != 1 {
		t.Errorf("Expected opcode 1, got %d", opcode)
	}

	length := binary.LittleEndian.Uint32(msg[4:8])
	if length != uint32(len(payload)) {
		t.Errorf("Expected length %d, got %d", len(payload), length)
	}

	body := string(msg[8:])
	if body != string(payload) {
		t.Errorf("Expected body %s, got %s", payload, body)
	}
}

func TestBuildActivityPayload(t *testing.T) {
	payloadBytes := BuildActivityPayload(1234, "Dan", "Sheila on 7", 240, 30, true)
	if len(payloadBytes) == 0 {
		t.Fatal("Payload is empty")
	}

	var data map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &data); err != nil {
		t.Fatalf("Payload is not valid JSON: %v", err)
	}

	if data["cmd"] != "SET_ACTIVITY" {
		t.Errorf("Expected cmd SET_ACTIVITY, got %v", data["cmd"])
	}

	args, ok := data["args"].(map[string]interface{})
	if !ok {
		t.Fatalf("args is not a map: %v", data["args"])
	}

	activity, ok := args["activity"].(map[string]interface{})
	if !ok {
		t.Fatalf("activity is not a map: %v", args["activity"])
	}

	if activity["details"] != "Dan" {
		t.Errorf("Expected details 'Dan', got %v", activity["details"])
	}

	if activity["state"] != "Sheila on 7" {
		t.Errorf("Expected state 'Sheila on 7', got %v", activity["state"])
	}
}
