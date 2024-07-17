package cooldown

import (
	"github.com/df-mc/atomic"
	"github.com/restartfu/gophig"
	"time"
)

// coolDownData represents the data of a CoolDown that is marshaled and unmarshaled.
type coolDownData struct {
	Expiration       time.Time
	Paused           bool
	RemainingAtPause time.Duration
}

func marshalCooldown(c *CoolDown, marshaler gophig.Marshaler) ([]byte, error) {
	d := coolDownData{
		Expiration:       c.expiration.Load(),
		Paused:           c.paused.Load(),
		RemainingAtPause: c.remainingAtPause.Load(),
	}
	return marshaler.Marshal(d)
}

func unmarshalCooldown(c *CoolDown, b []byte, marshaler gophig.Marshaler) error {
	d := coolDownData{}
	err := marshaler.Unmarshal(b, &d)
	c.expiration = *atomic.NewValue(d.Expiration)
	c.paused.Store(d.Paused)
	c.remainingAtPause.Store(d.RemainingAtPause)
	return err
}
