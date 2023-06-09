package goadb

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Device struct {
	id     string
	width  int
	height int
}

func NewDevice(id string) (*Device, error) {
	d := &Device{id: id}
	if d.checkIfExist() {
		out, err := exec.Command("adb", "-s", d.id, "shell", "wm size | cut -d' ' -f3").Output()
		if err != nil {
			return nil, err
		}
		out = TrimRightControlKey(out)
		point := strings.SplitN(string(out), "x", 2)
		width, _ := strconv.Atoi(point[0])
		height, _ := strconv.Atoi(point[1])
		d.width = width
		d.height = height
		return d, nil
	}
	return nil, fmt.Errorf("device %v not found", id)
}

func GetDevices() []string {
	return devices
}

func (d *Device) checkIfExist() bool {
	err := refreshDevices()
	if err != nil {
		return false
	}
	for _, device := range devices {
		if device == d.id {
			return true
		}
	}
	return false
}

func (d *Device) Click(x, y int) error {
	_, err := exec.Command("adb", "-s", d.id, "shell", "input", "tap", strconv.Itoa(x), strconv.Itoa(y)).Output()
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) ClickWithOffset(x, y, offset int) (int, int, error) {
	x = RandBetween(max(0, x-offset), min(d.width, x+offset))
	y = RandBetween(max(0, y-offset), min(d.height, y+offset))
	return x, y, d.Click(x, y)
}

func (d *Device) Swipe(x1, y1, x2, y2, duration int) error {
	_, err := exec.Command("adb", "-s", d.id, "shell", "input", "swipe", strconv.Itoa(x1), strconv.Itoa(y1), strconv.Itoa(x2), strconv.Itoa(y2), strconv.Itoa(duration)).Output()
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) SwipeWithOffset(x1, y1, x2, y2, duration, offset int) (int, int, int, int, error) {
	x1 = RandBetween(max(0, x1-offset), min(d.width, x1+offset))
	y1 = RandBetween(max(0, y1-offset), min(d.height, y1+offset))
	time.Sleep(time.Millisecond * 100)
	x2 = RandBetween(max(0, x2-offset), min(d.width, x2+offset))
	y2 = RandBetween(max(0, y2-offset), min(d.height, y2+offset))
	return x1, y1, x2, y2, d.Swipe(x1, y1, x2, y2, duration)
}

func (d *Device) GetWindowSize() (int, int) {
	return d.width, d.height
}

func (d *Device) GetScreenPicture() (image.Image, error) {
	out, err := exec.Command("adb", "-s", d.id, "exec-out", "screencap", "-p").Output()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(out)
	return png.Decode(buf)
}

func (d *Device) GetPixelColor(x, y int) string {
	img, err := d.GetScreenPicture()
	if err != nil {
		return ""
	}
	c := img.At(x, y).(color.NRGBA)
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

func (d *Device) GetPixelsColor(pixels [][2]int) []string {
	cap, _ := d.GetScreenPicture()
	pixelsColor := make([]string, len(pixels))
	for idx, p := range pixels {
		c := cap.At(p[0], p[1]).(color.NRGBA)
		pixelsColor[idx] = fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
	}
	return pixelsColor
}

func (d *Device) CurrentActivity() string {
	out, err := exec.Command("adb", "-s", d.id, "shell", "dumpsys activity top | grep ACTIVITY | cut -d' ' -f 4").Output()
	if err != nil {
		return ""
	}
	return string(TrimRightControlKey(out))
}

func (d *Device) IsCurrentActivity(name string) bool {
	return d.CurrentActivity() == name
}

func (d *Device) StartApplication(name string) error {
	_, err := exec.Command("adb", "-s", d.id, "shell am start -n", name).Output()
	return err
}

func (d *Device) StopApplication(name string) error {
	_, err := exec.Command("adb", "-s", d.id, "shell am force-stop", name).Output()
	return err
}

func (d *Device) PressButton(key KeyCode) error {
	_, err := exec.Command("adb", "-s", d.id, "shell", "input keyevent", strconv.Itoa(int(key))).Output()
	return err
}

func TrimRightControlKey(in []byte) []byte {
	return bytes.TrimRightFunc(in, func(r rune) bool {
		if r == '\n' || r == '\r' {
			return true
		}
		return false
	})
}
