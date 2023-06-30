package cpu

func add() {
	sum := Nyb2Byte[a] + Nyb2Byte[b]
	if carry {
		sum += 1
	}

	if sum >= 16 {
		carry = true
	} else {
		carry = false
	}

	(&a).Set(Byte2Nyb[sum])
}

func ror() {
	(&t0).Set(a)
	a[0] = t0[1]
	a[1] = t0[2]
	a[2] = t0[3]
	a[3] = t0[0]
}

func rol() {
	(&t0).Set(a)
	a[0] = t0[3]
	a[1] = t0[0]
	a[2] = t0[1]
	a[3] = t0[2]
}

func jcu(line Ptr) {
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

func lib(n Nyb) {
	(&b).Set(n)
}

func ld() {
	(&a).Set(RAM[MemCache[Nyb2Byte[b]]])
}

func ldb(ref Nyb) {
	(&b).Set(RAM[MemCache[Nyb2Byte[ref]]])
}

func st() {
	(&RAM[MemCache[Nyb2Byte[b]]]).Set(a)
}

func sta(ref Nyb) {
	(&RAM[MemCache[Nyb2Byte[ref]]]).Set(a)
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
		(&t1).Set(a)
		(&a).Set(b)
		(&b).Set(t1)
	} else if sum == 2 {
		(&t1).Set(a)
		(&a).Set(t0)
		(&t0).Set(t1)
	} else if sum == 3 {
		(&t1).Set(b)
		(&b).Set(t1)
		(&t0).Set(t1)
	}
	// TODO: more input modes?
}

func lc() {
	carry = b[3]
}

func mca(index Nyb, addr Ptr) {
	MemCache[Nyb2Byte[index]] = addr
}

func mri() {
	offset := 8 - int(Nyb2Byte[b])
	for i, addr := range MemCache {
		MemCache[i] = Ptr(int(addr) - offset)
	}
}
