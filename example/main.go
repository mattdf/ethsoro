package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/chzyer/readline"
	"github.com/conejoninja/trezor"
	"github.com/zserge/hid"
)

func main() {

	var c trezor.TrezorClient

	numberDevices := 0
	hid.UsbWalk(func(device hid.Device) {
		info := device.Info()
		// 0x534c : 21324 vendor
		// 0x0001 : 1     product
		if info.Vendor == 21324 && info.Product == 1 {
			numberDevices++
			c.SetTransport(device)
		}
	})
	if numberDevices == 0 {
		fmt.Println("No TREZOR devices found, make sure your TREZOR device is connected")
	} else {
		fmt.Printf("Found %d TREZOR devices connected\n", numberDevices)
		//c.Initialize()
		c.Ping("hola")
		c.GetAddress()
		shell(c)
		defer c.CloseTransport()
	}
}

func shell(c trezor.TrezorClient) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt: ">",
	})
	if err != nil {
		panic(err)
	}

	defer rl.Close()
	log.SetOutput(rl.Stderr())

out:
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		line = strings.ToLower(line)
		if line != "" {
			c.PinMatrixAck(line)
			continue out
		}
	}
}