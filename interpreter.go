package main

type opName string

const (
	opUndef           opName = ""
	opAdd             opName = "add"
	opRotateRight     opName = "ror"
	opRotateLeft      opName = "rol"
	opJump            opName = "jcu"
	opAnd             opName = "and"
	opOr              opName = "or"
	opXor             opName = "xor"
	opLoadImmediateB  opName = "lib"
	opLoad            opName = "ld"
	opLoadB           opName = "ldb"
	opLoadCarry       opName = "lc"
	opStore           opName = "st"
	opStoreA          opName = "sta"
	opMemoryCacheAdd  opName = "mca"
	opLoadCacheOffset opName = "mci"
	opSwapRegisters   opName = "swr"
)

func runInstruction(instr opName, arg1 nyb, arg2 ptr) {
	switch instr {
	case opAdd:
		add()
	case opRotateRight:
		ror()
	case opRotateLeft:
		rol()
	case opJump:
		jcu(arg2)
	case opAnd:
		and()
	case opOr:
		or()
	case opXor:
		xor()
	case opLoadImmediateB:
		lib(arg1)
	case opLoad:
		ld()
	case opLoadB:
		ldb(arg1)
	case opLoadCarry:
		lc()
	case opStore:
		st()
	case opStoreA:
		sta(arg1)
	case opMemoryCacheAdd:
		mca(arg1, ptr(arg2))
	case opLoadCacheOffset:
		mri(arg2)
	case opSwapRegisters:
		swr(arg1)
	}

	// inc the program counter after every instruction, unless the program just jumped
	if carry || instr != "jcu" {
		pc += 1
	}
}

var opMap = [16]opName{
	opAdd,
	opRotateRight,
	opRotateLeft,
	opJump,
	opAnd,
	opOr,
	opXor,
	opLoadImmediateB,
	opLoad,
	opLoadB,
	opStore,
	opStoreA,
	opSwapRegisters,
	opLoadCarry,
	opMemoryCacheAdd,
	opLoadCacheOffset,
}

var (
	tillNextOp = 0
	opStr      = opUndef
	nybArg     = byte2nyb[0]
	ptrArg     = ptr(0)
)

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
