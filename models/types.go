package models

import (
	"encoding/json"
	"math"
	"time"
)

type Resolution int

const (
	Sec15    = 15
	Minute   = 60
	Minute5  = 300
	Minute15 = 900
	Hour     = 3600
	Hour4    = 14400
	Day      = 86400
)

type FTXTime struct {
	Time time.Time
}

func (p *FTXTime) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	sec, nsec := math.Modf(f)
	p.Time = time.Unix(int64(sec), int64(nsec))
	return nil
}
