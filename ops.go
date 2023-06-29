package main

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
	(&a).set(ram[cache[nyb2byte[b]]])
}

func ldb(ref nyb) {
	(&b).set(ram[cache[nyb2byte[ref]]])
}

func st() {
	(&ram[cache[nyb2byte[b]]]).set(a)
}

func sta(ref nyb) {
	(&ram[cache[nyb2byte[ref]]]).set(a)
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

func mri() {
	offset := 8 - int(nyb2byte[b])
	for i, addr := range cache {
		cache[i] = ptr(int(addr) - offset)
	}
}
