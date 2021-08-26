#include "textflag.h"

GLOBL ·NameData(SB),NOPTR,$8
DATA  ·NameData(SB)/8,$"gopher"

GLOBL ·Name(SB),NOPTR,$16
DATA  ·Name+0(SB)/8,$·NameData(SB)
DATA  ·Name+8(SB)/8,$6

