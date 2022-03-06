# 汇编
![汇编到机器码过程](.introduction_images/process.png)

汇编语言（assembly language):用于电子计算机、微处理器、微控制器或其他可编程器件的低级语言，亦称为符号语言

1. 助记符（Mnemonics）代替机器指令的操作码
2. 用地址符号（Symbol）或标号（Label）代替指令或操作数的地址

## 汇编器和链接器：

- 汇编器（assembler）是一种工具程序，用于将汇编语言源程序转换为机器语言
  1. windows: Microsoft 宏汇编器（称为 MASM）,TASM（Turbo 汇编器），NASM（Netwide 汇编器）和 MASM32（MASM 的一种变体）
  2. linux: GAS（GNU 汇编器）和 NASM,NASM 的语法与 MASM 的最相似
- 链接器（linker）把汇编器生成的单个文件 组合 为一个可执行程序。
- 调试器（debugger），使程序员可以在程序运行时，单步执行程序并检查寄存器和内存状态。

## 汇编语言与机器语言有什么关系:

- 机器语言（machine language）是一种数字语言， 专门设计成能被计算机处理器（CPU）理解。所有 x86 处理器都理解共同的机器语言。
- 汇编语言（assembly language）包含用短助记符如 ADD、MOV、SUB 和 CALL 书写的语句。汇编语言与机器语言是一对一（one-to-one）的关系：
    每一条汇编语言指令对应一条机器语言指令。寄存器（register）是 CPU 中被命名的存储位置，用于保存操作的中间结果。

## 优点
    
1. 因为用汇编语言设计的程序最终被转换成机器指令，故能够保持机器语言的一致性，直接、简捷，并能像机器指令一样访问、控制计算机的各种硬件设备，
如磁盘、存储器、CPU、I/O端口等。使用汇编语言，可以访问所有能够被访问的软、硬件资源。

2. 目标代码简短，占用内存少，执行速度快，是高效的程序设计语言，经常与高级语言配合使用，以改善程序的执行速度和效率，弥补高级语言在硬件控制方面的不足，应用十分广泛


## 缺点
    
1. 汇编语言是面向机器的，处于整个计算机语言层次结构的底层，故被视为一种低级语言，通常是为特定的计算机或系列计算机专门设计的。
不同的处理器有不同的汇编语言语法和编译器，编译的程序无法在不同的处理器上执行，缺乏可移植性；
2. 难于从汇编语言代码上理解程序设计意图，可维护性差，即使是完成简单的工作也需要大量的汇编语言代码，很容易产生bug，难于调试；
3. 使用汇编语言必须对某种处理器非常了解，而且只能针对特定的体系结构和处理器进行优化，开发效率很低，周期长且单调。

## 语言组成
    
1. 数据传送指令:通用数据传送指令MOV、条件传送指令CMOVcc、堆栈操作指令PUSH/PUSHA/PUSHAD/POP/POPA/POPAD、
交换指令XCHG/XLAT/BSWAP、地址或段描述符选择子传送指令LEA/LDS/LES/LFS/LGS/LSS等

2. 整数和逻辑运算指令:加法指令ADD/ADC、减法指令SUB/SBB、加一指令INC、减一指令DEC、比较操作指令CMP、
乘法指令MUL/IMUL、除法指令DIV/IDIV、符号扩展指令CBW/CWDE/CDQE、十进制调整指令DAA/DAS/AAA/AAS、逻辑运算指令NOT/AND/OR/XOR/TEST等。

3. 移位指令:将寄存器或内存操作数移动指定的次数。包括逻辑左移指令SHL、逻辑右移指令SHR、算术左移指令SAL、算术右移指令SAR、循环左移指令ROL、循环右移指令ROR等

4. 位操作指令:位测试指令BT、位测试并置位指令BTS、位测试并复位指令BTR、位测试并取反指令BTC、位向前扫描指令BSF、位向后扫描指令BSR等

5. 条件设置指令: 这不是一条具体的指令，而是一个指令簇，包括大约30条指令，用于根据EFLAGS寄存器的某些位状态来设置一个8位的寄存器或者内存操作数。
比如SETE/SETNE/SETGE等等

