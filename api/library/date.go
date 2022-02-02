package library

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

//CustomDate struct is used to parse custom time format(YYYY-MM-DD) in json document
type CustomDate struct {
	Format string
	time.Time
}

//UnmarshalJSON CustomDate method
func (Date *CustomDate) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	Date.Format = "2006-01-02"
	t, err := time.Parse(Date.Format, s)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("You must provide a date with %s format", Date.Format)
	}

	Date.Time = t
	return nil
}

// MarshalJSON CustomDate method
func (Date CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(Date.Time.Format(Date.Format))
}
