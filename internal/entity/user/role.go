package user

import (
	"encoding/json"
	"errors"
)

type Role string

const (
	EmptyRole   Role = ""
	SuperUser   Role = "SUPER_USER"
	GeneralUser Role = "GENERAL_USER"
)

func (r *Role) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	default:
		return errors.New("invalid user role")
	case string(EmptyRole):
		*r = EmptyRole
	case string(SuperUser):
		*r = SuperUser
	case string(GeneralUser):
		*r = GeneralUser
	}

	return nil
}

func (r Role) MarshalJSON() ([]byte, error) {
	var s string
	switch r {
	default:
		return nil, errors.New("invalid user role")
	case EmptyRole:
		s = string(EmptyRole)
	case SuperUser:
		s = string(SuperUser)
	case GeneralUser:
		s = string(GeneralUser)
	}

	return json.Marshal(s)
}

func (r *Role) UnmarshalText(text []byte) error {
	switch string(text) {
	default:
		return errors.New("invalid user role")
	case string(EmptyRole):
		*r = EmptyRole
	case string(SuperUser):
		*r = SuperUser
	case string(GeneralUser):
		*r = GeneralUser
	}

	return nil
}

func (r Role) MarshalText() ([]byte, error) {
	var s string
	switch r {
	default:
		return nil, errors.New("invalid user role")
	case EmptyRole:
		s = string(EmptyRole)
	case SuperUser:
		s = string(SuperUser)
	case GeneralUser:
		s = string(GeneralUser)
	}

	return []byte(s), nil
}
