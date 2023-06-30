package cpu

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

func runInstruction(instr opName, arg1 Nyb, arg2 Ptr) {
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
		mca(arg1, Ptr(arg2))
	case opLoadCacheOffset:
		mri()
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
	nybArg     = Byte2Nyb[0]
	ptrArg     = Ptr(0)
)

func cycle() {
	i := ROM[pc]
	op := byte(i / 16)
	arg := byte(i % 16)

	if tillNextOp == 0 {
		opStr = opMap[op]
		nybArg = Byte2Nyb[arg]

		if opStr == "jcu" || opStr == "mca" || opStr == "mri" {
			tillNextOp = 2
			pc += 1
		} else {
			runInstruction(opStr, nybArg, ptrArg)
		}

	} else {
		if tillNextOp == 2 {
			tillNextOp -= 1
			ptrArg = Ptr(16 * i)
			pc += 1
		} else if tillNextOp == 1 {
			tillNextOp -= 1
			ptrArg += Ptr(i)
			runInstruction(opStr, nybArg, ptrArg)
		}
	}
}
