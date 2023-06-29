## Structure of the CPU

* 4x 4-bit registers: A, B, T1, T2
    - Operations generally use A and B, store results in A
    - Sometimes operations use T1 and T2 so B doesn't get disrupted
* 1-bit carry
* 16-bit program counter
    - Incremented after every non-jumping instruction
    - Might end up incrementing after every instruction and jumping to the spot before you want to land, depending on hardware implementation
* 16x16-bit memory cache
    - RAM and ROM are both more than 16-bits long, so you need addresses that don't fit in a register to index into them. The solution is to store longer addresses in a cache that is accessed by a 4-bit pointer-pointer.
    - A cache offset is also supported so you can index into structured data without having to modify the entire memory cache. Right now the cache offset is 16-bit. Might lower it to 4-bit if that seems good enough
    - Loading new addresses into the memory cache takes 3 cycles per value
    - Loading a new cache offset takes 3 cycles
* 2048x4-bit RAM
* 2048x4-bit VRAM
    - Used for the VPU, not yet implemented
* 49152x8-bit ROM
    - Each instruction is 8 bits wide (4 op bits + 4 arg bits)
    - For multi-cycle instructions, the extra instructions are read as halves of a 16-bit arg, rather than being split up into op and arg

## Operations
A list of all supported operations available in ops.txt
- All atomic operations preserve the state of B unless the explicit purpose is to modify B (e.g. ldb)
- By default, composite operations are safe and will preserve the state of B. Unsafe versions are also provided. They expand into fewer operations and should be preferred if you don't care what's in B
