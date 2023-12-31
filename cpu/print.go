package cpu

import "fmt"

func PrintRAM(from, to, width int) {
	out := ""
	for i := from; i < to; i += width {
		for j := 0; j < width; j += 1 {
			out += fmt.Sprintf("%x ", Nyb2Byte[RAM[i+j]])
		}
		out += "\n"
	}

	fmt.Print(out)
}

func PrintROM(from, to int) {
	out := ""
	for i := from; i < to; i += 8 {
		for j := 0; j < 8; j += 1 {
			out += fmt.Sprintf("%02x ", ROM[i+j])
		}
		out += "\n"
	}

	fmt.Print(out)
}

func PrintRegisters() {
	fmt.Printf("A: %x, B: %x | T: %x, %x\n", Nyb2Byte[A], Nyb2Byte[B], Nyb2Byte[T0], Nyb2Byte[T1])
}

func PrintMemCache() {

	for i := 0; i < 8; i += 1 {
		fmt.Printf("$%04x(%x) ", MemCache[i], Nyb2Byte[RAM[MemCache[i]]])
	}
	fmt.Println()
	for i := 8; i < 16; i += 1 {
		fmt.Printf("$%04x(%x) ", MemCache[i], Nyb2Byte[RAM[MemCache[i]]])
	}
	fmt.Println()
}

func PrintInstruction(i Instr) string {
	op := byte(i / 16)
	arg := byte(i % 16)

	opStr := opMap[op]

	if tillNextOp == 0 {
		return fmt.Sprintf("%s %x", opStr, arg)
	} else {
		return fmt.Sprintf("%x %x", op, arg)
	}
}

func PrintCallStack() {
	if len(CallStack) == 0 {
		fmt.Println("<empty call stack>")
		return
	}

	out := "Call stack: "

	for _, ptr := range CallStack {
		out += fmt.Sprintf("%d > ", ptr)
	}
	out = out[:len(out)-4]
	fmt.Println(out)
}
