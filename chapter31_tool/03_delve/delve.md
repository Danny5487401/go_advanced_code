<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Debugger 调试器 delve](#debugger-%E8%B0%83%E8%AF%95%E5%99%A8-delve)
  - [背景](#%E8%83%8C%E6%99%AF)
  - [Go 目前的调试器](#go-%E7%9B%AE%E5%89%8D%E7%9A%84%E8%B0%83%E8%AF%95%E5%99%A8)
  - [基本术语](#%E5%9F%BA%E6%9C%AC%E6%9C%AF%E8%AF%AD)
  - [调试器分类](#%E8%B0%83%E8%AF%95%E5%99%A8%E5%88%86%E7%B1%BB)
    - [Instruction level debugger 指令级调试器](#instruction-level-debugger-%E6%8C%87%E4%BB%A4%E7%BA%A7%E8%B0%83%E8%AF%95%E5%99%A8)
    - [Symbol level debugger 符号级调试器](#symbol-level-debugger-%E7%AC%A6%E5%8F%B7%E7%BA%A7%E8%B0%83%E8%AF%95%E5%99%A8)
  - [使用](#%E4%BD%BF%E7%94%A8)
  - [dlv前端(github.com/aarzilli/gdlv)](#dlv%E5%89%8D%E7%AB%AFgithubcomaarzilligdlv)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Debugger 调试器 delve


## 背景


开发人员定位代码中的问题时，通常会借助 Print语句（如fmt.Println） 来打印变量的值，进而推断程序执行结果是否符合预期。在某些更复杂的场景下，打印语句可能难以胜任，调试器会更好地协助我们定位问题。


调试器可以帮助我们控制tracee（被调试进程、线程）的执行，也可以观察tracee的运行时内存、寄存器状态，借此我们可以实现代码的逐语句执行、控制代码执行流程、检查变量值是否符合预期，等等



调试器要支持的重要操作，通常包括：

- 设置断点，在指定内存地址、函数、语句、文件行号处设置断点；
- 单步执行，单步执行一条指令，单步执行一条语句，或运行到下个断点处；
- 获取、设置寄存器信息；
- 获取、设置内存信息；
- 对表达式进行估值计算；
- 调用函数；
- 其他；
## Go 目前的调试器

- GDB 最早期的调试工具，现在用的很少。对goroutine场景支持不足，不能很好的应对goroutine的调试
- LLDB macOS 系统推荐的标准调试工具，但 Go 的一些专有特性支持的比较少。
- Delve 专门为 Go 语言打造的调试工具，使用最为广泛。


## 基本术语

DWARF：是一种调试信息标准，指导编译器将调试信息生成到目标文件中，指导链接器合并存储在多个目标文件中的调试信息，调试器将加载此调试信息。简言之，DWARF用来协调编译器、链接器和调试器之间的工作

## 调试器分类

调试器可以分为两种类型：指令级调试器和符号级调试器

### Instruction level debugger 指令级调试器

指令级调试器，其操作的对象是机器指令。通过处理器指令patch技术就可以实现指令级调试，不需要调试符号信息。
它仅适用于指令或汇编语言级别的操作，不支持源代码级别的操作。


### Symbol level debugger 符号级调试器

符号级调试器，其操作的对象不仅是机器指令，更重要的是支持源代码级的操作。它可以提取和解析调试符号信息，建立内存地址、指令地址和源代码之间的映射关系，支持在源代码语句上设置断点的时候，将其转换为精确的机器指令断点，也支持其他方便的操作


## 使用

```go
# Version: 1.22.1
➜  ~ dlv --help
Delve is a source level debugger for Go programs.

Delve enables you to interact with your program by controlling the execution of the process,
evaluating variables, and providing information of thread / goroutine state, CPU register state and more.

The goal of this tool is to provide a simple yet powerful interface for debugging Go programs.

Pass flags to the program you are debugging using `--`, for example:

`dlv exec ./hello -- server --config conf/config.toml`

Usage:
  dlv [command]

Available Commands:
  attach      Attach to running process and begin debugging.
  completion  Generate the autocompletion script for the specified shell
  connect     Connect to a headless debug server with a terminal client.
  core        Examine a core dump.
  dap         Starts a headless TCP server communicating via Debug Adaptor Protocol (DAP).
  debug       Compile and begin debugging main package in current directory, or the package specified.
  exec        Execute a precompiled binary, and begin a debug session.
  help        Help about any command
  test        Compile test binary and begin debugging program.
  trace       Compile and begin tracing program.
  version     Prints version.

Additional help topics:
  dlv backend    Help about the --backend flag.
  dlv log        Help about logging flags.
  dlv redirect   Help about file redirection.
```
|   命令   |                               解释                                |
|:------:|:---------------------------------------------------------------:| 
| attach | 这个命令将使Delve控制一个已经运行的进程，并开始一个新的调试会话。 当退出调试会话时，你可以选择让该进程继续运行或杀死它。 | 
| debug  |                                默认情况下，没有参数，Delve将编译当前目录下的 "main "包，并开始调试。或者，你可以指定一个包的名字，Delve将编译该包，并开始一个新的调试会话                                 |


进入 dlv 调试后的指令

```go
✗ dlv debug  
Type 'help' for list of commands.
(dlv) help
The following commands are available:

Running the program:
    call ------------------------ Resumes process, injecting a function call (EXPERIMENTAL!!!)
    continue (alias: c) --------- Run until breakpoint or program termination.
    next (alias: n) ------------- Step over to next source line.
    rebuild --------------------- Rebuild the target executable and restarts it. It does not work if the executable was not built by delve.
    restart (alias: r) ---------- Restart process.
    step (alias: s) ------------- Single step through program.
    step-instruction (alias: si)  Single step a single cpu instruction.
    stepout (alias: so) --------- Step out of the current function.

Manipulating breakpoints:
    break (alias: b) ------- Sets a breakpoint.
    breakpoints (alias: bp)  Print out info for active breakpoints.
    clear ------------------ Deletes breakpoint.
    clearall --------------- Deletes multiple breakpoints.
    condition (alias: cond)  Set breakpoint condition.
    on --------------------- Executes a command when a breakpoint is hit.
    toggle ----------------- Toggles on or off a breakpoint.
    trace (alias: t) ------- Set tracepoint.
    watch ------------------ Set watchpoint.

Viewing program variables and memory:
    args ----------------- Print function arguments.
    display -------------- Print value of an expression every time the program stops.
    examinemem (alias: x)  Examine raw memory at the given address.
    locals --------------- Print local variables.
    print (alias: p) ----- Evaluate an expression.
    regs ----------------- Print contents of CPU registers.
    set ------------------ Changes the value of a variable.
    vars ----------------- Print package variables.
    whatis --------------- Prints type of an expression.

Listing and switching between threads and goroutines:
    goroutine (alias: gr) -- Shows or changes current goroutine
    goroutines (alias: grs)  List program goroutines.
    thread (alias: tr) ----- Switch to the specified thread.
    threads ---------------- Print out info for every traced thread.

Viewing the call stack and selecting frames:
    deferred --------- Executes command in the context of a deferred call.
    down ------------- Move the current frame down.
    frame ------------ Set the current frame, or execute command on a different frame.
    stack (alias: bt)  Print stack trace.
    up --------------- Move the current frame up.

Other commands:
    config --------------------- Changes configuration parameters.
    disassemble (alias: disass)  Disassembler.
    dump ----------------------- Creates a core dump from the current process state
    edit (alias: ed) ----------- Open where you are in $DELVE_EDITOR or $EDITOR
    exit (alias: quit | q) ----- Exit the debugger.
    funcs ---------------------- Print list of functions.
    help (alias: h) ------------ Prints the help message.
    libraries ------------------ List loaded dynamic libraries
    list (alias: ls | l) ------- Show source code.
    packages ------------------- Print list of packages.
    source --------------------- Executes a file containing a list of delve commands
    sources -------------------- Print list of source files.
    target --------------------- Manages child process debugging.
    transcript ----------------- Appends command output to a file.
    types ---------------------- Print list of types

Type help followed by a command for full documentation.
```


## dlv前端(github.com/aarzilli/gdlv)



## 参考
- [delve源码](github.com/go-delve/delve)
- 如何开发Go语言调试器: https://github.com/hitzhangjie/golang-debugger-book

