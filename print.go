package main

import "fmt"

func printRam(from, to int) {
	out := ""
	for i := from; i < to; i += 8 {
		for j := 0; j < 8; j += 1 {
			out += fmt.Sprintf("%x ", nyb2byte[ram[i+j]])
		}
		out += "\n"
	}

	fmt.Print(out)
}

func printRom(from, to int) {
	out := ""
	for i := from; i < to; i += 8 {
		for j := 0; j < 8; j += 1 {
			out += fmt.Sprintf("%02x ", rom[i+j])
		}
		out += "\n"
	}

	fmt.Print(out)
}

func printRegs() {
	fmt.Printf("A: %x, B: %x | T: %x\n", nyb2byte[a], nyb2byte[b], nyb2byte[t0])
}

func printCache() {

	for i := 0; i < 8; i += 1 {
		fmt.Printf("$%04x(%x) ", cache[i], nyb2byte[ram[cache[i]]])
	}
	fmt.Println()
	for i := 8; i < 16; i += 1 {
		fmt.Printf("$%04x(%x) ", cache[i], nyb2byte[ram[cache[i]]])
	}
	fmt.Println()
}

func printInstruction(i instr) string {
	op := byte(i / 16)
	arg := byte(i % 16)

	opStr := opMap[op]

	if tillNextOp == 0 {
		return fmt.Sprintf("%s %x", opStr, arg)
	} else {
		return fmt.Sprintf("%x %x", op, arg)
	}
}
