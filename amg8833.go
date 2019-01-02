// Package amg8833 for interacting with the AMG8833 8x8 Thermal Camera Sensor
//
// Mostly ported from Adafruit's AMG88xx Library here:
// https://github.com/adafruit/Adafruit_AMG88xx/
package amg8833

import (
	"time"

	"golang.org/x/exp/io/i2c"
)

// General purpose consts
const (
	AMG88xxADDR                 = 0x69 // default I2C address
	AMG88xxPCTL                 = 0x00
	AMG88xxRST                  = 0x01
	AMG88xxFPSC                 = 0x02
	AMG88xxINTC                 = 0x03 // interrupt control register
	AMG88xxSTAT                 = 0x04
	AMG88xxSCLR                 = 0x05
	AMG88xxAVE                  = 0x07 // average register (for setting moving average output mode)
	AMG88xxINTHL                = 0x08 // interrupt level registers (for setting upper / lower limit hysteresis on interrupt level)
	AMG88xxINTHH                = 0x09
	AMG88xxINTLL                = 0x0A // interrupt level lower limit. Interrupt output (and interrupt pixel table are set when value is lower than set value)
	AMG88xxINTLH                = 0x0B
	AMG88xxIHYSL                = 0x0C // setting of interrupt hysteresis level when interrupt is generated. (should not be higher than interrupt level)
	AMG88xxIHYSH                = 0x0D
	AMG88xxTTHL                 = 0x0E // status register
	AMG88xxTTHH                 = 0x0F
	AMG88xxIntOffset            = 0x010
	AMG88xxPixelOffset          = 0x80 // temperature registers 0x80 - 0xFF
	AMG88xxNormalMode           = 0x00 // power: normal mode
	AMG88xxSleepMode            = 0x01 // power: sleep mode
	AMG88xxStandBy60            = 0x20 // power: stand-by mode (60 sec intermittence)
	AMG88xxStandBy10            = 0x21 // power: stand-by mode (10 sec intermittence)
	AMG88xxFlagReset            = 0x30 // resets: flag reset (all clear status reg 0x04, interrupt flag and interrupt table)
	AMG88xxInitialReset         = 0x3F // resets: initial reset (brings flag reset and returns to initial setting)
	AMG88xxFPS10                = 0x00 // frame rate: 10FPS
	AMG88xxFPS1                 = 0x01 // frame rate: 1FPS
	AMG88xxIntDisabled          = 0x00 // int enables
	AMG88xxIntEnabled           = 0x01
	AMG88xxDifference           = 0x00 // int modes
	AMG88xxAbsoluteValue        = 0x01
	AMG88xxPixelArraySize       = 64
	AMG88xxPixelTempConversion  = .25
	AMG88xxThermistorConversion = .0625
)

// Opts - set options when using the AMG88xx sensor
type Opts struct {
	Device  string
	Mode    byte
	Reset   byte
	Disable byte
	FPS     byte
}

// AMG88xx - blah
type AMG88xx struct {
	mode    byte
	reset   byte
	disable byte
	fps     byte
	dev     *i2c.Device
}

// NewAMG8833 - Begin sets up a AMG88xx chip via the I2C protocol, sets its
// gain and timing attributes, and returns an error if any occurred in that
// process or if the AMG88xx was not found.
func NewAMG8833(opts *Opts) (*AMG88xx, error) {

	device, err := i2c.Open(&i2c.Devfs{Dev: opts.Device}, int(AMG88xxADDR))
	if err != nil {
		return nil, err
	}

	amg := &AMG88xx{
		dev: device,
	}

	// enter normal mode
	amg.SetMode(opts.Mode)

	// software reset
	amg.Reset(opts.Reset)

	// disable interrupts
	amg.DisableInterrupts(0)

	// set FPS
	amg.SetFPS(opts.FPS)

	// warmup delay
	time.Sleep(100 * time.Millisecond)

	return amg, nil
}

// DisableInterrupts - disable interrupts
func (amg *AMG88xx) DisableInterrupts(disable byte) {

	write := []byte{
		disable | amg.disable,
	}
	if err := amg.dev.WriteReg(AMG88xxINTC, write); err != nil {
		panic(err)
	}

	amg.disable = disable
}

// Reset - software reset
func (amg *AMG88xx) Reset(reset byte) {

	write := []byte{
		reset | amg.reset,
	}
	if err := amg.dev.WriteReg(AMG88xxRST, write); err != nil {
		panic(err)
	}

	amg.reset = reset
}

// SetMode - set power mode
func (amg *AMG88xx) SetMode(mode byte) {

	write := []byte{
		mode | amg.mode,
	}
	if err := amg.dev.WriteReg(AMG88xxPCTL, write); err != nil {
		panic(err)
	}

	amg.mode = mode
}

// SetFPS - set frame rate
func (amg *AMG88xx) SetFPS(fps byte) {

	write := []byte{
		fps | amg.fps,
	}
	if err := amg.dev.WriteReg(AMG88xxFPSC, write); err != nil {
		panic(err)
	}

	amg.fps = fps
}

// ReadPixels - return slice of floats
func (amg *AMG88xx) ReadPixels() []float64 {

	buf := make([]float64, 64)

	for i := 0; i < 64; i++ {
		pixel := make([]byte, 1)
		amg.dev.ReadReg(AMG88xxPixelOffset+byte(i<<1), pixel)

		recast := uint16(pixel[0]<<8) | uint16(pixel[0])
		buf[i] = amg.signedMag12ToFloat(recast) * AMG88xxPixelTempConversion
	}

	return buf

}

func (amg *AMG88xx) signedMag12ToFloat(val uint16) float64 {

	//take first 11 bits as absolute val
	absVal := (val & 0x7FF)

	if val&0x8000 == 1 {
		return 0 - float64(absVal)
	}
	return float64(absVal)

}
