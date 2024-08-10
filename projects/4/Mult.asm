// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/4/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
// The algorithm is based on repetitive addition.

@R0
D=M
@R1
D=D-M
@SMALLER
D;JLE //if R0 <= R1
@R1
D=M
@small
M=D
@R0
D=M
@great
M=D
@END_BRANCH
0;JMP
(SMALLER)
@R0
D=M
@small
M=D
@R1
D=M
@great
M=D
(END_BRANCH)
@R2
M=0

(LOOP)
@small
M=M-1
D=M
@EXIT
D;JLT
@great
D=M
@R2
M=D+M
@LOOP
0;JMP

(EXIT)


