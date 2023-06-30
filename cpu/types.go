package cpu

type Ptr int16
type Instr uint8
type Nyb [4]bool

func (n *Nyb) Set(to Nyb) {
	n[0] = to[0]
	n[1] = to[1]
	n[2] = to[2]
	n[3] = to[3]
}

var Nyb2Byte = map[Nyb]byte{
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

var Byte2Nyb = map[byte]Nyb{
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
