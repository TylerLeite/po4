package cpu

import (
	"bufio"
	"fmt"
	"os"
)

var (
	RAM       [2048]Nyb
	ROM       [49152]Instr
	MemCache  [16]Ptr
	CallStack = make([]Ptr, 0)

	// Registers
	A  Nyb
	B  Nyb
	T0 Nyb // tmp 1
	T1 Nyb // tmp 2

	Carry          bool
	ProgramCounter Ptr
)

func PowerOn() {
	// populate Byte2Nyb
	fmt.Println("Populating byte -> nyb map")
	for i := byte(16); i > 0; i += 1 {
		Byte2Nyb[i] = Byte2Nyb[i%16]
	}
}

func Load(filename string) int {
	fmt.Printf("Loading binary program: %s\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		panic("Error loading file :(")
	}
	defer file.Close()

	stats, _ := file.Stat()
	byts := make([]byte, stats.Size())

	reader := bufio.NewReader(file)
	reader.Read(byts)

	for i, byt := range byts {
		ROM[i] = Instr(byt)
	}

	programSize := len(byts)
	PrintROM(0, programSize)
	return programSize
}

func Run(programSize int) {
	fmt.Printf("Running program from ROM (%d bytes)...\n", programSize)

	for int(ProgramCounter) < programSize {
		// TODO: debug mode, pause at breakpoints, allow line-by-line stepping
		cycle()
	}
}

func PowerOff() {
	fmt.Println("Final state of the machine:")
	PrintRegisters()
	PrintMemCache()
	PrintRAM(0, 2048, 64)
	fmt.Printf("PC: %d\nCarry: %v\n", ProgramCounter, Carry)
}
