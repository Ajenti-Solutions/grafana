package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana/pkg/services/featuremgmt/strcase"
)

type SettingsSource int

const (
	DB = iota
	System
)

func (s SettingsSource) MarshalJSON() ([]byte, error) {
	switch s {
	case DB:
		return json.Marshal("database")
	case System:
		return json.Marshal("system")
	default:
		return nil, fmt.Errorf("unknown source: %d", s)
	}
}

type SSOSetting struct {
	ID        string                 `xorm:"id pk" json:"-"`
	Provider  string                 `xorm:"provider" json:"provider"`
	Settings  map[string]interface{} `xorm:"settings" json:"settings"`
	Created   time.Time              `xorm:"created" json:"-"`
	Updated   time.Time              `xorm:"updated" json:"-"`
	IsDeleted bool                   `xorm:"is_deleted" json:"-"`
	Source    SettingsSource         `xorm:"-" json:"source"`
}

// TableName returns the table name (needed for Xorm)
func (s SSOSetting) TableName() string {
	return "sso_setting"
}

// MarshalJSON implements the json.Marshaler interface and converts the s.Settings from map[string]any to map[string]any in camelCase
func (s SSOSetting) MarshalJSON() ([]byte, error) {
	type Alias SSOSetting
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(&s),
	}

	settings := make(map[string]any)
	for k, v := range aux.Settings {
		settings[strcase.ToLowerCamel(k)] = v
	}

	aux.Settings = settings
	return json.Marshal(aux)
}

// UnmarshalJSON implements the json.Unmarshaler interface and converts the settings from map[string]any camelCase to map[string]interface{} snake_case
func (s *SSOSetting) UnmarshalJSON(data []byte) error {
	type Alias SSOSetting
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	settings := make(map[string]any)
	for k, v := range aux.Settings {
		settings[strcase.ToSnake(k)] = v
	}

	s.Settings = settings
	return nil
}

type SSOSettingsResponse struct {
	Settings map[string]interface{} `json:"settings"`
	Provider string                 `json:"type"`
}
