#include "textflag.h"

GLOBL text<>(SB),NOPTR,$16
DATA text<>+0(SB)/8,$"Hello Wo"
DATA text<>+8(SB)/8,$"rld!"

GLOBL ·Helloworld(SB),NOPTR,$24            // var helloworld []byte("Hello World!")
DATA ·Helloworld+0(SB)/8,$text<>(SB) // StringHeader.Data
DATA ·Helloworld+8(SB)/8,$12         // StringHeader.Len
DATA ·Helloworld+16(SB)/8,$16        // StringHeader.Cap


