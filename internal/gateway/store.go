package gateway

import (
	"sync"
	"time"
)

// LinkState is the last known state of the vehicle's cellular link.
type LinkState string

const (
	LinkOnline   LinkState = "online"
	LinkAsleep   LinkState = "asleep"
	LinkOffline  LinkState = "offline"
	LinkUpdating LinkState = "updating"
)

// Vehicle is the connectivity record the gateway keeps for one vehicle.
type Vehicle struct {
	VIN          string
	Link         LinkState
	SignalBars   int
	LastSeen     time.Time
	FirmwareTrim string
}

// Store looks up connectivity records by VIN.
type Store interface {
	Get(vin string) (Vehicle, bool)
}

// MemoryStore is an in-memory Store used by the demo and the tests.
type MemoryStore struct {
	mu       sync.RWMutex
	vehicles map[string]Vehicle
}

func NewMemoryStore(vehicles []Vehicle) *MemoryStore {
	m := &MemoryStore{vehicles: make(map[string]Vehicle, len(vehicles))}
	for _, v := range vehicles {
		m.vehicles[v.VIN] = v
	}
	return m
}

func (m *MemoryStore) Get(vin string) (Vehicle, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.vehicles[vin]
	return v, ok
}

// SeedVehicles returns invented vehicles. The VINs are not real.
func SeedVehicles() []Vehicle {
	seen := time.Date(2026, 11, 1, 8, 30, 0, 0, time.UTC)
	return []Vehicle{
		{VIN: "5LDEMAA0000000001", Link: LinkOnline, SignalBars: 4, LastSeen: seen, FirmwareTrim: "2.6.1"},
		{VIN: "5LDEMAA0000000002", Link: LinkAsleep, SignalBars: 0, LastSeen: seen.Add(-6 * time.Hour), FirmwareTrim: "2.6.1"},
		{VIN: "5LDEMAA0000000003", Link: LinkUpdating, SignalBars: 3, LastSeen: seen.Add(-2 * time.Minute), FirmwareTrim: "2.5.9"},
		{VIN: "5LDEMAA0000000004", Link: LinkOffline, SignalBars: 0, LastSeen: seen.Add(-72 * time.Hour), FirmwareTrim: "2.4.0"},
	}
}
