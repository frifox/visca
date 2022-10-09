package main

import (
	"fmt"
	"github.com/frifox/visca"
	"time"
)

func main() {
	cam := &visca.Device{
		Path: "udp://10.0.0.213:52381",
		Config: visca.Config{
			LocalUDP:  ":52381",
			XMaxSpeed: 20.0 / 24.0,
			YMaxSpeed: 20.0 / 24.0,
			ZMaxSpeed: 8.0 / 8.0,
		},
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
	time.Sleep(time.Millisecond * 100)

	cam.Do(&visca.InqKneeSlope{})
	time.Sleep(time.Millisecond * 100)
	fmt.Printf("%s\n", cam.Inquiry.InqKneeSlope)

	//cam.Close()
	//fmt.Printf("Camera done\n")
}
