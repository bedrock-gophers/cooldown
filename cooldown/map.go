package cooldown

import (
	"time"
)

// MappedCoolDown represents a cool-down mapped to a key.
type MappedCoolDown[T comparable] map[T]*CoolDown

// NewMappedCoolDown returns a new mapped cool-down.
func NewMappedCoolDown[T comparable]() MappedCoolDown[T] {
	return make(map[T]*CoolDown)
}

// Active returns true if the cool-down is active.
func (m MappedCoolDown[T]) Active(key T) bool {
	coolDown, ok := m[key]
	return ok && coolDown.Active()
}

// Set sets the cool-down.
func (m MappedCoolDown[T]) Set(key T, d time.Duration) {
	coolDown := m.Key(key)
	coolDown.Set(d)
	m[key] = coolDown
}

// Reduce reduces the cool-down.
func (m MappedCoolDown[T]) Reduce(key T, d time.Duration) {
	coolDown := m.Key(key)
	coolDown.Reduce(d)
	m[key] = coolDown
}

// Key returns the cool-down for the key.
func (m MappedCoolDown[T]) Key(key T) *CoolDown {
	coolDown, ok := m[key]
	if !ok {
		newCD := NewCoolDown()
		m[key] = newCD
		return newCD
	}
	return coolDown
}

// Reset resets the cool-down.
func (m MappedCoolDown[T]) Reset(key T) {
	delete(m, key)
}

// Remaining returns the remaining time of the cool-down.
func (m MappedCoolDown[T]) Remaining(key T) time.Duration {
	coolDown, ok := m[key]
	if !ok {
		return 0
	}
	return coolDown.Remaining()
}

// All returns all cool-downs.
func (m MappedCoolDown[T]) All() (coolDowns []*CoolDown) {
	for _, coolDown := range m {
		coolDowns = append(coolDowns, coolDown)
	}
	return coolDowns
}

// ActiveKeys returns all active keys
func (m MappedCoolDown[T]) ActiveKeys() (keys []T) {
	for key, coolDown := range m {
		if coolDown.Active() {
			keys = append(keys, key)
		}
	}
	return keys
}
