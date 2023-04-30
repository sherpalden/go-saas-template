package user

import (
	"encoding/json"
	"errors"
)

type AccountStatus string

const (
	Active   AccountStatus = "ACTIVE"
	InActive AccountStatus = "INACTIVE"
)

func (r *AccountStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	default:
		return errors.New("invalid user account status")
	case string(Active):
		*r = Active
	case string(InActive):
		*r = InActive
	}

	return nil
}

func (r AccountStatus) MarshalJSON() ([]byte, error) {
	var s string
	switch r {
	default:
		return nil, errors.New("invalid user account status")
	case Active:
		s = string(Active)
	case InActive:
		s = string(InActive)
	}

	return json.Marshal(s)
}

func (r *AccountStatus) UnmarshalText(text []byte) error {
	switch string(text) {
	default:
		return errors.New("invalid user account status")
	case string(Active):
		*r = Active
	case string(InActive):
		*r = InActive
	}

	return nil
}

func (r AccountStatus) MarshalText() ([]byte, error) {
	var s string
	switch r {
	default:
		return nil, errors.New("invalid user account status")
	case Active:
		s = string(Active)
	case InActive:
		s = string(InActive)
	}

	return []byte(s), nil
}
