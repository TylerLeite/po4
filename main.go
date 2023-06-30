package main

import (
	"flag"
	"fmt"

	"github.com/TylerLeite/po4/cpu"
)

func main() {
	cpu.PowerOn()

	var program string
	flag.StringVar(&program, "load", "fib", "load this program from build/<arg>.bin")
	flag.StringVar(&program, "l", "fib", "load this program from build/<arg>.bin")
	flag.Parse()

	filename := fmt.Sprintf("./build/%s.bin", program)

	programSize := cpu.Load(filename)
	cpu.Run(programSize)
	cpu.PowerOff()
}
