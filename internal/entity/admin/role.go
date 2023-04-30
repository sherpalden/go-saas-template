package admin

import (
	"encoding/json"
	"errors"
)

type Role string

const (
	EmptyRole    Role = ""
	SuperAdmin   Role = "SUPER_ADMIN"
	GeneralAdmin Role = "GENERAL_ADMIN"
)

func (r *Role) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	default:
		return errors.New("invalid admin role")
	case string(EmptyRole):
		*r = EmptyRole
	case string(SuperAdmin):
		*r = SuperAdmin
	case string(GeneralAdmin):
		*r = GeneralAdmin
	}

	return nil
}

func (r Role) MarshalJSON() ([]byte, error) {
	var s string
	switch r {
	default:
		return nil, errors.New("invalid admin role")
	case EmptyRole:
		s = string(EmptyRole)
	case SuperAdmin:
		s = string(SuperAdmin)
	case GeneralAdmin:
		s = string(GeneralAdmin)
	}

	return json.Marshal(s)
}

func (r *Role) UnmarshalText(text []byte) error {
	switch string(text) {
	default:
		return errors.New("invalid admin role")
	case string(EmptyRole):
		*r = EmptyRole
	case string(SuperAdmin):
		*r = SuperAdmin
	case string(GeneralAdmin):
		*r = GeneralAdmin
	}

	return nil
}

func (r Role) MarshalText() ([]byte, error) {
	var s string
	switch r {
	default:
		return nil, errors.New("invalid admin role")
	case EmptyRole:
		s = string(EmptyRole)
	case SuperAdmin:
		s = string(SuperAdmin)
	case GeneralAdmin:
		s = string(GeneralAdmin)
	}

	return []byte(s), nil
}
