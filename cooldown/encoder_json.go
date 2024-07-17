package cooldown

import (
	"encoding/json"
)

// jsonMarshaler is a Marshaler that uses the encoding/json package to marshal and unmarshal data.
type jsonMarshaler struct{}

// Marshal ...
func (jsonMarshaler) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal ...
func (jsonMarshaler) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// MarshalJSON ...
func (c *CoolDown) MarshalJSON() ([]byte, error) {
	return marshalCooldown(c, jsonMarshaler{})
}

// UnmarshalJSON ...
func (c *CoolDown) UnmarshalJSON(b []byte) error {
	return unmarshalCooldown(c, b, jsonMarshaler{})
}
