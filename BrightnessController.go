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
	backlight = "/sys/class/backlight/"
)

func getBacklightDevice() (string, error) {
	blDir := exec.Command("ls", "/sys/class/backlight/")
	op, err := blDir.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(op), "\n \r"), nil
}

type BrightnessController struct {
	device string
}

func NewBrightnessController() *BrightnessController {
	device, err := getBacklightDevice()
	if err != nil || device == "" {
		fmt.Println("no backlight devices were found!")
		return nil
	}

	return &BrightnessController{
		device: backlight + device,
	}
}

func (b *BrightnessController) GetMaxBrits() int {
	maxB, _ := os.ReadFile(path.Join(b.device, "max_brightness"))
	maxBInt, _ := strconv.ParseInt(strings.Trim(string(maxB), "\r \n"), 10, 64)
	return int(maxBInt)
}

func (b *BrightnessController) GetCurrentBrits() int {
	currB, _ := os.ReadFile(path.Join(b.device, "brightness"))
	currBInt, _ := strconv.ParseInt(strings.Trim(string(currB), "\r \n"), 10, 64)
	return int(currBInt)
}

func (b *BrightnessController) DecBrits() {
	_5britSteps := b.GetMaxBrits() / 50

	newBrit := fmt.Sprint(b.GetCurrentBrits() - _5britSteps)
	err := os.WriteFile(path.Join(b.device, "brightness"), []byte(newBrit), 0644)
	if err != nil {
		fmt.Println("ooowieee looks like you fucked up some permissions :)")
	}
}

func (b *BrightnessController) IncBrits() {
	_5britSteps := b.GetMaxBrits() / 50

	newBrit := fmt.Sprint(b.GetCurrentBrits() + _5britSteps)

	err := os.WriteFile(path.Join(b.device, "brightness"), []byte(newBrit), 0644)
	if err != nil {
		fmt.Println("ooowieee looks like you fucked up some permissions :)")
	}
}

func (b *BrightnessController) SetBrightness(brit int) {
	newBrit := (int(b.GetMaxBrits()) * brit) / 100
	err := os.WriteFile(path.Join(b.device, "brightness"), []byte(strconv.Itoa(newBrit)), 0664)
	if err != nil {
		fmt.Println("ooowieee looks like you fucked up some permissions :)")
	}
}
