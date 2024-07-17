package cooldown

import (
	"github.com/df-mc/atomic"
	"time"
)

// CoolDown represents a time cooldown.
type CoolDown struct {
	expiration       atomic.Value[time.Time]
	paused           atomic.Bool
	remainingAtPause atomic.Value[time.Duration]
}

// NewCoolDown returns a new CoolDown instance.
func NewCoolDown() *CoolDown {
	return &CoolDown{}
}

// TogglePause toggles the pause state of the cooldown.
func (c *CoolDown) TogglePause() {
	currentPaused := c.paused.Load()
	c.paused.Store(!currentPaused)

	if currentPaused { // If currently paused, resume
		remaining := c.remainingAtPause.Load()
		c.expiration.Store(time.Now().Add(remaining))
		c.remainingAtPause.Store(0)
	} else { // If currently active, pause
		remaining := time.Until(c.expiration.Load())
		c.remainingAtPause.Store(remaining)
		c.expiration.Store(time.Time{}) // Clear expiration on pause
	}
}

// Paused returns true if the cooldown is paused.
func (c *CoolDown) Paused() bool {
	return c.paused.Load()
}

// Set sets the cooldown duration.
func (c *CoolDown) Set(dur time.Duration) {
	if c.paused.Load() {
		c.remainingAtPause.Store(dur)
		return
	}

	c.expiration.Store(time.Now().Add(dur))
}

// Reduce reduces the remaining cooldown duration by the specified amount.
func (c *CoolDown) Reduce(dur time.Duration) {
	if c.paused.Load() {
		remaining := c.remainingAtPause.Load()
		if remaining <= dur {
			c.Reset()
			return
		}
		c.remainingAtPause.Store(remaining - dur)
		return
	}

	exp := c.expiration.Load()
	if time.Until(exp) <= dur {
		c.Reset()
		return
	}
	c.expiration.Store(exp.Add(-dur))
}

// Active returns true if the cooldown is currently active.
func (c *CoolDown) Active() bool {
	if c.paused.Load() {
		return c.remainingAtPause.Load() > 0
	}
	return c.expiration.Load().After(time.Now())
}

// Remaining returns the remaining cooldown duration.
func (c *CoolDown) Remaining() time.Duration {
	if c.paused.Load() {
		return c.remainingAtPause.Load()
	}
	exp := c.expiration.Load()
	return time.Until(exp)
}

// Reset resets the cooldown.
func (c *CoolDown) Reset() {
	c.paused.Store(false)
	c.remainingAtPause.Store(0)
	c.expiration.Store(time.Time{}) // Clear expiration
}
