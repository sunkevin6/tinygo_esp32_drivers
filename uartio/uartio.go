package uartio

import (
	"machine"

	"github.com/sunkevin6/tinygo_esp32_drivers/numfmt"
)

type Port struct {
	uart machine.Serialer
}

func New(uart machine.Serialer) Port {
	return Port{uart: uart}
}

func (port Port) Configure(config machine.UARTConfig) error {
	return port.uart.Configure(config)
}

func (port Port) Write(data []byte) (int, error) {
	return port.uart.Write(data)
}

func (port Port) WriteByte(value byte) error {
	return port.uart.WriteByte(value)
}

func (port Port) Buffered() int {
	return port.uart.Buffered()
}

func (port Port) ReadByte() (byte, error) {
	return port.uart.ReadByte()
}

func (port Port) Clear() {
	for port.Buffered() > 0 {
		_, _ = port.ReadByte()
	}
}

func (port Port) Line(message string) {
	port.String(message)
	port.String("\r\n")
}

func (port Port) String(message string) {
	port.uart.Write([]byte(message))
}

func (port Port) Bool(value bool) {
	if value {
		port.String("true")
	} else {
		port.String("false")
	}
}

func (port Port) BoolInt(value bool) {
	digitByte := byte('0')
	if value {
		digitByte = '1'
	}

	port.uart.WriteByte(digitByte)
}

func (port Port) Int(value int) {
	var buffer [11]byte // enough for int32
	digitCount := numfmt.AppendInt(buffer[:0], value)
	port.uart.Write(buffer[:digitCount])
}

func (port Port) Int32(value int32) {
	var buffer [11]byte // enough for int32
	digitCount := numfmt.AppendInt32(buffer[:0], value)
	port.uart.Write(buffer[:digitCount])
}

func (port Port) Uint32(value uint32) {
	var buffer [10]byte // enough for uint32
	digitCount := numfmt.AppendUint32(buffer[:0], value)
	port.uart.Write(buffer[:digitCount])
}

func (port Port) HexByte(value uint16) {
	port.uart.WriteByte(numfmt.HexDigit(uint8(value >> 4 & 0x0F)))
	port.uart.WriteByte(numfmt.HexDigit(uint8(value & 0x0F)))
}

func (port Port) Hex8(value uint8) {
	if value >= 16 {
		port.uart.WriteByte(numfmt.HexDigit(value >> 4))
	}
	port.uart.WriteByte(numfmt.HexDigit(value & 0x0F))
}

func (port Port) Hex32(value uint32) {
	started := false
	for shift := int32(28); shift >= 0; shift -= 4 {
		nibble := uint8(value >> shift & 0x0F)
		if nibble != 0 || started || shift == 0 {
			port.uart.WriteByte(numfmt.HexDigit(nibble))
			started = true
		}
	}
}

func (port Port) FloatCenti(value float64) {
	var buffer [16]byte // enough for float64
	digitCount := numfmt.AppendFloatCenti(buffer[:], value)
	port.uart.Write(buffer[:digitCount])
}

func (port Port) FixedCenti(value int32) {
	var buffer [12]byte // enough for int32 with decimal point and sign
	digitCount := numfmt.AppendFixedCenti(buffer[:], value)
	port.uart.Write(buffer[:digitCount])
}

func (port Port) FixedCentiScaled(value int32, scale int32) {
	var buffer [12]byte // enough for int32 with decimal point and sign
	digitCount := numfmt.AppendFixedCentiScaled(buffer[:], value, scale)
	port.uart.Write(buffer[:digitCount])
}

func (port Port) FixedDeci(value int32) {
	var buffer [12]byte // enough for int32 with decimal point and sign
	digitCount := numfmt.AppendFixedDeci(buffer[:], value)
	port.uart.Write(buffer[:digitCount])
}

func (port Port) Pin(label string, pin machine.Pin) {
	port.String(label)
	port.String(": ")
	if pin == machine.NoPin {
		port.String("NoPin")
		return
	}
	port.String("GPIO")
	port.Int(int(pin))
	port.Line("")
}
