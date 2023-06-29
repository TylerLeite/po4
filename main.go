package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type nyb [4]bool
type ptr int16
type instr uint8

func (n *nyb) set(to nyb) {
	n[0] = to[0]
	n[1] = to[1]
	n[2] = to[2]
	n[3] = to[3]
}

var nyb2byte = map[nyb]byte{
	{false, false, false, false}: 0,
	{false, false, false, true}:  1,
	{false, false, true, false}:  2,
	{false, false, true, true}:   3,
	{false, true, false, false}:  4,
	{false, true, false, true}:   5,
	{false, true, true, false}:   6,
	{false, true, true, true}:    7,
	{true, false, false, false}:  8,
	{true, false, false, true}:   9,
	{true, false, true, false}:   10,
	{true, false, true, true}:    11,
	{true, true, false, false}:   12,
	{true, true, false, true}:    13,
	{true, true, true, false}:    14,
	{true, true, true, true}:     15,
}

var byte2nyb = map[byte]nyb{
	0:  {false, false, false, false},
	1:  {false, false, false, true},
	2:  {false, false, true, false},
	3:  {false, false, true, true},
	4:  {false, true, false, false},
	5:  {false, true, false, true},
	6:  {false, true, true, false},
	7:  {false, true, true, true},
	8:  {true, false, false, false},
	9:  {true, false, false, true},
	10: {true, false, true, false},
	11: {true, false, true, true},
	12: {true, true, false, false},
	13: {true, true, false, true},
	14: {true, true, true, false},
	15: {true, true, true, true},
}

func populateByte2Nyb() {
	for i := byte(16); i > 0; i += 1 {
		byte2nyb[i] = byte2nyb[i%16]
	}
}

var (
	ram [2048]nyb
	// vram        [2048]nyb
	rom         [49152 * 2]instr
	cache       [16]ptr
	cacheOffset ptr

	a  nyb // output register
	b  nyb // output register
	t0 nyb
	t1 nyb

	carry bool
	pc    ptr
)

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

// atomic functions

func add() {
	sum := nyb2byte[a] + nyb2byte[b]
	if carry {
		sum += 1
	}

	if sum >= 16 {
		carry = true
	} else {
		carry = false
	}

	(&a).set(byte2nyb[sum])
}

// func mul() {
// 	product := nyb2byte[a] * nyb2byte[b]

// 	(&a).set(byte2nyb[product])
// 	(&b).set(byte2nyb[product/16])
// }

// func div() {
// 	ratio := nyb2byte[a] / nyb2byte[b]

// 	(&a).set(byte2nyb[ratio])
// }

func ror() {
	(&t0).set(a)
	a[0] = t0[1]
	a[1] = t0[2]
	a[2] = t0[3]
	a[3] = t0[0]
}

func rol() {
	(&t0).set(a)
	a[0] = t0[3]
	a[1] = t0[0]
	a[2] = t0[1]
	a[3] = t0[2]
}

func jcu(line ptr) {
	if !carry {
		pc = line
	}
}

func and() {
	a[0] = a[0] && b[0]
	a[1] = a[1] && b[1]
	a[2] = a[2] && b[2]
	a[3] = a[3] && b[3]
}

func or() {
	a[0] = a[0] || b[0]
	a[1] = a[1] || b[1]
	a[2] = a[2] || b[2]
	a[3] = a[3] || b[3]
}

func xor() {
	a[0] = a[0] != b[0]
	a[1] = a[1] != b[1]
	a[2] = a[2] != b[2]
	a[3] = a[3] != b[3]
}

func lib(n nyb) {
	(&b).set(n)
}

func ld() {
	(&a).set(ram[cache[nyb2byte[b]]+cacheOffset])
}

func ldb(ref nyb) {
	(&b).set(ram[cache[nyb2byte[ref]]+cacheOffset])
}

func st() {
	(&ram[cache[nyb2byte[b]]+cacheOffset]).set(a)
}

func sta(ref nyb) {
	(&ram[cache[nyb2byte[ref]]+cacheOffset]).set(a)
}

func swr(n nyb) {
	sum := 0
	if n[3] {
		sum += 1
	}
	if n[2] {
		sum += 2
	}
	if n[1] {
		sum += 1
	}
	if n[0] {
		sum += 2
	}

	if sum == 1 {
		(&t1).set(a)
		(&a).set(b)
		(&b).set(t1)
	} else if sum == 2 {
		(&t1).set(a)
		(&a).set(t0)
		(&t0).set(t1)
	} else if sum == 3 {
		(&t1).set(b)
		(&b).set(t1)
		(&t0).set(t1)
	} else {
		// TODO: more input modes?
	}
}

func lc() {
	carry = b[3]
}

func mca(index nyb, addr ptr) {
	cache[nyb2byte[index]] = addr
}

func mri(offset ptr) {
	cacheOffset = offset
}

// composite functions

func sub() {
	swp()
	neg()
	add()
}

func addi(n nyb) {
	lib(n)
	add()
}

func subi(n nyb) {
	lib(n)
	sub()
}

func neg() {
	not()
	addi(byte2nyb[1])
}

func rsh() {
	ror()
	lib(byte2nyb[7])
	and()
}

func lsh() {
	rol()
	lib(byte2nyb[14])
	and()
}

func not() {
	lib(byte2nyb[15])
	xor()
}

func andi(n nyb) {
	lib(n)
	and()
}

func ori(n nyb) {
	lib(n)
	or()
}

func xori(n nyb) {
	lib(n)
	xor()
}

func lia(n nyb) {
	swp()
	lib(n)
	swp()
}

