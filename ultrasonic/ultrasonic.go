package ultrasonic

import (
	"machine"
	"time"
)

const (
	defaultTriggerPulse = 10 * time.Microsecond
	defaultEchoTimeout  = 30 * time.Millisecond
	microsecondsPerInch = 148
	microsecondsPerCM   = 58
)

type Config struct {
	TriggerPulse time.Duration
	EchoTimeout  time.Duration
}

type Sensor struct {
	pin    machine.Pin
	config Config
}

func New(pin machine.Pin) Sensor {
	return Sensor{pin: pin, config: DefaultConfig()}
}

func DefaultConfig() Config {
	return Config{
		TriggerPulse: defaultTriggerPulse,
		EchoTimeout:  defaultEchoTimeout,
	}
}

func (sensor Sensor) Configure(config Config) Sensor {
	sensor.config = config
	return sensor
}

func (sensor Sensor) Inches() uint32 {
	return sensor.PulseMicroseconds() / microsecondsPerInch
}

func (sensor Sensor) Centimeters() uint32 {
	return sensor.PulseMicroseconds() / microsecondsPerCM
}

func (sensor Sensor) PulseMicroseconds() uint32 {
	pulseWidth := sensor.EchoPulse()
	if pulseWidth <= 0 {
		return 0
	}
	return uint32(pulseWidth / time.Microsecond)
}

func (sensor Sensor) EchoPulse() time.Duration {
	sensor.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	sensor.pin.Low()
	time.Sleep(2 * time.Microsecond)
	sensor.pin.High()
	time.Sleep(sensor.config.TriggerPulse)
	sensor.pin.Low()

	sensor.pin.Configure(machine.PinConfig{Mode: machine.PinInput})
	if !sensor.waitForLevel(true, sensor.config.EchoTimeout) {
		return 0
	}

	start := time.Now()
	if !sensor.waitForLevel(false, sensor.config.EchoTimeout) {
		return 0
	}
	return time.Since(start)
}

func (sensor Sensor) waitForLevel(level bool, timeout time.Duration) bool {
	start := time.Now()
	for sensor.pin.Get() != level {
		if time.Since(start) > timeout {
			return false
		}
	}
	return true
}
