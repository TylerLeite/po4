pass_1 = []

with open("programs/fib.4sm", "r") as file:
    for line in file:
        line = line.strip()
        if not line or line[0] == ";":
            continue
        else:
            line = line.split(";")[0]
            pass_1.append(line)

# now need to expand computed ops into atomics
def expand_computed(pass_n):
    pass_n_plus_1 = []
    for line in pass_n:
        parts = line.split(" ")
        if parts[0] == "sub":
            pass_n_plus_1.append("swp")
            pass_n_plus_1.append("neg")
            pass_n_plus_1.append("add")
        elif parts[0] == "addi":
            pass_n_plus_1.append("lib " + parts[1])
            pass_n_plus_1.append("add")
        elif parts[0] == "subi":
            pass_n_plus_1.append("lib " + parts[1])
            pass_n_plus_1.append("sub")
        elif parts[0] == "neg":
            pass_n_plus_1.append("not")
            pass_n_plus_1.append("addi 1")
        elif parts[0] == "rsh":
            pass_n_plus_1.append("ror")
            pass_n_plus_1.append("lib 7")
            pass_n_plus_1.append("and")
        elif parts[0] == "lsh":
            pass_n_plus_1.append("rol")
            pass_n_plus_1.append("lib 14")
            pass_n_plus_1.append("and")
        elif parts[0] == "not":
            pass_n_plus_1.append("lib 15")
            pass_n_plus_1.append("xor")
        elif parts[0] == "andi":
            pass_n_plus_1.append("lib " + parts[1])
            pass_n_plus_1.append("and")
        elif parts[0] == "ori":
            pass_n_plus_1.append("lib " + parts[1])
            pass_n_plus_1.append("or")
        elif parts[0] == "xori":
            pass_n_plus_1.append("lib " + parts[1])
            pass_n_plus_1.append("xor")
        elif parts[0] == "lia":
            pass_n_plus_1.append("swp")
            pass_n_plus_1.append("lib " + parts[1])
            pass_n_plus_1.append("swp")
        elif parts[0] == "lda":
            pass_n_plus_1.append("swp")
            pass_n_plus_1.append("ldb " + parts[1])
            pass_n_plus_1.append("swp")
        elif parts[0] == "ulia":
            pass_n_plus_1.append("lib " + parts[1])
            pass_n_plus_1.append("swp")
        elif parts[0] == "ulda":
            pass_n_plus_1.append("ldb " + parts[1])
            pass_n_plus_1.append("swp")
        elif parts[0] == "clr":
            pass_n_plus_1.append("andi 0")
        elif parts[0] == "clc":
            pass_n_plus_1.append("lib 0")
            pass_n_plus_1.append("lc")
        elif parts[0] == "stc":
            pass_n_plus_1.append("lib 1")
            pass_n_plus_1.append("lc")
        elif parts[0] == "swp":
            pass_n_plus_1.append("swr 1")
        elif parts[0] == "swa":
            pass_n_plus_1.append("swr 2")
        elif parts[0] == "swb":
            pass_n_plus_1.append("swr 3")
        else:
            pass_n_plus_1.append(line)
    return pass_n_plus_1

# maximum computed nesting is 3 (sub -> neg -> addi -> add), run expand 3 times
pass_1 = expand_computed(pass_1)
pass_1 = expand_computed(pass_1)
pass_1 = expand_computed(pass_1)

with open("build/fib.x.4sm", "w") as file:
    file.write("\n".join(pass_1))

def ptr_to_2_nybs(ptr):
    a = ptr//256
    a0 = a//16
    a1 = a%16

    b = ptr%256
    b0 = b//16
    b1 = b%16

    return f"{a0:x} {a1:x}", f"{b0:x} {b1:x}"

# also expand ptr ops into 3 lines
def expand_large_ops(pass_n):
    pass_n_plus_1 = []
    for line in pass_n:
        parts = line.split(" ")
        if parts[0] == "jcu":
            pass_n_plus_1.append(line)
            # l2, l3 = ptr_to_2_nybs(int(parts[1], 16))
            pass_n_plus_1.append("?")
            pass_n_plus_1.append("?")
        elif parts[0] == "mca":
            pass_n_plus_1.append("mca " + parts[1])
            l2, l3 = ptr_to_2_nybs(int(parts[2], 16))
            pass_n_plus_1.append(l2)
            pass_n_plus_1.append(l3)
        elif parts[0] == "mri":
            pass_n_plus_1.append("mri")
            l2, l3 = ptr_to_2_nybs(int(parts[1], 16))
            pass_n_plus_1.append(l2)
            pass_n_plus_1.append(l3)
        else:
            pass_n_plus_1.append(line)
    return pass_n_plus_1


pass_1 = expand_large_ops(pass_1)
pass_2 = []
labels = {}

def my_hex(n):
    out = hex(i)[2:]
    while len(out) < 4:
        out = "0" + out
    return out

i = 0
# build label-to-line map
for line in pass_1:
    # NOTE: two labels in consecutive lines will mess everything up
    if ":" in line:
        labels[line[:-1]] = my_hex(i)
    else:
        i += 1
        pass_2.append(line)

# expand labels
pass_3 = []
for line in pass_2:
    for label in labels:
        if label in line:
            line = line.replace(label, labels[label])
    pass_3.append(line.strip())

# expand jumps
for i, line in enumerate(pass_3):
    parts = line.split(" ")
    if parts[0] == "jcu":
        l2, l3 = ptr_to_2_nybs(int(parts[1], 16))
        pass_3[i] = "jcu"
        pass_3[i+1] = l2
        pass_3[i+2] = l3

opMap = {
	"add": "0",
	"ror": "1",
	"rol": "2",
	"jcu": "3",
	"and": "4",
	"or": "5",
	"xor": "6",
	"lib": "7",
	"ld": "8",
	"ldb": "9",
	"st": "A",
	"sta": "B",
	"swr": "C",
	"lc": "D",
	"mca": "E",
	"mri": "F",
    "0": "0",
    "1": "1",
    "2": "2",
    "3": "3",
    "4": "4",
    "5": "5",
    "6": "6",
    "7": "7",
    "8": "8",
    "9": "9",
    "A": "A",
    "a": "A",
    "B": "B",
    "b": "B",
    "C": "C",
    "c": "C",
    "D": "D",
    "d": "D",
    "E": "E",
    "e": "E",
    "F": "F",
    "f": "F"
}

pass_4 = []
for line in pass_3:
    parts = line.split(" ")
    if len(parts) == 1:
        parts = [parts[0], "0"]
    pass_4.append("".join([opMap[part] for part in parts]))

with open("build/fib.4bc", "w") as file:
    hexstr = ""
    for i in range(0, len(pass_4), 8):
        hexstr += " ".join(pass_4[i:i+8]) + "\n"
    # hexstr = "\n".join(pass_4)
    file.write(hexstr)

with open("build/fib.bin", "wb") as file:
    hexstr = "".join(pass_4)
    byts = bytearray.fromhex(hexstr)
    file.write(byts)