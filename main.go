package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	ram [2048]nyb
	// vram        [2048]nyb
	rom         [49152 * 2]instr
	cache       [16]ptr
	cacheOffset ptr

	// Registers
	a  nyb
	b  nyb
	t0 nyb // tmp 1
	t1 nyb // tmp 2

	carry bool
	pc    ptr
)

func main() {
	fmt.Println("Populating byte -> nyb map")
	populateByte2Nyb()

	fmt.Println("Loading binary program:")
	file, _ := os.Open("./build/fib.bin")
	defer file.Close()

	stats, _ := file.Stat()
	byts := make([]byte, stats.Size())

	reader := bufio.NewReader(file)
	reader.Read(byts)

	for i, byt := range byts {
		rom[i] = instr(byt)
	}

	printRom(0, 128)

	fmt.Printf("Running program from ROM (%d bytes)...\n", len(byts))
	for int(pc) < len(byts) {
		// TODO: debug mode, pause at breakpoints, allow line-by-line stepping
		cycle()
	}

	fmt.Println("Final state of the machine:")
	printRegs()
	printCache()
	printRam(0, 32)
	fmt.Printf("PC: %d\nCarry: %v\nCache offset: %d\n", pc, carry, cacheOffset)

}