6. 控制转移指令:这部分包括无条件转移指令JMP、条件转移指令Jcc/JCXZ、循环指令LOOP/LOOPE/LOOPNE、过程调用指令CALL、子过程返回指令RET、
中断指令INTn、INT3、INTO、IRET等。注意，Jcc是一个指令簇，包含了很多指令，用于根据EFLAGS寄存器的某些位状态来决定是否转移；
INT n是软中断指令，n可以是0到255之间的数，用于指示中断向量号。

7. 串操作指令:这部分指令用于对数据串进行操作，包括串传送指令MOVS、串比较指令CMPS、串扫描指令SCANS、串加载指令LODS、串保存指令STOS，
这些指令可以有选择地使用REP/REPE/REPZ/REPNE和REPNZ的前缀以连续操作
![](.introduction_images/movsb.png)
si:源寄存器  di:目标寄存器
movsb -> move string byte(以字节单元传送): 将ds:si指向的内存单元中的字节 送入 es:di中，然后根据标志寄存器DF位的值，将si和di递增1或则递减1。
movsw -> move string word(以字单元传送): 将ds:si指向的内存单元中的字节 送入 es:di中，然后根据标志寄存器DF位的值，将si和di递增2或则递减2。
一般和rep搭配，如rep movsb，rep的作用是根据cx的值，重复执行后面的 串操作指令。

8. 输入输出指令: 这部分指令用于同外围设备交换数据，包括端口输入指令IN/INS、端口输出指令OUT/OUTS

## cpu可以直接读取数据的地方
1. cpu 内部寄存器
2. 内存单元
![img.png](chapter11_assembly_language/01asm/.introduction_images/cpu_read_mem.png)
3. 硬件端口
![img_1.png](chapter11_assembly_language/01asm/.introduction_images/cpu_read_port.png)

## cpu对存储器(内存)的读写操作过程

![](.introduction_images/asm_process.png)   
地址总线：64位cpu，代表查找能力

### 1. 数据总线:宽度决定与外界的数据传输速度
![](.introduction_images/address_bus.png)

### 2. 控制总线：读和写
![](.introduction_images/control_bus.png)

### 3. 数据总线
![](.introduction_images/data_address.png)


## 存储器
![](.introduction_images/rom_n_ram.png)  
- 从读写属性上分：随机存储器(ram)和只读存储器(rom).
- 功能上分：bios(BASIC INPUT/OUTPUT SYSTEM)基本输入输出系统上的rom.

## 寄存器
寄存器是中央处理器内的组成部分。寄存器是有限存贮容量的高速存贮部件，它们可用来暂存指令、数据和地址。

CPU 本身只负责运算，不负责储存数据。数据一般都储存在内存之中，CPU 要用的时候就去内存读写数据。但是，CPU的运算速度远高于内存的读写速度，
为了避免被拖慢，CPU 都自带一级缓存和二级缓存。基本上，CPU 缓存可以看作是读写速度较快的内存。

但是，CPU 缓存还是不够快，另外数据在缓存里面的地址是不固定的，CPU 每次读写都要寻址也会拖慢速度。
因此，除了缓存之外，CPU 还自带了寄存器（register），用来储存最常用的数据。也就是说，那些最频繁读写的数据（比如循环变量），都会放在寄存器里面，
CPU 优先读写寄存器，再由寄存器跟内存交换数据.

### 寄存器位数区分
```
ah/al => 8 位
ax/bx => 16 位
eax/ebx => 32 位
rax/rbx => 64 位，8 byte

rax : 0x0000000000000000

```

### 应用层代码一般只会用到如下分为三类的19个寄存器

1. 通用寄存器：rax, rbx, rcx, rdx, rsi, rdi, rbp, rsp, r8, r9, r10, r11, r12, r13, r14, r15寄存器。
CPU对这16个通用寄存器的用途没有做特殊规定，程序员和编译器可以自定义其用途（rsp/rbp寄存器其实是有特殊用途的), si和di是8086cpu与bx相近的寄存器，只是不能分为两个8位寄存器使用

