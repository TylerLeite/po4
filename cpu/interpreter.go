package cpu

import "fmt"

type opName string

const (
	opUndef           opName = ""
	opAdd             opName = "add"
	opMod             opName = "mod"
	opRotateLeft      opName = "rol"
	opJump            opName = "jcu"
	opFunctionCall    opName = "fnc"
	opReturn          opName = "return"
	opAnd             opName = "and"
	opOr              opName = "or"
	opXor             opName = "xor"
	opLoadImmediateB  opName = "lib"
	opReloadB         opName = "rlb"
	opLoadCarry       opName = "lc"
	opStore           opName = "st"
	opSwapRegisters   opName = "swr"
	opMemoryCacheAdd  opName = "mca"
	opLoadCacheOffset opName = "mci"
)

func runInstruction(instr opName, arg1 Nyb, arg2 Ptr) {
	switch instr {
	case opAdd:
		add()
	case opMod:
		mod()
	case opRotateLeft:
		rol()
	case opJump:
		jcu(arg2)
	case opFunctionCall:
		fnc(arg2)
	case opReturn:
		ret()
	case opAnd:
		and()
	case opOr:
		or()
	case opXor:
		xor()
	case opLoadImmediateB:
		lib(arg1)
	case opReloadB:
		rlb()
	case opLoadCarry:
		lc()
	case opStore:
		st()
	case opMemoryCacheAdd:
		mca(arg1, Ptr(arg2))
	case opLoadCacheOffset:
		mri()
	case opSwapRegisters:
		swr(arg1)
	}

	PrintMemCache()
	PrintRegisters()
	PrintRAM(0, 16, 8)

	// inc the program counter after every instruction, unless the program just jumped
	if instr != opFunctionCall && (Carry || instr != opJump) {
		ProgramCounter += 1
	}
}

var opMap = [16]opName{
	opAdd,
	opMod,
	opRotateLeft,
	opJump,
	opFunctionCall,
	opReturn,
	opAnd,
	opOr,
	opXor,
	opLoadImmediateB,
	opReloadB,
	opStore,
	opSwapRegisters,
	opLoadCarry,
	opMemoryCacheAdd,
	opLoadCacheOffset,
}

var (
	tillNextOp = 0
	opStr      = opUndef
	nybArg     = Byte2Nyb[0]
	ptrArg     = Ptr(0)
)

func cycle() {
	i := ROM[ProgramCounter]
	op := byte(i / 16)
	arg := byte(i % 16)

	fmt.Println(PrintInstruction(i))

	if tillNextOp == 0 {
		opStr = opMap[op]
		nybArg = Byte2Nyb[arg]

		if opStr == opJump || opStr == opMemoryCacheAdd || opStr == opLoadCacheOffset || opStr == opFunctionCall {
			tillNextOp = 2
			ProgramCounter += 1
		} else {
			runInstruction(opStr, nybArg, ptrArg)
		}
	} else {
		if tillNextOp == 2 {
			tillNextOp -= 1
			ptrArg = Ptr(16 * i)
			ProgramCounter += 1
		} else if tillNextOp == 1 {
			tillNextOp -= 1
			ptrArg += Ptr(i)
			runInstruction(opStr, nybArg, ptrArg)
		}
	}
}
