package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/bitinvert/Many_Sim/helper"
	"github.com/bitinvert/Many_Sim/processor"
)

func newCPU() processor.Cpu {
	return processor.Cpu{
		PC:     0,
		DStack: helper.NewStack(processor.DataStackDepth),
		RStack: helper.NewStack(processor.ReturnStackDepth),
		Memory: make([]uint16, processor.MemorySize, processor.MemorySize),
		Tick:   0}
}

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1:]

		fmt.Println("Starting CPU!")
		core := newCPU()

		image := helper.LoadFile(args[0])

		core.LoadImage(image)
		processor.Loading = len(args) == 3
		if processor.Loading {

			if len(args) > 1 {
				f1, err := os.Open(args[1])

				defer f1.Close()

				if err != nil {
					panic(err)
				}
				processor.Sourcecode = bufio.NewReader(f1)
			}
			for /*loading*/ i := 0; i < (600 - 44); i++ {
				start := time.Now()
				core.CpuStep()
				fmt.Println(time.Since(start))

			}

			fmt.Println("\n\nDumping memory image!")
			core.SaveImage(args[2])

		} else {
			for processor.Currentkey != 27 {
				core.CpuStep()
			}
		}

	} else {
		fmt.Println("Many_Sim: A Simulator for the J1a")
		fmt.Println("Usage:")
		fmt.Println("go run main.go image                     - Interactive Mode")
		fmt.Println("go run main.go image-in source image-out - Image generation mode")
	}
}