rsp 栈顶寄存器和rbp栈基址寄存器:   
这两个寄存器都跟函数调用栈有关，其中rsp寄存器一般用来存放函数调用栈的栈顶地址，而rbp寄存器通常用来存放函数的栈帧起始地址，
编译器一般使用这两个寄存器加一定偏移的方式来访问函数局部变量或函数参数，比如
```shell
mov    0x8(%rsp),%rdx
```
这条指令把地址为 0x8(%rsp) 的内存中的值拷贝到rdx寄存器，这里的0x8(%rsp) 就利用了 rsp 寄存器加偏移 8 的方式来读取内存中的值.

2. 程序计数寄存器(PC寄存器，有时也叫IP寄存器,rip寄存器): 它用来存放下一条即将执行的指令的地址，这个寄存器决定了程序的执行流程.表示当前我的代码运行到哪里了.

rip寄存器里面存放的是CPU即将执行的下一条指令在内存中的地址。看如下汇编语言代码片段：
```shell
0x0000000000400770: add   %rdx,%rax
0x0000000000400773: mov   $0x0,%ecx
```

假设当前CPU正在执行第一条指令，这条指令在内存中的地址是0x0000000000400770，紧接它后面的下一条指令的地址是0x0000000000400773，所以此时rip寄存器里面存放的值是0x0000000000400773。
这里需要牢记的就是rip寄存器的值不是正在被CPU执行的指令在内存中的地址，而是紧挨这条正在被执行的指令后面那一条指令的地址.

在前面的两个汇编指令片段中并没有指令修改 rip寄存器的值，是怎么做到让它一直指向下一条即将执行的指令的呢？其实修改rip寄存器的值是CPU自动控制的，不需要我们用指令去修改，
当然CPU也提供了几条可以间接修改rip寄存器的指令.

3. 段寄存器：
- cs代码段寄存器:	存放当前正在运行的程序代码所在段的段基址，表示当前使用的指令代码可以从该段寄存器指定的存储器段中取得，相应的偏移量则由IP提供。
- ds数据段寄存器： 当前程序使用的数据所存放段的最低地址，即存放数据段的段基址
- ss是堆栈段寄存器： 当前堆栈的底部地址，即存放堆栈段的段基址
- es是扩展段寄存器： 当前程序使用附加数据段的段基址，该段是串操作指令中目的串所在的段
- fs标志段寄存器: fs是80386起增加的两个辅助段寄存器之一,在这之前只有一个辅助段寄存器ES,FS寄存器指向当前活动线程的TEB结构（线程结构)
- gs全局段寄存器。
fs,gs 一般用它来实现线程本地存储（TLS），比如AMD64 linux平台下go语言和pthread都使用fs寄存器来实现系统线程的TLS
```css
偏移  说明
000  指向SEH链指针
004  线程堆栈顶部
008  线程堆栈底部
00C  SubSystemTib
010  FiberData
014  ArbitraryUserPointer
018  FS段寄存器在内存中的镜像地址
020  进程PID
024  线程ID
02C  指向线程局部存储指针
030  PEB结构地址（进程结构）
034  上个错误号
```
上述这些寄存器除了fs和gs段寄存器是16位的，其它都是64位的，也就是8个字节，其中的16个通用寄存器还可以作为32/16/8位寄存器使用，
只是使用时需要换一个名字，比如可以用eax这个名字来表示一个32位的寄存器，它使用的是rax寄存器的低32位

Note:标志寄存器如 eflags:记录各种运算结果(是否为 0，是否发生溢出)的标志位。


## 汇编指令格式   
AT&T格式说明:

![](.introduction_images/direct_addr.png)
![](.introduction_images/indirect_addr.png)

1. 立即操作数需要加上$符号做前缀，如  "mov $0x1 %rdi" 这条指令中第一个操作数不是寄存器，也不是内存地址，而是直接写在指令中的一个常数，
这种操作数叫做立即操作数。这条指令表示把数值0x1放入rdi寄存器中

