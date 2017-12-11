package fortnox

import (
	"encoding/json"
	"fmt"
)

// Floatish type to allow unmarshalling from either string or float
type Floatish float64

func unmarshalIsh(data []byte, receiver interface{}) error {
	if len(data) == 0 {
		return nil
	}
	if data[0] == '"' {
		data = data[1:]
		data = data[:len(data)-1]
	}

	if len(data) < 1 {
		return nil
	}
	return json.Unmarshal(data, receiver)
}

// UnmarshalJSON to allow unmarshalling from either string or float
func (f *Floatish) UnmarshalJSON(data []byte) error {
	var newF float64
	err := unmarshalIsh(data, &newF)
	*f = Floatish(newF)
	return err
}

// Intish type to allow unmarshalling from either string or int
type Intish int

// UnmarshalJSON to allow unmarshalling from either string or int
func (f *Intish) UnmarshalJSON(data []byte) error {
	var newI int
	err := unmarshalIsh(data, &newI)
	*f = Intish(newI)
	return err
}

// Date simple fortnox date holder
type Date struct {
	Year  int
	Month int
	Date  int
}

// String representation of fnox date
func (d *Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Date)
}

// MarshalJSON marshals date to json
func (d *Date) MarshalJSON() ([]byte, error) {
	// sure about this??
	if d.Year == 0 || d.Month == 0 || d.Date == 0 {
		return nil, nil
	}
	return []byte(d.String()), nil
}

// UnmarshalJSON of fnox date
func (d *Date) UnmarshalJSON(data []byte) error {

	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if len(v) != 10 {
		return nil
	}

	if _, err := fmt.Sscanf(v, "%d-%d-%d", &d.Year, &d.Month, &d.Date); err != nil {
		return err
	}

	return nil
}
