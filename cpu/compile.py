import sys

if len(sys.argv) < 2:
    print("No input file provided")
    sys.exit(0)

program = sys.argv[1]
code = []

with open(f"programs/{program}.4sm", "r") as file:
    for line in file:
        code.append(line)

# Remove comments, empty lines
def strip_comments(lines):
    new_lines = []
    for line in lines:
        line = line.strip()
        if not line or line[0] == ";":
            continue
        else:
            line = line.split(";")[0]
            new_lines.append(line)
    return new_lines


# TODO: support multiple constants in a single line e.g. adding constants together
def replace_constants(lines):
    def get_constants(lines):
        constants = {}
        new_lines = []
        for line in lines:
            if line[0] == "$":
                parts = line.split(" ")
                constants[parts[0]] = int(parts[-1], 16)
            else:
                new_lines.append(line)
        return constants, new_lines
    
    constants, lines = get_constants(lines)
    new_lines = []
    for line in lines:
        new_line = line
        for name, value in constants.items():
            if name in line:
                if "+" in line:
                    value += int(line.split('+')[-1], 16)
                    new_line = line.split('+')[0]
                elif "-" in line:
                    value -= int(line.split('-')[-1], 16)
                    new_line = line.split('-')[0]

                new_line = new_line.replace(name, my_hex(value))
        new_lines.append(new_line)
    return new_lines

# TODO: support aliasing (mostly to simulate var names for memory cache address)
# e.g. .x 1 ... lda .x -> lda 1
def replace_aliases(lines):
    return lines            

# Some optimizations are easier with composite functions still un-expanded
# e.g. replacing safe composite functions with unsafe versions when B will be loaded anyway
def optimize_1(lines):
    return lines #lol

# Expand computed ops into atomics
def expand_computed(lines):
    new_lines = []
    for line in lines:
        parts = line.split(" ")
        if parts[0] == "sub":
            new_lines.append("swb")
            new_lines.append("swp")
            new_lines.append("neg")
            new_lines.append("add")
            new_lines.append("swb")
        elif parts[0] == "usub":
            new_lines.append("swp")
            new_lines.append("neg")
            new_lines.append("add")
        elif parts[0] == "addi":
            new_lines.append("swb")
            new_lines.append("lib " + parts[1])
            new_lines.append("add")
            new_lines.append("swb")
        elif parts[0] == "uaddi":
            new_lines.append("lib " + parts[1])
            new_lines.append("add")
        elif parts[0] == "subi":
            new_lines.append("swb")
            new_lines.append("lib " + parts[1])
            new_lines.append("usub")
            new_lines.append("swb")
        elif parts[0] == "usubi":
            new_lines.append("lib " + parts[1])
            new_lines.append("usub")
        elif parts[0] == "neg":
            new_lines.append("not")
            new_lines.append("addi 1")
        elif parts[0] == "uneg":
            new_lines.append("unot")
            new_lines.append("uaddi 1")
        elif parts[0] == "rsh":
            new_lines.append("ror")
            new_lines.append("swb")
            new_lines.append("lib 7")
            new_lines.append("and")
            new_lines.append("swb")
        elif parts[0] == "ursh":
            new_lines.append("ror")
            new_lines.append("lib 7")
            new_lines.append("and")
        elif parts[0] == "lsh":
            new_lines.append("rol")
            new_lines.append("swb")
            new_lines.append("lib E")
            new_lines.append("and")
            new_lines.append("swb")
        elif parts[0] == "ulsh":
            new_lines.append("rol")
            new_lines.append("lib E")
            new_lines.append("and")
        elif parts[0] == "not":
            new_lines.append("swb")
            new_lines.append("lib F")
            new_lines.append("xor")
            new_lines.append("swb")
        elif parts[0] == "unot":
            new_lines.append("lib F")
            new_lines.append("xor")
        elif parts[0] == "andi":
            new_lines.append("swb")
            new_lines.append("lib " + parts[1])
            new_lines.append("and")
            new_lines.append("swb")
        elif parts[0] == "uandi":
            new_lines.append("lib " + parts[1])
            new_lines.append("and")
        elif parts[0] == "ori":
            new_lines.append("swb")
            new_lines.append("lib " + parts[1])
            new_lines.append("or")
            new_lines.append("swb")
        elif parts[0] == "uori":
            new_lines.append("lib " + parts[1])
            new_lines.append("or")
        elif parts[0] == "xori":
            new_lines.append("swb")
            new_lines.append("lib " + parts[1])
            new_lines.append("xor")
            new_lines.append("swb")
        elif parts[0] == "uxori":
            new_lines.append("lib " + parts[1])
            new_lines.append("xor")
        elif parts[0] == "lia":
            new_lines.append("swp")
            new_lines.append("lib " + parts[1])
            new_lines.append("swp")
        elif parts[0] == "ulia":
            new_lines.append("lib " + parts[1])
            new_lines.append("swp")
        elif parts[0] == "lda":
            new_lines.append("swp")
            new_lines.append("ldb " + parts[1])
            new_lines.append("swp")
        elif parts[0] == "ulda":
            new_lines.append("ldb " + parts[1])
            new_lines.append("swp")
        elif parts[0] == "clr":
            new_lines.append("andi 0")
        elif parts[0] == "clc":
            new_lines.append("lib 0")
            new_lines.append("lc")
        elif parts[0] == "stc":
            new_lines.append("lib 1")
            new_lines.append("lc")
        elif parts[0] == "swp":
            new_lines.append("swr 1")
        elif parts[0] == "swa":
            new_lines.append("swr 2")
        elif parts[0] == "swb":
            new_lines.append("swr 3")
        else:
            new_lines.append(line)
    return new_lines