2. 寄存器间接寻址的格式为 offset(%register)，如果offset为0，则可以略去偏移不写直接写成(%register)。何为间接寻址呢？

其实就是指指令中的寄存器并不是真正的源操作数或目的操作数，寄存器的值是一个内存地址，这个地址对应的内存才是真正的源或目的操作数，比如 mov %rax, (%rsp)这条指令，
第二个操作数(%rsp)中的寄存器的名字用括号括起来了，表示间接寻址，rsp的值是一个内存地址，这条指令的真实意图是把rax寄存器中的值赋值给rsp寄存器的值（内存地址）对应的内存，
rsp寄存器本身的值不会被修改，

作为比较，我们看一下 mov%rax, %rsp 这条指令 ，这里第二个操作数仅仅少了个括号，变成了直接寻址，意思完全不一样了，
这条指令的意思是把rax的值赋给rsp，这样rsp寄存器的值被修改为跟rax寄存器一样的值了。上面面的2张图展示了这两种寻址方式的不同

3. 与内存相关的一些指令的操作码会加上b, w, l和q字母分别表示操作的内存是1，2，4还是8个字节，比如指令 movl $0x0,-0x8(%rbp) ，
这条指令操作码movl的后缀字母l说明我们要把从-0x8(%rbp) 这个地址开始的4个内存单元赋值为0。可能有读者会问，那如果我要操作3个，或5个内存单元呢？
很遗憾的是cpu没有提供相应的单条指令，我们只能通过多条指令组合起来达到目的

主要专注于AMD64 Linux平台下的go调度器，因此下面我们只介绍该平台下所使用的AT&T格式的汇编指令，AT&T汇编指令的基本格式为：
```shell
操作码  [操作数]
add   %rdx,%rax  
```
可以看到每一条汇编指令通常都由两部分组成：
    
- 操作码：操作码指示CPU执行什么操作，比如是执行加法，减法还是读写内存。每条指令都必须要有操作码。

- 操作数：操作数是操作的对象，比如加法操作需要两个加数，这两个加数就是这条指令的操作数。操作数的个数一般是0个，1个或2个。

AT&T格式的汇编指令中，寄存器名需要加%作为前缀.
这条指令的操作码是add，表示执行加法操作，它有两个操作数，rdx和rax。
如果一条指令有两个操作数，那么第一个操作数叫做源操作数，第二个操作数叫做目的操作数，

    顾名思义，目的操作数表示这条指令执行完后结果应该保存的地方。所以上面这条指令表示对rax和rdx寄存器里面的值求和，并把结果保存在rax寄存器中。
    其实这条指令的第二个操作数rax寄存器既是源操作数也是目的操作数，因为rax既是加法操作的两个加数之一，又得存放加法操作的结果。
    这条指令执行完后rax寄存器的值发生了改变，指令执行前的值被覆盖而丢失了，如果rax寄存器之前的值还有用，那么就得先用指令把它保存到其它寄存器或内存之中

再来看一个只有一个操作数的例子
```shell
callq  0x400526 
#这条指令的操作码是callq，表示调用函数，操作数是0x400526，它是被调用函数的地址。
```

最后来看一条没有操作数的指令：
```shell
retq
#这条指令只有操作码retq，表示从被调用函数返回到调用函数继续执行
```

## 常用指令详解
### mov指令-传送指令
```shell

mov 源操作数 目的操作数
mov %rsp,%rbp       # 直接寻址，把rsp的值拷贝给rbp，相当于 rbp = rsp
mov -0x8(%rbp),%edx # 源操作数间接寻址，目的操作数直接寻址。从内存中读取4个字节到edx寄存器
mov %rsi,-0x8(%rbp) # 源操作数直接寻址，目的操作数间接寻址。把rsi寄存器中的8字节值写入内存
```
### add/sub指令
```shell
add 源操作数 目的操作数
sub 源操作数 目的操作数
sub $0x350,%rsp  # 源操作数是立即操作数，目的操作数直接寻址。rsp = rsp - 0x350
add %rdx,%rax    # 直接寻址。rax = rax + rdx
addl $0x1,-0x8(%rbp) # 源操作数是立即操作数，目的操作数间接寻址。内存中的值加1（addl后缀字母l表示操作内存中的4个字节）
```

