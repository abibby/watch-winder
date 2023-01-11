package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/easystepper"
)

func main() {
	machine.InitADC()

	led := machine.LED
	in1 := machine.Pin(14)
	in2 := machine.Pin(15)
	in3 := machine.Pin(16)
	in4 := machine.Pin(17)
	speed := machine.ADC{machine.ADC2}

	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	in1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	in2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	in3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	in4.Configure(machine.PinConfig{Mode: machine.PinOutput})

	stepper, err := easystepper.New(easystepper.DeviceConfig{
		Pin1:      in1,
		Pin2:      in3,
		Pin3:      in2,
		Pin4:      in4,
		StepCount: 2038,
		RPM:       5,
		Mode:      easystepper.ModeFour,
	})
	if err != nil {
		for {
			blink(led, ".-")
		}
	}

	turnsPerDay := float32(0.0)
	for {
		turnsPerDay = normaliseInput(speed.Get())
		stepper.Move(2038)
		stepper.Off()
		time.Sleep((time.Hour * 24) / time.Duration(turnsPerDay))

		turnsPerDay = normaliseInput(speed.Get())
		stepper.Move(-2038)
		stepper.Off()
		time.Sleep((time.Hour * 24) / time.Duration(turnsPerDay))
	}
}

func blink(led machine.Pin, str string) {
	for _, c := range str {
		switch c {
		case '.', '0':
			led.Low()
			time.Sleep(time.Millisecond * 500)

			led.High()
			time.Sleep(time.Millisecond * 250)
		case '-', '1':
			led.Low()
			time.Sleep(time.Millisecond * 500)

			led.High()
			time.Sleep(time.Millisecond * 750)
		}
	}
	led.Low()
}

// ////////////.... 65520
// //./.... 208

func normaliseInput(inputValue uint16) float32 {
	max := float32(1000)
	min := float32(500)
	return (float32(inputValue)/float32(0xffff))*(max-min) + min // ADC ranges from 0..0xffff
}