# Optimize again, expanding prcess is dumb and sometimes produces redundant code
def optimize_2(lines):
    return lines #lol

# Expand ptr ops, add placeholders for jumps
def expand_large_ops(lines):
    new_lines = []
    for line in lines:
        parts = line.split(" ")
        if parts[0] == "jcu":
            new_lines.append(line)
            new_lines.append("?")
            new_lines.append("?")
        elif parts[0] == "mca":
            new_lines.append("mca " + parts[1])
            l2, l3 = ptr_to_2_nybs(int(parts[2], 16))
            new_lines.append(l2)
            new_lines.append(l3)
        elif parts[0] == "mri":
            new_lines.append("mri")
            l2, l3 = ptr_to_2_nybs(int(parts[1], 16))
            new_lines.append(l2)
            new_lines.append(l3)
        else:
            new_lines.append(line)
    return new_lines

# Convert string labels to rom addresses
def expand_labels(lines):
    # build label-to-line map
    def get_labels(lines):
        i = 0
        labels = {}
        new_lines = []
        for line in lines:
            # NOTE: two labels in consecutive lines will mess everything up
            if ":" in line:
                labels[line[:-1]] = my_hex(i)
            else:
                i += 1
                new_lines.append(line)
        return labels, new_lines
    
    # use label-to-line map
    labels, lines = get_labels(lines)
    new_lines = []
    for line in lines:
        for label in labels:
            if label in line:
                line = line.replace(label, labels[label])
        new_lines.append(line.strip())
    return new_lines

# Place actual rom addresses in jump instructions
def expand_jumps(lines):
    for i, line in enumerate(lines):
        parts = line.split(" ")
        if parts[0] == "jcu":
            l2, l3 = ptr_to_2_nybs(int(parts[1], 16))
            lines[i] = "jcu"
            lines[i+1] = l2
            lines[i+2] = l3
    return lines

# Convert string op names to hex codes
def convert_to_hex(lines):
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

    new_lines = []
    for line in lines:
        parts = line.split(" ")
        if len(parts) == 1:
            parts = [parts[0], "0"]
        new_lines.append("".join([opMap[part] for part in parts]))
    return new_lines

# Util functions
def ptr_to_2_nybs(ptr):
    a = ptr//256
    a0 = a//16
    a1 = a%16

    b = ptr%256
    b0 = b//16
    b1 = b%16

    return f"{a0:x} {a1:x}", f"{b0:x} {b1:x}"

def my_hex(n):
    out = hex(n)[2:]
    while len(out) < 4:
        out = "0" + out
    return out


code = strip_comments(code)
code = replace_constants(code)
code = optimize_1(code)
# maximum computed nesting is 3 (sub -> neg -> addi -> add), run expand 3 times
code = expand_computed(code)
code = expand_computed(code)
code = expand_computed(code)
code = optimize_2(code)

with open(f"build/{program}.x.4sm", "w") as file:
    file.write("\n".join(code))

code = expand_large_ops(code)
code = expand_labels(code)
expand_jumps(code)
code = convert_to_hex(code)

with open(f"build/{program}.4bc", "w") as file:
    hexstr = ""
    for i in range(0, len(code), 8):
        hexstr += " ".join(code[i:i+8]) + "\n"
    file.write(hexstr)

with open(f"build/{program}.bin", "wb") as file:
    hexstr = "".join(code)
    byts = bytearray.fromhex(hexstr)
    file.write(byts)