### call/ret指令
![](.introduction_images/retf.png) 
```shell
# 解析
call 目标地址  (相当于=>) push pc(push ip); jmp to callee addr;
ret (相当于=>) pop pc;

```
call指令执行函数调用。CPU执行call指令时首先会把rip寄存器中的值入栈，然后设置rip值为目标地址，又因为rip寄存器决定了下一条需要执行的指令，
所以当CPU执行完当前call指令后就会跳转到目标地址去执行。

ret指令从被调用函数返回调用函数，它的实现原理是把call指令入栈的返回地址弹出给rip寄存器

```shell
#调用函数片段
0x0000000000400559 : callq 0x400526 <sum>
0x000000000040055e : mov   %eax,-0x4(%rbp)

#被调用函数片段
0x0000000000400526 : push   %rbp
......
0x000000000040053f : retq 
```
![](.introduction_images/call_n_ret.png)

![](.introduction_images/callq_0x400526.png)

    从上图可以看到call指令执行之初rip寄存器的值是紧跟call后面那一条指令的地址，即0x40055e，但当call指令完成后但还未开始执行下一条指令之前，
    rip寄存器的值变成了call指令的操作数，即被调用函数的地址0x400526，这样CPU就会跳转到被调用函数去执行了。
    
    同时还需要注意的是这里的call指令执行时把call指令后面那一条指令的地址 0x40055e PUSH到了栈上，所以一条call指令修改了3个地方的值：rip寄存器、rsp和栈
![](.introduction_images/retq.png)

    可以看到ret指令执行的操作跟call指令执行的操作完全相反，ret指令开始执行时rip寄存器的值是紧跟ret指令后面的那个地址，也就是0x400540，
    但ret指令执行过程中会把之前call指令PUSH到栈上的返回地址 0x40055e POP给rip寄存器，这样，当ret执行完成后就会从被调用函数返回到调用函数的call指令的下一条指令继续执行。
    这里同样要注意的是retq指令也会修改rsp寄存器的

### cmp指令：影响标志寄存器
![](.introduction_images/cmp.png)
通过做减法运算，影响标志寄存器，标志寄存器的相关位记录了比较的结果

### jmp/je/jle/jg/jge等等j开头的指令--转移指令，例如可以修改8086cpu的cs段寄存器，ip指令寄存器
![](.introduction_images/jmp.png)
![](.introduction_images/jmp_english.png)
这些都属于跳转指令，操作码后面直接跟要跳转到的地址或存有地址的寄存器，这些指令与高级编程语言中的 goto 和 if 等语句对应。用法示例：
```shell
jmp    0x4005f2 #-->相当于jmp IP  0x4005f2,仅仅修改ip指令寄存器
jle    0x4005ee
jl     0x4005b8
```
通常cmp和跳转指令为一对，就像call和ret指令一样.

### push/pop指令-可以直接操作段寄存器
```shell
push 源操作数
pop 目的操作数
```
![](.introduction_images/push.png)

专用于函数调用栈的入栈出栈指令，这两个指令都会自动修改rsp寄存器
push入栈时rsp寄存器的值先减去8把栈位置留出来，然后把操作数复制到rsp所指位置。push指令相当于

    sub $8,%rsp
    mov 源操作数,(%rsp)
    push指令需要重点注意rsp寄存器的变化。

pop出栈时先把rsp寄存器所指位置的数据复制到目的操作数中，然后rsp寄存器的值加8。pop指令相当于

    mov (%rsp),目的操作数
    add $8,%rsp
    同样，pop指令也需要重点注意rsp寄存器的变化

### leave指令
leave指令没有操作数，它一般放在函数的尾部ret指令之前，用于调整rsp和rbp，这条指令相当于如下两条指令

    mov %rbp,%rsp
    pop %rbp
    