func lda(ref nyb) {
	swp()
	ldb(ref)
	swp()
}

func ulia(n nyb) {
	lib(n)
	swp()
}

func ulda(ref nyb) {
	ldb(ref)
	swp()
}

func clr() {
	andi(byte2nyb[0])
}

func clc() {
	lib(byte2nyb[0])
	lc()
}

func stc() {
	lib(byte2nyb[1])
	lc()
}

func swp() {
	swr(byte2nyb[1])
}

func swa() {
	swr(byte2nyb[2])
}

func swb() {
	swr(byte2nyb[3])
}

func parseLine(line string) {
	parts := strings.Split(line, " ")

	var (
		arg1 nyb
		arg2 ptr
	)
	if len(parts) > 1 {
		arg1_64, _ := strconv.ParseInt(parts[1], 16, 16)
		arg1 = byte2nyb[byte(arg1_64)]
		if len(parts) > 2 {
			arg2_64, _ := strconv.ParseInt(parts[2], 16, 16)
			arg2 = ptr(arg2_64)
		}
	}

	// fmt.Println("parsed line:", parts[0], arg1, arg2)
	runInstruction(parts[0], arg1, arg2)
}

func runInstruction(instr string, arg1 nyb, arg2 ptr) {
	switch instr {
	case "add":
		add()
	// case "mul":
	// 	mul()
	// case "div":
	// 	div()
	case "ror":
		ror()
	case "rol":
		rol()
	case "jcu":
		jcu(arg2)
	case "and":
		and()
	case "or":
		or()
	case "xor":
		xor()
	case "lib":
		lib(arg1)
	case "ld":
		ld()
	case "ldb":
		ldb(arg1)
	case "lc":
		lc()
	case "st":
		st()
	case "sta":
		sta(arg1)
	case "mca":
		mca(arg1, ptr(arg2))
	case "mri":
		mri(arg2)
	case "swr":
		swr(arg1)
	case "sub":
		sub()
	case "subi":
		subi(arg1)
	case "addi":
		addi(arg1)
	case "neg":
		neg()
	case "rsh":
		rsh()
	case "lsh":
		lsh()
	case "not":
		not()
	case "andi":
		andi(arg1)
	case "ori":
		ori(arg1)
	case "xori":
		xori(arg1)
	case "lia":
		lia(arg1)
	case "lda":
		lda(arg1)
	case "ulia":
		ulia(arg1)
	case "ulda":
		ulda(arg1)
	// case "sti":
	// 	sti(byte2nyb[byte(arg1)])
	case "stc":
		stc()
	case "clc":
		clc()
	case "clr":
		clr()
	case "swp":
		swp()
	case "swa":
		swa()
	case "swb":
		swb()
	}

	if carry || instr != "jcu" {
		pc += 1
	}
}

var opMap = [16]string{
	"add",
	"ror",
	"rol",
	"jcu",
	"and",
	"or",
	"xor",
	"lib",
	"ld",
	"ldb",
	"st",
	"sta",
	"swr",
	"lc",
	"mca",
	"mri",
}

var (
	tillNextOp = 0
	opStr      = ""
	nybArg     = byte2nyb[0]
	ptrArg     = ptr(0)
)

func decode(i instr) string {
	op := byte(i / 16)
	arg := byte(i % 16)

	opStr := opMap[op]

	if tillNextOp == 0 {
		return fmt.Sprintf("%s %x", opStr, arg)
	} else {
		return fmt.Sprintf("%x %x", op, arg)
	}
}

func cycle() {
	i := rom[pc]
	op := byte(i / 16)
	arg := byte(i % 16)

	if tillNextOp == 0 {
		opStr = opMap[op]
		nybArg = byte2nyb[arg]

		if opStr == "jcu" || opStr == "mca" || opStr == "mri" {
			tillNextOp = 2
			pc += 1
		} else {
			runInstruction(opStr, nybArg, ptrArg)
		}

	} else {
		if tillNextOp == 2 {
			tillNextOp -= 1
			ptrArg = ptr(16 * i)
			pc += 1
		} else if tillNextOp == 1 {
			tillNextOp -= 1
			ptrArg += ptr(i)
			runInstruction(opStr, nybArg, ptrArg)
		}
	}
}

const DEBUG_MODE = false

func RunProgramString() {
	fmt.Println("Loading test program:")
	file, _ := os.Open("./build/fib.4smproc")
	defer file.Close()

	reader := bufio.NewReader(file)
	program := make([]string, 0)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		program = append(program, string(line))
		fmt.Println(string(line))
	}

	fmt.Println("Running test program...")

	for int(pc) < len(program) {
		if DEBUG_MODE {
			fmt.Println("after line:", program[pc])
		}

		parseLine(program[pc])

		if DEBUG_MODE {
			printRegs()
			printCache()
			printRam(0, 32)
		}
	}

	fmt.Println("Final state of the machine:")
	printRegs()
	printCache()
	printRam(0, 32)
	fmt.Printf("PC: %d\nCarry: %v\nCache offset: %d\n", pc, carry, cacheOffset)
}

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
		if DEBUG_MODE {
			fmt.Printf("after line: %s\n", decode(rom[pc]))
		}

		cycle()

		if DEBUG_MODE {
			printRegs()
			printCache()
			printRam(0, 32)
		}
	}

	fmt.Println("Final state of the machine:")
	printRegs()
	printCache()
	printRam(0, 32)
	fmt.Printf("PC: %d\nCarry: %v\nCache offset: %d\n", pc, carry, cacheOffset)

}
