package pdmmic

import "machine"

type Config struct {
	ClockPin machine.Pin
	DataPin  machine.Pin
}

func (config Config) PDMConfig() machine.PDMConfig {
	return machine.PDMConfig{
		Stereo: false,
		DIN:    config.DataPin,
		CLK:    config.ClockPin,
	}
}
