package irremote

import (
	"machine"
	"time"
)

const (
	leaderLowMin  = 8 * time.Millisecond
	leaderLowMax  = 10 * time.Millisecond
	leaderHighMin = 3500 * time.Millisecond
	leaderHighMax = 3500 * time.Millisecond
	repeatHighMin = 1800 * time.Millisecond
	repeatHighMax = 2800 * time.Millisecond
	bitLowMin     = 350 * time.Microsecond
	bitLowMax     = 800 * time.Microsecond
	zeroHighMin   = 350 * time.Microsecond
	zeroHighMax   = 800 * time.Microsecond
	oneHighMin    = 1400 * time.Microsecond
	oneHighMax    = 1900 * time.Microsecond
)

type Receiver struct {
	pin machine.Pin
}

func New(pin machine.Pin) Receiver {
	return Receiver{pin: pin}
}

func (receiver Receiver) Configure() {
	receiver.pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
}

func (receiver Receiver) ReadNECCommand(timeout time.Duration) (uint8, bool) {
	raw, ok := receiver.readNECRaw(timeout)
	if !ok {
		return 0, false
	}

	command := uint8(raw >> 16)
	invertedCommand := uint8(raw >> 24)
	if command^invertedCommand != 0xFF {
		return 0, false
	}

	return command, true
}

func (receiver Receiver) ReadNECCode(timeout time.Duration) (uint32, bool) {
	raw, ok := receiver.readNECRaw(timeout)
	if !ok {
		return 0, false
	}

	code := reverseBytes(raw)
	if code == 0 {
		return 0, false
	}
	return code, true
}

func (receiver Receiver) readNECRaw(timeout time.Duration) (uint32, bool) {
	if !receiver.waitForLow(timeout) {
		return 0, false
	}

	leaderLow, ok := receiver.measureLevel(false, 12*time.Millisecond)
	if !ok || !inRange(leaderLow, leaderLowMin, leaderLowMax) {
		return 0, false
	}

	leaderHigh, ok := receiver.measureLevel(true, 6*time.Millisecond)
	if !ok {
		return 0, false
	}
	if inRange(leaderHigh, repeatHighMin, repeatHighMax) {
		return 0, false
	}
	if !inRange(leaderHigh, leaderHighMin, leaderHighMax) {
		return 0, false
	}

	var raw uint32
	for bitIndex := uint8(0); bitIndex < 32; bitIndex++ {
		bitLow, ok := receiver.measureLevel(false, 2*time.Millisecond)
		if !ok || !inRange(bitLow, bitLowMin, bitLowMax) {
			return 0, false
		}

		bitHigh, ok := receiver.measureLevel(true, 3*time.Millisecond)
		if !ok {
			return 0, false
		}
		if inRange(bitHigh, oneHighMin, oneHighMax) {
			raw |= 1 << bitIndex
		} else if !inRange(bitHigh, zeroHighMin, zeroHighMax) {
			return 0, false
		}
	}

	return raw, true
}

func (receiver Receiver) waitForLow(timeout time.Duration) bool {
	start := time.Now()
	for receiver.pin.Get() {
		if time.Since(start) > timeout {
			return false
		}
	}
	return true
}

func (receiver Receiver) measureLevel(level bool, timeout time.Duration) (time.Duration, bool) {
	start := time.Now()
	for receiver.pin.Get() != level {
		if time.Since(start) > timeout {
			return 0, false
		}
	}

	start = time.Now()
	for receiver.pin.Get() == level {
		if time.Since(start) > timeout {
			return 0, false
		}
	}
	return time.Since(start), true
}

func inRange(value, minimum, maximum time.Duration) bool {
	return value >= minimum && value <= maximum
}

func reverseBytes(value uint32) uint32 {
	return value<<24 | value<<8&0x00FF0000 | value>>8&0x0000FF00 | value>>24
}
