package rounddisplay

import "machine"

type Config struct {
	Width      int
	Height     int
	Rotation   int
	TouchInt   machine.Pin
	TouchLabel string
	SDCS       machine.Pin
	SDCSLabel  string
}

func DefaultConfig() Config {
	return Config{
		Width:      240,
		Height:     240,
		Rotation:   1,
		TouchInt:   machine.D7,
		TouchLabel: "D7/GPIO44",
		SDCS:       machine.GPIO21,
		SDCSLabel:  "GPIO21",
	}
}
