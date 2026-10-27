package backend

import (
	"sync"
)

var defaultPipedInstances = []string{
	"https://pipedapi.kavin.rocks",
	"https://api.piped.privacydev.net",
	"https://piped-api.lunar.icu",
	"https://piped-api.garudalinux.org",
}

var defaultInvidiousInstances = []string{
	"https://inv.tux.pizza",
	"https://invidious.nerdvpn.de",
	"https://vid.puffyan.us",
	"https://invidious.flokinet.to",
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
	if len(p) == 0 {
		p = append([]string{}, defaultPipedInstances...)
	}
	inv := customInvidious
	if len(inv) == 0 {
		inv = append([]string{}, defaultInvidiousInstances...)
	}
	return &InstancePool{
		piped:     p,
		invidious: inv,
	}
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
