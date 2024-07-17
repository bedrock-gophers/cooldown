package cooldown

import "github.com/rcrowley/go-bson"

// bsonMarshaler is a Marshaler that uses the go-bson package to marshal and unmarshal data.
type bsonMarshaler struct{}

// Marshal ...
func (bsonMarshaler) Marshal(v interface{}) ([]byte, error) {
	return bson.Marshal(v)
}

// Unmarshal ...
func (bsonMarshaler) Unmarshal(data []byte, v interface{}) error {
	return bson.Unmarshal(data, v)
}

// MarshalBSON ...
func (c *CoolDown) MarshalBSON() ([]byte, error) {
	return marshalCooldown(c, bsonMarshaler{})
}

// UnmarshalBSON ...
func (c *CoolDown) UnmarshalBSON(b []byte) error {
	return unmarshalCooldown(c, b, bsonMarshaler{})
}

// MarshalBSON ...
func (m MappedCoolDown[T]) MarshalBSON() ([]byte, error) {
	return marshalMappedCooldown(m, bsonMarshaler{})
}

// UnmarshalBSON ...
func (m MappedCoolDown[T]) UnmarshalBSON(b []byte) error {
	return unmarshalMappedCooldown(m, b, bsonMarshaler{})
}
