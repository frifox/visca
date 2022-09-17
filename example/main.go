package main

import (
	"fmt"
	"github.com/frifox/visca"
	"time"
)

func main() {
	cam := visca.Device{
		Path: "udp://192.168.88.10:52381",
	}
	cam.Move.StepsY = 0x17
	cam.Move.StepsX = 0x17
	cam.Zoom.StepsZ = 7

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

	time.Sleep(time.Millisecond * 100)

	fmt.Printf("SeqReset\n")
	cam.Apply(&cam.SeqReset)

	time.Sleep(time.Millisecond * 100)

	//steps := []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 30, 40, 60, 80, 100, 120, 140, 160, 180, 200, 220, 240}
	//for _, step := range steps {
	//	cam.ZoomTo.Step = step
	//	cam.Apply(&cam.ZoomTo)
	//	time.Sleep(time.Second)
	//}

	//cam.ZoomTo.Step = 14
	//cam.Apply(&cam.ZoomTo)
	//time.Sleep(time.Second)

	//cam.Raw.Bytes = []byte{0x81, 0x1, 0x7e, 0x4, 0x3b, 0x2, 0xff}
	//cam.Apply(&cam.Raw)
	//time.Sleep(time.Second / 2)

	cam.Move.Y = 1.0 / 0x17
	cam.Apply(&cam.Move)
	time.Sleep(time.Second / 2)
	cam.Move.Y = 0.0
	cam.Apply(&cam.Move)
	time.Sleep(time.Second / 2)

	cam.Move.Y = -1.0 / 0x17
	cam.Apply(&cam.Move)
	time.Sleep(time.Second / 2)
	cam.Move.Y = 0.0
	cam.Apply(&cam.Move)
	time.Sleep(time.Second / 2)
	//
	//cam.Move.X = 1.0 / 0x17
	//cam.Apply(&cam.Move)
	//time.Sleep(time.Second / 2)
	//cam.Move.X = 0.0
	//cam.Apply(&cam.Move)
	//time.Sleep(time.Second / 2)
	//
	//cam.Move.X = -1.0 / 0x17
	//cam.Apply(&cam.Move)
	//time.Sleep(time.Second / 2)
	//cam.Move.X = 0.0
	//cam.Apply(&cam.Move)
	//time.Sleep(time.Second / 2)

	//fmt.Println("Move Left")
	//cam.Move.X = -1
	//cam.Apply(&cam.Move)
	//time.Sleep(time.Second / 2)
	//
	//fmt.Println("Stop")
	//cam.Move.X = 0
	//cam.Apply(&cam.Move)
	//time.Sleep(time.Second / 2)
	//
	//fmt.Println("Move Right")
	//cam.Move.X = 1
	//cam.Apply(&cam.Move)
	//time.Sleep(time.Second / 2)
	//
	//fmt.Println("Stop")
	//cam.Move.X = 0
	//cam.Apply(&cam.Move)
	//time.Sleep(time.Second / 2)

	time.Sleep(time.Second / 2)
	cam.Close()
	fmt.Printf("Camera quit\n")
}
