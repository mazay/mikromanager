package db

import (
	"fmt"
	"strings"
	"time"
)

type FirmwareBuildTime time.Time

func (bt *FirmwareBuildTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		return err
	}
	*bt = FirmwareBuildTime(t)
	return nil
}

func (bt FirmwareBuildTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(bt).Format("2006-01-02 15:04:05") + `"`), nil
}

//nolint:golint,errcheck
func (bt FirmwareBuildTime) Format(f fmt.State, c rune) {
	f.Write([]byte(time.Time(bt).Format("2006-01-02 15:04:05")))
}
