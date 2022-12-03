# link


## 选项
```shell
➜  03_n git:(feature/memory) ✗ go tool link --help
usage: link [options] main.o
  -B note
        add an ELF NT_GNU_BUILD_ID note when using ELF
  -E entry
        set entry symbol name
  -H type
        set header type
  -I linker
        use linker as ELF dynamic linker
  -L directory
        add specified directory to library path
  -R quantum
        set address rounding quantum (default -1)
  -T address
        set text segment address (default -1)
  -V    print version and exit
  -X definition
        add string value definition of the form importpath.name=value
  -a    no-op (deprecated)
  -asan
        enable ASan interface
  -aslr
        enable ASLR for buildmode=c-shared on windows (default true)
  -benchmark string
        set to 'mem' or 'cpu' to enable phase benchmarking
  -benchmarkprofile base
        emit phase profiles to base_phase.{cpu,mem}prof
  -buildid id
        record id as Go toolchain build id
  -buildmode mode
        set build mode
  -c    dump call graph
  -compressdwarf
        compress DWARF if possible (default true)
  -cpuprofile file
        write cpu profile to file
  -d    disable dynamic executable
  -debugtextsize int
        debug text section max size
  -debugtramp int
        debug trampolines
  -dumpdep
        dump symbol dependency graph
  -extar string
        archive program for buildmode=c-archive
  -extld linker
        use linker when linking in external mode
  -extldflags flags
        pass flags to external linker
  -f    ignore version mismatch
  -g    disable go package data checks
  -h    halt on error
  -importcfg file
        read import configuration from file
  -installsuffix suffix
        set package directory suffix
  -k symbol
        set field tracking symbol
  -libgcc string
        compiler support lib for internal linking; use "none" to disable
  -linkmode mode
        set link mode
  -linkshared
        link against installed Go shared libraries
  -memprofile file
        write memory profile to file
  -memprofilerate rate
        set runtime.MemProfileRate to rate
  -msan
        enable MSan interface
  -n    dump symbol table
  -o file
        write output to file
  -pluginpath string
        full path name for plugin
  -r path
        set the ELF dynamic linker search path to dir1:dir2:...
  -race
        enable race detector
  -s    disable symbol table
  -strictdups int
        sanity check duplicate symbol contents during object file reading (1=warn 2=err).
  -tmpdir directory
        use directory for temporary files
  -v    print link trace
  -w    disable DWARF generation

```