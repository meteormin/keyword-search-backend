package utils

import (
	"strings"
	"time"
)

type JsonTime time.Time

func (jt *JsonTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value) //parse time
	if err != nil {
		return err
	}

	*jt = JsonTime(t) //set result using the pointer
	return nil
}

func (jt *JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(*jt).Format(DefaultDateLayout) + `"`), nil
}
