package db

import (
	"net/url"
	"strings"
)

type DeviceBoardName string

func (bn *DeviceBoardName) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	*bn = DeviceBoardName(value)
	return nil
}

func (bn *DeviceBoardName) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strings.ReplaceAll(string(*bn), "^", "") + `"`), nil
}

func (bn *DeviceBoardName) UrlEncode() string {
	return url.QueryEscape(string(*bn))
}
