package processor

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/bitinvert/Many_Sim/helper"
)

// Constants and globals needed for the simulator
const DataStackDepth uint16 = 256
const ReturnStackDepth uint16 = 32
const MemorySize uint16 = 0x1FFF
const KiloByte uint32 = 1024

var Keypending bool = true
var Currentkey byte = 0
var Loading bool = false
var Sourcecode *bufio.Reader

// Cpu is the heart of the project
type Cpu struct {
	PC     uint16
	DStack *helper.Stack
	RStack *helper.Stack
	Memory []uint16
	Tick   uint16
}

// LoadImage loads the image into memory
func (c *Cpu) LoadImage(fulldata [][]byte) {
	addr := 0
	for i := 0; i < len(fulldata); i++ {
		c.Memory[addr] = helper.Hex2word(fulldata[i])
		addr++
	}
}

// SaveImage takes the current memory and creates a memorydump, this is only needed
func (c *Cpu) SaveImage(fname string) {
	file, err := os.Create(fname)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	for k := 0; k < len(c.Memory); k++ {
		value := helper.Word2hex(c.Memory[k])
		file.WriteString(value + "\n")
	}
}

// WriteIO simulates UART an memory mapped IO
func (c *Cpu) WriteIO(addr uint16, data uint16) {
	if (addr & 0x1000) != 0 {
		fmt.Println((data & 0xFF))
	}
	if (addr & 0x4000) != 0 {
		c.Tick = data
	}
}

// ReadIO simulates the reading part UART. In this case it is only for fetching the source file specified at the start of the program.
// It could be also used for fetching user input, though not implemented
func (c *Cpu) ReadIO(addr uint16) uint16 {
	var ortogether uint16 = 0
	r := rand.New(rand.NewSource(0))
	addr = helper.WordMalformed(addr)
	if (addr & 0x1000) != 0 {
		if Keypending == true {
			Keypending = false
		} else {
			if Loading == true {
				Currentkey = 0
				if Sourcecode.Buffered() == 1 {
					Loading = false
				} else {
					value, err := Sourcecode.ReadByte()
					if err != nil {
						panic(err)
					}
					Currentkey = value
					fmt.Printf("%s", string([]byte{Currentkey}))
					os.Stdout.Sync()

				}
			} else {
				// Do nothing as no userinput is allowed right now

			}
			ortogether = ortogether | uint16(Currentkey)
		}
	}

	if (addr & 0x2000) != 0 {

		ortogether = ortogether | 0x0003
		if r.Intn(3) == 1 {
			ortogether = ortogether | 0x0004
		}
	}

	if (addr & 0x4000) != 0 {

		ortogether = ortogether | c.Tick
	}
	return ortogether
}

// CpuStep is the heart of the simulator. It decides what opcodes do
func (c *Cpu) CpuStep() {

	var tosN uint16
	c.Tick = (c.Tick + 1) & 0xFFFF
	insn := c.Memory[c.PC&uint16(MemorySize)] // High-Call = Memory Fetch, mask away the topmost address bit.

	if (c.PC & 0x2000) != 0 { // Memory fetch
		c.DStack.Push(insn)
		c.PC = c.RStack.Pop() >> 1
	} else {
		if (insn & 0x8000) != 0 { // Literal
			c.DStack.Push(insn & 0x7FFF)
			c.PC++
		} else {
			switch insn & 0xE000 {
			case 0x0000: // Jump
				c.PC = insn & uint16(MemorySize)
			case 0x2000: // Conditional Jump
				if c.DStack.Pop() == 0 {
					c.PC = insn & uint16(MemorySize)
				} else {
					c.PC++
				}
			case 0x4000: // Call
				c.RStack.Push((c.PC + 1) << 1)
				c.PC = insn & uint16(MemorySize)
			case 0x6000: // ALU
				rtos := c.RStack.Tos()
				tos := c.DStack.Tos()
				nos := c.DStack.Nos()
				switch insn & 0x0F00 {
				case 0x0000:
					tosN = tos
				case 0x0100:
					tosN = nos
				case 0x0200:
					tosN = tos + nos
				case 0x0300:
					tosN = tos & nos
				case 0x0400:
					tosN = tos | nos
				case 0x0500:
					tosN = tos ^ nos
				case 0x0600:
					tosN = ^tos
				case 0x0700:
					tosN = helper.Flag(nos == tos)
				case 0x0800:
					tosN = helper.Flag(helper.Signed(nos) < helper.Signed(tos))
				case 0x0900:
					tosN = (tos >> 1) | (tos & 0x8000)
				case 0x0A00:
					tosN = tos << 1
				case 0x0B00:
					tosN = rtos
				case 0x0C00:
					tosN = nos - tos
				case 0x0D00:
					tosN = c.ReadIO(tos)
				case 0x0E00:
					tosN = c.DStack.Cap
				case 0x0F00:
					tosN = helper.Flag(nos < tos)
				}

				// ALU exit bit
				if (insn & 0x0080) != 0 {
					c.PC = rtos >> 1
				} else {
					c.PC++
				}

				// Data stack movement (d-1, d-2, d+1)
				if (insn & 0x0003) == 0x0003 {
					c.DStack.Pop()
				}
				if (insn & 0x0003) == 0x0002 {
					c.DStack.Pop()
					c.DStack.Pop()
				}
				if (insn & 0x0003) == 0x0001 {
					c.DStack.Push(tos)
				}

				// Return stack movement (r-1, r-2, r+1)
				if (insn & 0x000C) == 0x000C {
					c.RStack.Pop()
				}
				if (insn & 0x000C) == 0x0008 {
					c.RStack.Pop()
					c.RStack.Pop()
				}
				if (insn & 0x000C) == 0x0004 {
					c.RStack.Push(tos)
				}

				// Read & Write operations (T -> N, T -> R, Memory write, IO write)
				if (insn & 0x0070) == 0x0010 {
					c.DStack.Elements[1] = tos
				}
				if (insn & 0x0070) == 0x0020 {
					c.RStack.Push(tos)
				}
				if (insn & 0x0070) == 0x0030 {
					//fmt.Printf("writing at %#x: %#x\n", tos>>1, nos)
					c.Memory[(tos >> 1)] = nos
				}
				if (insn & 0x0070) == 0x0040 {
					c.WriteIO(tos, nos)
				}
				c.DStack.Push(tosN)
			}
		}
	}
}
