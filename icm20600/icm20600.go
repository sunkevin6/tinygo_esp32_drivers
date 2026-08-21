package icm20600

import (
	"machine"
	"time"
)

const (
	AddressLow  = 0x68
	AddressHigh = 0x69
)

const (
	regAccelXOutH  = 0x3B
	regPwrMgmt1    = 0x6B
	regAccelConfig = 0x1C
	regWhoAmI      = 0x75
)

const accelRange16G = 0x18

type Device struct {
	bus     *machine.I2C
	address uint16
}

func New(bus *machine.I2C) Device {
	return Device{bus: bus, address: AddressHigh}
}

func (device *Device) Configure() bool {
	if !device.probe(device.address) {
		device.address = AddressLow
		if !device.probe(device.address) {
			return false
		}
	}

	if !device.writeRegister(regPwrMgmt1, 0x01) {
		return false
	}
	time.Sleep(10 * time.Millisecond)

	return device.writeRegister(regAccelConfig, accelRange16G)
}

func (device Device) Address() uint16 {
	return device.address
}

func (device *Device) probe(address uint16) bool {
	device.address = address
	var data [1]byte
	return device.readRegister(regWhoAmI, data[:])
}

func (device Device) ReadAccelerationMG() (int32, int32, int32) {
	var data [6]byte
	if !device.readRegister(regAccelXOutH, data[:]) {
		return 0, 0, 0
	}

	x := int16(uint16(data[0])<<8 | uint16(data[1]))
	y := int16(uint16(data[2])<<8 | uint16(data[3]))
	z := int16(uint16(data[4])<<8 | uint16(data[5]))

	return Raw16GToMG(x), Raw16GToMG(y), Raw16GToMG(z)
}

func Raw16GToMG(value int16) int32 {
	return int32(value) * 16000 / 32768
}

func (device Device) readRegister(register uint8, data []byte) bool {
	if err := device.bus.Tx(device.address, []byte{register}, data); err != nil {
		return false
	}
	return true
}

func (device Device) writeRegister(register uint8, value uint8) bool {
	if err := device.bus.Tx(device.address, []byte{register, value}, nil); err != nil {
		return false
	}
	return true
}
