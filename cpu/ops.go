package cpu

func add() {
	sum := Nyb2Byte[A] + Nyb2Byte[B]
	if Carry {
		sum += 1
	}

	if sum >= 16 {
		Carry = true
	} else {
		Carry = false
	}

	(&A).Set(Byte2Nyb[sum])
}

func mod() {
	rem := Nyb2Byte[A] % Nyb2Byte[B]
	(&A).Set(Byte2Nyb[rem])
}

func rol() {
	(&T0).Set(A)
	A[0] = T0[3]
	A[1] = T0[0]
	A[2] = T0[1]
	A[3] = T0[2]
}

func jcu(line Ptr) {
	if !Carry {
		ProgramCounter = line
	}
}

func fnc(line Ptr) {
	CallStack = append(CallStack, ProgramCounter)
	ProgramCounter = line
}

func ret() {
	i := len(CallStack) - 1
	ProgramCounter = CallStack[i]
	CallStack = CallStack[:i]
}

func and() {
	A[0] = A[0] && B[0]
	A[1] = A[1] && B[1]
	A[2] = A[2] && B[2]
	A[3] = A[3] && B[3]
}

func or() {
	A[0] = A[0] || B[0]
	A[1] = A[1] || B[1]
	A[2] = A[2] || B[2]
	A[3] = A[3] || B[3]
}

func xor() {
	A[0] = A[0] != B[0]
	A[1] = A[1] != B[1]
	A[2] = A[2] != B[2]
	A[3] = A[3] != B[3]
}

func lib(n Nyb) {
	(&B).Set(n)
}

func rlb() {
	(&B).Set(RAM[MemCache[Nyb2Byte[B]]])
}

func st() {
	(&RAM[MemCache[Nyb2Byte[B]]]).Set(A)
}

func swr(n Nyb) {
	// This gives some unintuitive behavior like swr 0101 will be interpreted as swr 0010
	// but there's no reason to call swr 0101 anyway and this is a lot easier to implement in hardware
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
		(&T1).Set(A)
		(&A).Set(B)
		(&B).Set(T1)
	} else if sum == 2 {
		(&T1).Set(A)
		(&A).Set(T0)
		(&T0).Set(T1)
	} else if sum == 3 {
		(&T1).Set(B)
		(&B).Set(T0)
		(&T0).Set(T1)
	}
	// TODO: more input modes?
}

func lc() {
	Carry = B[3]
}

func mca(index Nyb, addr Ptr) {
	MemCache[Nyb2Byte[index]] = addr
}

func mri() {
	offset := 8 - int(Nyb2Byte[B])
	for i, addr := range MemCache {
		MemCache[i] = Ptr(int(addr) - offset)
	}
}
