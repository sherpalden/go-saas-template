package model

import (
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type ID uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

// scan ID from db
func (id *ID) Scan(value interface{}) error {
	parsedID, err := StringToID(value.(string))
	if err != nil {
		return err
	}
	*id = parsedID
	return nil
}

func StringToID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(id).String())
}

func (id *ID) UnmarshalJSON(bt []byte) error {
	s, err := uuid.ParseBytes(bt)
	if err != nil {
		return err
	}
	*id = ID(s)
	return nil
}
