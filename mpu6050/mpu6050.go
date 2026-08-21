package mpu6050

import (
	"machine"
	"time"
)

const Address = 0x68

const (
	regXAOffsetH   = 0x06
	regYAOffsetH   = 0x08
	regZAOffsetH   = 0x0A
	regXGOffsetH   = 0x13
	regYGOffsetH   = 0x15
	regZGOffsetH   = 0x17
	regAccelConfig = 0x1C
	regAccelXOutH  = 0x3B
	regPwrMgmt1    = 0x6B
	regWhoAmI      = 0x75
)

type Device struct {
	bus     *machine.I2C
	address uint16
}

type Offsets struct {
	AccelX int32
	AccelY int32
	AccelZ int32
	GyroX  int32
	GyroY  int32
	GyroZ  int32
}

func New(bus *machine.I2C) Device {
	return Device{bus: bus, address: Address}
}

func (device Device) Initialize() {
	device.writeRegister(regPwrMgmt1, 0x01)
	time.Sleep(10 * time.Millisecond)
}

func (device Device) Connected() bool {
	var data [1]byte
	if !device.readRegister(regWhoAmI, data[:]) {
		return false
	}
	return data[0] == 0x68
}

func (device Device) ReadMotion16() (int16, int16, int16, int16, int16, int16) {
	var data [14]byte
	if !device.readRegister(regAccelXOutH, data[:]) {
		return 0, 0, 0, 0, 0, 0
	}

	ax := int16(uint16(data[0])<<8 | uint16(data[1]))
	ay := int16(uint16(data[2])<<8 | uint16(data[3]))
	az := int16(uint16(data[4])<<8 | uint16(data[5]))
	gx := int16(uint16(data[8])<<8 | uint16(data[9]))
	gy := int16(uint16(data[10])<<8 | uint16(data[11]))
	gz := int16(uint16(data[12])<<8 | uint16(data[13]))
	return ax, ay, az, gx, gy, gz
}

func (device Device) ReadAccelerationRaw() (int16, int16, int16) {
	var data [6]byte
	if !device.readRegister(regAccelXOutH, data[:]) {
		return 0, 0, 0
	}

	ax := int16(uint16(data[0])<<8 | uint16(data[1]))
	ay := int16(uint16(data[2])<<8 | uint16(data[3]))
	az := int16(uint16(data[4])<<8 | uint16(data[5]))
	return ax, ay, az
}

func (device Device) SetOffsets(offsets Offsets) {
	device.SetAccelOffset(offsets.AccelX, offsets.AccelY, offsets.AccelZ)
	device.SetGyroOffset(offsets.GyroX, offsets.GyroY, offsets.GyroZ)
}

func (device Device) SetAccelOffset(x int32, y int32, z int32) {
	device.writeInt16Register(regXAOffsetH, int16(x))
	device.writeInt16Register(regYAOffsetH, int16(y))
	device.writeInt16Register(regZAOffsetH, int16(z))
}

func (device Device) SetGyroOffset(x int32, y int32, z int32) {
	device.writeInt16Register(regXGOffsetH, int16(x))
	device.writeInt16Register(regYGOffsetH, int16(y))
	device.writeInt16Register(regZGOffsetH, int16(z))
}

func (device Device) SetFullScaleAccelRange(registerValue uint8) {
	device.writeRegister(regAccelConfig, registerValue)
}

func AccelRangeRegister(rangeValue uint8) uint8 {
	return rangeValue << 3
}

func RawAccelToContiMS2(value int16, accelRange uint8) uint32 {
	sensitivity := int32(16384 / (1 + uint16(accelRange)))
	return int32(value) * 981 / sensitivity
}

func (device Device) readRegister(register uint8, data []byte) bool {
	return device.bus.Tx(device.address, []byte{register}, data) == nil
}

func (device Device) writeRegister(register uint8, value uint8) bool {
	return device.bus.Tx(device.address, []byte{register, value}, nil) == nil
}

func (device Device) writeInt16Register(register uint8, value int16) bool {
	encoded := uint16(value)
	return device.bus.Tx(device.address, []byte{register, byte(encoded >> 8), byte(encoded)}, nil) == nil
}
