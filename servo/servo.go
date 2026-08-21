package servo

import "machine"

const (
	DefaultPeriodUS       = uint32(20_000)
	DefaultMinPulseUS     = uint32(544)
	DefaultMaxPulseUS     = uint32(2400)
	DefaultMaxAngleTenths = int32(1800)
)

type Config struct {
	PeriodUS       uint32
	MinPulseUS     uint32
	MaxPulseUS     uint32
	MaxAngleTenths int32
}

type PWM interface {
	Configure(config machine.PWMConfig) error
	Change(pin machine.Pin) (uint8, error)
	Set(channel uint8, value uint32)
	Top() uint32
}

type Device struct {
	pwm     PWM
	pin     machine.Pin
	config  Config
	channel uint8
}

func New(pwm PWM, pin machine.Pin) Device {
	return Device{
		pwm:    pwm,
		pin:    pin,
		config: DefaultConfig(),
	}
}

func DefaultConfig() Config {
	return Config{
		PeriodUS:       DefaultPeriodUS,
		MinPulseUS:     DefaultMinPulseUS,
		MaxPulseUS:     DefaultMaxPulseUS,
		MaxAngleTenths: DefaultMaxAngleTenths,
	}
}

func (device *Device) Configure() error {
	return device.ConfigureWith(DefaultConfig())
}

func (device *Device) ConfigureWith(config Config) error {
	if err := device.pwm.Configure(machine.PWMConfig{Period: uint64(config.PeriodUS) * 1000}); err != nil {
		return err
	}

	channel, err := device.pwm.Change(device.pin)
	if err != nil {
		return err
	}

	device.channel = channel
	device.config = config

	return nil
}

func (device Device) WriteAngle(angleDegrees int32) {
	device.WriteAngleTenths(angleDegrees * 10)
}

func (device Device) WriteAngleTenths(angleTenths int32) {
	if angleTenths < 0 {
		angleTenths = 0
	} else if angleTenths > device.config.MaxAngleTenths {
		angleTenths = device.config.MaxAngleTenths
	}

	pulseWidthUS := device.config.MinPulseUS + uint32(angleTenths)*(device.config.MaxPulseUS-device.config.MinPulseUS)/uint32(device.config.MaxAngleTenths)
	device.pwm.Set(device.channel, pulseWidthUS*device.pwm.Top()/device.config.PeriodUS)
}
