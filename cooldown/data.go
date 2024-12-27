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

func marshalMappedCooldown[T comparable](m MappedCoolDown[T], marshaler gophig.Marshaler) ([]byte, error) {
	d := map[T]coolDownData{}
	for k, cd := range m {
		d[k] = coolDownData{
			Expiration:       cd.expiration.Load(),
			Paused:           cd.paused.Load(),
			RemainingAtPause: cd.remainingAtPause.Load(),
		}
	}
	return marshaler.Marshal(d)
}

func unmarshalMappedCooldown[T comparable](m MappedCoolDown[T], b []byte, marshaler gophig.Marshaler) error {
	if m == nil {
		m = make(MappedCoolDown[T])
	}
	d := map[T]coolDownData{}
	err := marshaler.Unmarshal(b, &d)
	if err != nil {
		return err
	}
	for k, cd := range d {
		m[k] = &CoolDown{
			expiration:       *atomic.NewValue(cd.Expiration),
			paused:           *atomic.NewBool(cd.Paused),
			remainingAtPause: *atomic.NewValue(cd.RemainingAtPause),
		}
	}
	return nil
}