### loop指令
![](.introduction_images/loop.png)
cx存储循环的次数,s为标号

### shl和shr逻辑移位指令
shl: 将一个寄存器或内存单元中的数据向左移位，将最后移出的一位写入CF中，最低位用0补充。
```assembly
mov al,01001000b;
shl al, 1;
```
结果al=1001000b,CF=0

如果超过1位，移动位数放在cl中.
```assembly
mov al, 01001000b;
mov cl, 3;
shl al, cl;
```


## 定位方式
![](.introduction_images/locate_addr.png)
bx+idata，为高级语言提供方便
![](.introduction_images/bx+idata.png)

## 案例:c语言中
![](.introduction_images/location_in_c.png)

### 和系统打交道
程序的基本分段
```assembly
.data : 有初始化值的全局变量；定义常量。 .bss : 没有初始化值的全局变量。
.text : 代码段。
.rodata: 只读数据段。

```

```assembly
# 数据段
section .data
message: db 'hello, world!', 10

# 代码段
section .text
global _start

_start:
    mov     rax, 1              ; 'write' syscall number
    mov     rdi, 1              ; stdout descriptor
    mov     rsi, message        ; string address
    mov     rdx, 14             ; string length in bytes
    syscall   ; 重点在这里

```

## 案例:go语言编写
```go
c=a+b
```
对应AMD64 Linux平台代码

```shell
mov   (%rsp),%rdx          //把变量a的值从内存中读取到寄存器rdx中
mov    0x8(%rsp),%rax   //把变量b的值从内存中读取到寄存器rax中
add   %rdx,%rax             //把寄存器rdx和rax中的值相加，并把结果放回rax寄存器中
mov   %rax,0x10(%rsp)  //把寄存器rax中的值写回变量c所在的内存
```
![](.introduction_images/asm_in_memory.png)     

对这个图做个简单的说明：

    这里假定rsp寄存器的值是X
    
    图中的内存部分，每一行有8个内存单元，它们的地址从右向左依次加一，即如果最右边的内存单元的地址为X的话，则同一行最左边的内存单元的地址为X+7。
    
    灰色箭头表述数据流动方向
    
    紫红色数字n表示上述代码片段中的第n条指令

对内存部分介绍
    
1. 内存中的每个字节都有一个地址；

2. 任何大于一个字节的变量在内存中都存储在相邻连续的的几个内存单元之中；

3. 大端存储模式指数据的高字节保存在内存的低地址中，低字节保存在内存的高地址中；小端存储模式指数据的高字节保存在内存的高地址中，低字节保存在内存的低地址中。

![](.introduction_images/big_small_end_in_mem.png)
- 大端存储模式：数据的高字节保存在内存的低地址中，低字节保存在内存的高地址中。
- 小端存储模式：数据的高字节保存在内存的高地址中，低字节保存在内存的低地址中。

注意的是大小端存储模式与CPU相关，而与内存无关，内存只管保存数据而不关心数据是什么以及怎么解释这些数据.

## 中断
### 内中断过程
内中断的中断类型码是由cpu内部产生的.
![img.png](.introduction_images/inter_interupt.png)

### 外中断
外中断分为可屏蔽和不可屏蔽，可屏蔽中断信息来自cpu外部，中断类型码通过数据总线送入cpu。
#### 可屏蔽
与内中断就第一步中断类型码不同。

案例：pc端键盘
![img.png](.introduction_images/keyboard_process.png)
![img.png](.introduction_images/keyboard_process2.png)

#### 不可屏蔽
![img.png](.introduction_images/unstop_outer_interupt.png)
标志寄存器IF设置位0，禁止其他的可屏蔽中断。8086cpu，sti设置IF=1,cli设置IF=0。


## 标号
![img.png](.introduction_images/code_before.png)
![img.png](.introduction_images/code_after.png)
start,s仅仅表示内存单元的地址。还有的可以内存单元的长度，注意冒号变化。
![img.png](.introduction_images/assume_data.png)



