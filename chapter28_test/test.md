
# test


## go test接受的参数
```shell
➜  03_n git:(feature/memory) ✗ go help test                                                                             
usage: go test [build/test flags] [packages] [build/test flags & test binary flags]

'Go test' automates testing the packages named by the import paths.
It prints a summary of the test results in the format:

        ok   archive/tar   0.011s
        FAIL archive/zip   0.022s
        ok   compress/gzip 0.033s
        ...

followed by detailed output for each failed package.

'Go test' recompiles each package along with any files with names matching
the file pattern "*_test.go".
These additional files can contain test functions, benchmark functions, fuzz
tests and example functions. See 'go help testfunc' for more.
Each listed package causes the execution of a separate test binary.
Files whose names begin with "_" (including "_test.go") or "." are ignored.

Test files that declare a package with the suffix "_test" will be compiled as a
separate package, and then linked and run with the main test binary.

The go tool will ignore a directory named "testdata", making it available
to hold ancillary data needed by the tests.

As part of building a test binary, go test runs go vet on the package
and its test source files to identify significant problems. If go vet
finds any problems, go test reports those and does not run the test
binary. Only a high-confidence subset of the default go vet checks are
used. That subset is: 'atomic', 'bool', 'buildtags', 'errorsas',
'ifaceassert', 'nilfunc', 'printf', and 'stringintconv'. You can see
the documentation for these and other vet tests via "go doc cmd/vet".
To disable the running of go vet, use the -vet=off flag. To run all
checks, use the -vet=all flag.

All test output and summary lines are printed to the go command's
standard output, even if the test printed them to its own standard
error. (The go command's standard error is reserved for printing
errors building the tests.)

Go test runs in two different modes:

The first, called local directory mode, occurs when go test is
invoked with no package arguments (for example, 'go test' or 'go
test -v'). In this mode, go test compiles the package sources and
tests found in the current directory and then runs the resulting
test binary. In this mode, caching (discussed below) is disabled.
After the package test finishes, go test prints a summary line
showing the test status ('ok' or 'FAIL'), package name, and elapsed
time.

The second, called package list mode, occurs when go test is invoked
with explicit package arguments (for example 'go test math', 'go
test ./...', and even 'go test .'). In this mode, go test compiles
and tests each of the packages listed on the command line. If a
package test passes, go test prints only the final 'ok' summary
line. If a package test fails, go test prints the full test output.
If invoked with the -bench or -v flag, go test prints the full
output even for passing package tests, in order to display the
requested benchmark results or verbose logging. After the package
tests for all of the listed packages finish, and their output is
printed, go test prints a final 'FAIL' status if any package test
has failed.

In package list mode only, go test caches successful package test
results to avoid unnecessary repeated running of tests. When the
result of a test can be recovered from the cache, go test will
redisplay the previous output instead of running the test binary
again. When this happens, go test prints '(cached)' in place of the
elapsed time in the summary line.

The rule for a match in the cache is that the run involves the same
test binary and the flags on the command line come entirely from a
restricted set of 'cacheable' test flags, defined as -benchtime, -cpu,
-list, -parallel, -run, -short, -timeout, -failfast, and -v.
If a run of go test has any test or non-test flags outside this set,
the result is not cached. To disable test caching, use any test flag
or argument other than the cacheable flags. The idiomatic way to disable
test caching explicitly is to use -count=1. Tests that open files within
the package's source root (usually $GOPATH) or that consult environment
variables only match future runs in which the files and environment
variables are unchanged. A cached test result is treated as executing
in no time at all,so a successful package test result will be cached and
reused regardless of -timeout setting.

In addition to the build flags, the flags handled by 'go test' itself are:

        -args
           -args 该参数会被原封不动 (uninterpreted and unchanged) 的传递给测试二进制文件
            Because this flag consumes the remainder of the command line,
            the package list (if present) must appear before this flag.

        -c
            只编译不执行
            (where pkg is the last element of the package's import path).
            The file name can be changed with the -o flag.

        -exec xprog
            使用xprog来执行测试二进制程序. The behavior is the same as
            in 'go run'. See 'go help run' for details.

        -i
            安装测试依赖的包，不执行测试
            The -i flag is deprecated. Compiled packages are cached automatically.

        -json
            将输出格式化为json格式，用于自动化处理结果
            See 'go doc test2json' for the encoding details.

        -o file
            Compile the test binary to the named file.
            The test still runs (unless -c or -i is specified).

The test binary also accepts flags that control execution of the test; these
flags are also accessible by 'go test'. See 'go help testflag' for details.

For more about build flags, see 'go help build'.
For more about specifying packages, see 'go help packages'.

See also: go build, go vet.
```



