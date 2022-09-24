package main

import (
	"fmt"
	"github.com/frifox/visca"
	"time"
)

func main() {
	cam := visca.Device{
		Path: "udp://10.0.0.10:52381",
		Type: visca.SonySRGX400,
	}

	fmt.Printf("Looking for camera at %s\n", cam.Path)
	err := cam.Find()
	if err != nil {
		fmt.Printf("Find: %v\n", err)
		return
	}
	if !cam.Found() {
		fmt.Printf("Camera 404\n")
		return
	}

	fmt.Printf("Running camera\n")
	go cam.Run()
	cam.Booting.Wait()

	// init
	cam.Do(&visca.SeqReset{})
	cam.Do(&visca.Power{On: true})

	// for manual PTZ recall
	pos := visca.InqPanTiltPosition{}
	cam.Do(&pos)
	zoom := visca.InqZoom{}
	cam.Do(&zoom)

	// zoom
	cam.Do(&visca.Zoom{Z: 1})
	time.Sleep(time.Second / 2)
	cam.Do(&visca.Zoom{Z: -1})
	time.Sleep(time.Second / 2)
	cam.Do(&visca.Zoom{})
	time.Sleep(time.Second / 2)

	// pan
	cam.Do(&visca.PanTiltDrive{X: 0.5})
	time.Sleep(time.Second / 2)
	cam.Do(&visca.PanTiltDrive{X: -0.5})
	time.Sleep(time.Second / 2)
	cam.Do(&visca.PanTiltDrive{})
	time.Sleep(time.Second / 2)

	// tilt
	cam.Do(&visca.PanTiltDrive{Y: 0.5})
	time.Sleep(time.Second / 2)
	cam.Do(&visca.PanTiltDrive{Y: -0.5})
	time.Sleep(time.Second / 2)
	cam.Do(&visca.PanTiltDrive{})
	time.Sleep(time.Second / 2)

	// set preset
	cam.Do(&visca.PresetSet{
		ID: 1,
	})

	// recall abs PTZ
	cam.Do(&visca.PanTiltDriveAbs{
		X: pos.X, SpeedX: 1,
		Y: pos.Y, SpeedY: 1,
	})
	cam.Do(&visca.ZoomAbs{
		Z: zoom.Z,
	})
	time.Sleep(time.Second)

	// recall preset
	cam.Do(&visca.PresetRecall{
		ID: 1,
	})
	time.Sleep(time.Second)

	// TODO
	// menu toggle
	// enter
	// arrow keys

	// TODO
	// state OnChange closures

	cam.Close()
	fmt.Printf("Camera done\n")
}
