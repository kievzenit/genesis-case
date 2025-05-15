package models

type Frequency string

const (
	Hourly Frequency = "hourly"
	Daily  Frequency = "daily"
)

func (f Frequency) IsValid() bool {
	switch f {
	case Hourly, Daily:
		return true
	default:
		return false
	}
}
