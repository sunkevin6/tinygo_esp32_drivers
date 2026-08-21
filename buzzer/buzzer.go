package buzzer

import (
	"machine"
	"time"
)

const (
	DefaultPeriodNS = uint64(5_000_000) // 5ms
	OneSecondNS     = uint64(1_000_000_000)
)

type PWM interface {
	Configure(config machine.PWMConfig) error
	Channel(pin machine.Pin) (uint8, error)
	Set(channel uint8, value uint32)
	Top() uint32
}

type Device struct {
	pwm     PWM
	pin     machine.Pin
	channel uint8
}

func New(pwm PWM, pin machine.Pin) Device {
	return Device{pwm: pwm, pin: pin}
}

type Output struct {
	pin machine.Pin
}

func NewOutput(pin machine.Pin) Output {
	return Output{pin: pin}
}

func (output Output) Configure() {
	output.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func (output Output) On() {
	output.pin.High()
}

func (output Output) Off() {
	output.pin.Low()
}

func (device Device) Configure() error {
	return device.ConfigureWithPeriod(DefaultPeriodNS)
}

func (device Device) ConfigureWithPeriod(periodNS uint64) error {
	if err := device.pwm.Configure(machine.PWMConfig{Period: periodNS}); err != nil {
		return err
	}

	channel, err := device.pwm.Channel(device.pin)
	if err != nil {
		return err
	}

	device.channel = channel
	device.NoTone()
	return nil
}

func (device Device) Tone() {
	device.pwm.Set(device.channel, device.pwm.Top()/2)
}

func (device Device) NoTone() {
	device.pwm.Set(device.channel, 0)
}

func (device Device) ToneFrequency(frequencyHz uint32) {
	if frequencyHz == 0 {
		device.NoTone()
		return
	}

	if err := device.pwm.Configure(machine.PWMConfig{Period: OneSecondNS / uint64(frequencyHz)}); err != nil {
		device.NoTone()
		return
	}

	device.Tone()
}

func (device Device) Play(duration time.Duration) {
	device.Tone()
	time.Sleep(duration)
	device.NoTone()
}

func (device Device) PlayFrequency(frequencyHz uint32, duration time.Duration) {
	device.ToneFrequency(frequencyHz)
	time.Sleep(duration)
	device.NoTone()
}
