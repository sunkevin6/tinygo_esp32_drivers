package pcf8563

const (
	address               = 0x51
	controlStatus1Reg     = 0x00
	secondsReg            = 0x02
	controlStatus1StopBit = 0x20
)

type Time struct {
	Year   uint8
	Month  uint8
	Day    uint8
	Hour   uint8
	Minute uint8
	Second uint8
}

type I2C interface {
	Tx(addr uint16, write []byte, read []byte) error
}

type Device struct {
	bus I2C
}

func New(bus I2C) Device {
	return Device{bus: bus}
}

func (device Device) Init() error {
	return device.writeRegister(controlStatus1Reg, 0)
}

func (device Device) Stop() error {
	return device.writeRegister(controlStatus1Reg, controlStatus1StopBit)
}

func (device Device) Start() error {
	return device.writeRegister(controlStatus1Reg, 0)
}

func (device Device) SetTime(value Time) error {
	data := []byte{
		secondsReg,
		decimalToBCD(value.Second),
		decimalToBCD(value.Minute),
		decimalToBCD(value.Hour),
		decimalToBCD(value.Day),
		0,
		decimalToBCD(value.Month),
		decimalToBCD(value.Year),
	}
	return device.bus.Tx(address, data, nil)
}

func (device Device) GetTime() (Time, error) {
	data := []byte{0, 0, 0, 0, 0, 0, 0}
	if err := device.bus.Tx(address, []byte{secondsReg}, data); err != nil {
		return Time{}, err
	}

	return Time{
		Second: bcdToDecimal(data[0] & 0x7f),
		Minute: bcdToDecimal(data[1] & 0x7f),
		Hour:   bcdToDecimal(data[2] & 0x3f),
		Day:    bcdToDecimal(data[3] & 0x3f),
		Month:  bcdToDecimal(data[5] & 0x1f),
		Year:   bcdToDecimal(data[6]),
	}, nil
}

func (device Device) writeRegister(register byte, value byte) error {
	return device.bus.Tx(address, []byte{register, value}, nil)
}

func decimalToBCD(value uint8) uint8 {
	return value/10<<4 | value%10
}

func bcdToDecimal(value uint8) uint8 {
	return (value>>4)*10 + (value & 0x0f)
}
