package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var (
	ram [2048]nyb
	// vram        [2048]nyb
	rom         [49152]instr
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

	var program string
	flag.StringVar(&program, "load", "fib", "load this program from build/<arg>.bin")
	flag.StringVar(&program, "l", "fib", "load this program from build/<arg>.bin")
	flag.Parse()

	filename := fmt.Sprintf("./build/%s.bin", program)
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
