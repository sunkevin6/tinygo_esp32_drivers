package ov2640

import "machine"

type Pins struct {
	PWDN  machine.Pin
	Reset machine.Pin
	XCLK  machine.Pin
	SIOD  machine.Pin
	SIOC  machine.Pin
	Y9    machine.Pin
	Y8    machine.Pin
	Y7    machine.Pin
	Y6    machine.Pin
	Y5    machine.Pin
	Y4    machine.Pin
	Y3    machine.Pin
	Y2    machine.Pin
	VSYNC machine.Pin
	HREF  machine.Pin
	PCLK  machine.Pin
	LED   machine.Pin
}

type Config struct {
	Pins             Pins
	XCLKFreqHz       int
	FrameSize        string
	PixelFormat      string
	GrabMode         string
	FrameBufferPlace string
	JPEGQuality      int
	FrameBufferCount int
}

var DefaultPins = Pins{
	PWDN:  machine.NoPin,
	Reset: machine.NoPin,
	XCLK:  machine.GPIO10,
	SIOD:  machine.GPIO40,
	SIOC:  machine.GPIO39,
	Y9:    machine.GPIO48,
	Y8:    machine.GPIO11,
	Y7:    machine.GPIO12,
	Y6:    machine.GPIO14,
	Y5:    machine.GPIO16,
	Y4:    machine.GPIO18,
	Y3:    machine.GPIO17,
	Y2:    machine.GPIO15,
	VSYNC: machine.GPIO38,
	HREF:  machine.GPIO47,
	PCLK:  machine.GPIO13,
	LED:   machine.GPIO21,
}

func DefaultConfig() Config {
	return Config{
		Pins:             DefaultPins,
		XCLKFreqHz:       20000000,
		FrameSize:        "UXGA",
		PixelFormat:      "JPEG",
		GrabMode:         "CAMERA_GRAB_WHEN_EMPTY",
		FrameBufferPlace: "PSRAM",
		JPEGQuality:      12,
		FrameBufferCount: 1,
	}
}

func RoundDisplayConfig() Config {
	return Config{
		Pins:             DefaultPins,
		XCLKFreqHz:       10000000,
		FrameSize:        "FRAMSIZE_240X240",
		PixelFormat:      "RGB565",
		GrabMode:         "CAMERA_GRAB_WHEN_EMPTY",
		FrameBufferPlace: "PSRAM",
		JPEGQuality:      12,
		FrameBufferCount: 2,
	}
}
