// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/4/Fill.asm

// Runs an infinite loop that listens to the keyboard input. 
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel. When no key is pressed, 
// the screen should be cleared.


(LOOP)
@KBD
D=M
@WHITE
D;JEQ
//fill with black
@8192
D=A
(FILL_BLACK)
@LOOP
D=D-1;JLT
@SCREEN
A=D+A
M=-1
@FILL_BLACK
0;JMP

(WHITE)
@8192
D=A
(FILL_WHITE)
@LOOP
D=D-1;JLT
@SCREEN
A=D+A
M=0
@FILL_WHITE
0;JMP

