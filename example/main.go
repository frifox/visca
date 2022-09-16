package main

import (
	"fmt"
	"github.com/frifox/visca"
	"time"
)

func main() {
	cam := visca.Device{
		Path: "/dev/cu.usbserial-12KBBLUS",
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

	time.Sleep(time.Second)

	fmt.Printf("Applying bytes\n")
	cam.Raw.Bytes = []byte{0x1, 0x6, 0x6, 0x10}
	cam.Apply(&cam.Raw)
	time.Sleep(time.Second)

	//fmt.Println("Move Left")
	//cam.Move.X = -1
	//cam.Apply(&cam.Move)
	//time.Sleep(time.Second * 2)
	//
	//fmt.Println("Stop")
	//cam.Move.X = 0
	//cam.Apply(&cam.Move)
	//time.Sleep(time.Second)
	//
	//fmt.Println("Move Right")
	//cam.Move.X = 1
	//cam.Apply(&cam.Move)
	//time.Sleep(time.Second * 2)
	//
	//fmt.Println("Stop")
	//cam.Move.X = 0
	//cam.Apply(&cam.Move)
	//time.Sleep(time.Second)

	cam.Close()
	fmt.Printf("Camera quit\n")
}
