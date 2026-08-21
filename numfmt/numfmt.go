package numfmt

func AppendInt(buffer []byte, value int) int {
	return AppendIntAt(buffer, 0, value)
}

func AppendInt32(buffer []byte, value int32) int {
	position := 0
	if value < 0 {
		position = AppendByte(buffer, position, '-')
		return appendUint32At(buffer, position, uint32(-(value+1))+1)
	}
	return appendUint32At(buffer, position, uint32(value))
}

func AppendUint32(buffer []byte, value uint32) int {
	return appendUint32At(buffer, 0, value)
}

func Unit32(value uint32) string {
	var digits [10]byte // enough for uint32
	digitCount := AppendUint32(digits[:], value)
	return string(digits[:digitCount])
}

func FixedDeci(value int32) string {
	var digits [12]byte // enough for int32 with decimal point and sign
	digitCount := AppendFixedDeci(digits[:], value)
	return string(digits[:digitCount])
}

func AppendByte(buffer []byte, position int, value byte) int {
	if position < len(buffer) {
		buffer[position] = value
		position++
	}
	return position
}

func AppendString(buffer []byte, position int, value string) int {
	for i := 0; i < len(value); i++ {
		position = AppendByte(buffer, position, value[i])
	}
	return position
}

func AppendIntAt(buffer []byte, position int, value int) int {
	if value < 0 {
		position = AppendByte(buffer, position, '-')
		value = -value
	}
	return appendUint32At(buffer, position, uint32(value))
}

func AppendFixedCenti(buffer []byte, value int32) int {
	position := 0
	if value < 0 {
		position = AppendByte(buffer, position, '-')
		value = -value
	}

	position = AppendInt32At(buffer, position, value/100)
	position = AppendByte(buffer, position, '.')
	fraction := value % 100
	if fraction < 10 {
		position = AppendByte(buffer, position, '0')
	}
	return AppendInt32At(buffer, position, fraction)
}

func AppendFixedCentiScaled(buffer []byte, value int32, scale int32) int {
	return AppendFixedDeci(buffer, RoundScaled(value, scale))
}

func AppendFixedDeci(buffer []byte, value int32) int {
	position := 0
	if value < 0 {
		position = AppendByte(buffer, position, '-')
		value = -value
	}

	position = AppendInt32At(buffer, position, value/10)
	position = AppendByte(buffer, position, '.')
	return AppendInt32At(buffer, position, value%10)
}

func AppendFloatCenti(buffer []byte, value float64) int {
	position := 0
	if value < 0 {
		position = AppendByte(buffer, position, '-')
		value = -value
	}

	whole := int(value)
	fraction := int((value - float64(whole)) * 100)
	position = AppendIntAt(buffer, position, whole)
	position = AppendByte(buffer, position, '.')
	if fraction < 10 {
		position = AppendByte(buffer, position, '0')
	}
	return AppendIntAt(buffer, position, fraction)
}

func HexDigit(value uint8) byte {
	if value < 10 {
		return '0' + value
	}
	return 'A' + (value - 10)
}

func RoundScaled(value int32, scale int32) int32 {
	if scale == 0 {
		return value
	}

	result := value / scale
	remainder := value % scale
	halfScale := scale / 2
	if value >= 0 && remainder >= halfScale {
		result++
	} else if value < 0 && remainder <= -halfScale {
		result--
	}
	return result
}

func AppendInt32At(buffer []byte, position int, value int32) int {
	if value < 0 {
		position = AppendByte(buffer, position, '-')
		return appendUint32At(buffer, position, uint32(-(value+1))+1)
	}
	return appendUint32At(buffer, position, uint32(value))
}

func appendUint32At(buffer []byte, position int, value uint32) int {
	if value == 0 {
		return AppendByte(buffer, position, '0')
	}

	var digits [10]byte // enough for uint32
	digitPosition := len(digits)
	for value > 0 {
		digitPosition--
		digits[digitPosition] = byte('0' + (value % 10))
		value /= 10
	}

	for digitPosition < len(digits) {
		position = AppendByte(buffer, position, digits[digitPosition])
		digitPosition++
	}

	return position
}
