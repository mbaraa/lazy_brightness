package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const (
	BACKLIGHT_DIR = "/sys/class/backlight/"
)

func getBacklightDevices() ([]string, error) {
	blDir := exec.Command("ls", BACKLIGHT_DIR)
	op, err := blDir.Output()
	if err != nil {
		return nil, err
	}

	op = []byte(strings.Trim(string(op), "\n\r"))
	return strings.Split(string(op), "\n"), nil
}

type BrightnessController struct {
	devices       []string
	currentDevice string
}

func NewBrightnessController() *BrightnessController {
	devices, err := getBacklightDevices()
	if err != nil || devices == nil {
		fmt.Println("no backlight devices were found!")
		return nil
	}

	return &BrightnessController{
		devices:       devices,
		currentDevice: devices[0],
	}
}

func (b *BrightnessController) GetMaxBrits() int {
	maxB, _ := os.ReadFile(path.Join(BACKLIGHT_DIR+b.currentDevice, "max_brightness"))
	maxBInt, _ := strconv.ParseInt(strings.Trim(string(maxB), "\r \n"), 10, 64)
	return int(maxBInt)
}

func (b *BrightnessController) GetCurrentBrits() int {
	currB, _ := os.ReadFile(path.Join(BACKLIGHT_DIR+b.currentDevice, "brightness"))
	currBInt, _ := strconv.ParseInt(strings.Trim(string(currB), "\r \n"), 10, 64)
	return int(currBInt)
}

func (b *BrightnessController) DecBrits() {
	_5britSteps := b.GetMaxBrits() / 50

	newBrit := fmt.Sprint(b.GetCurrentBrits() - _5britSteps)
	err := os.WriteFile(path.Join(BACKLIGHT_DIR+b.currentDevice, "brightness"), []byte(newBrit), 0644)
	if err != nil {
		fmt.Printf("ooowieee looks like you fucked up some permissions :), err: %s\n", err.Error())
	}
}

func (b *BrightnessController) IncBrits() {
	_5britSteps := b.GetMaxBrits() / 50

	newBrit := fmt.Sprint(b.GetCurrentBrits() + _5britSteps)

	err := os.WriteFile(path.Join(BACKLIGHT_DIR+b.currentDevice, "brightness"), []byte(newBrit), 0644)
	if err != nil {
		fmt.Printf("ooowieee looks like you fucked up some permissions :), err: %s\n", err.Error())
	}
}

func (b *BrightnessController) SetBrightness(brit int) {
	newBrit := (int(b.GetMaxBrits()) * brit) / 100
	err := os.WriteFile(path.Join(BACKLIGHT_DIR+b.currentDevice, "brightness"), []byte(strconv.Itoa(newBrit)), 0664)
	if err != nil {
		fmt.Printf("ooowieee looks like you fucked up some permissions :), err: %s\n", err.Error())
	}
}

func (b *BrightnessController) SelectDevice(name string) {
	for _, device := range b.devices {
		if device == name {
			b.currentDevice = device
			return
		}
	}
}

func (b *BrightnessController) GetDevices() []string {
	return b.devices
}
