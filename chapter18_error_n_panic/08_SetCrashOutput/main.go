package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

func main() {
	monitor()

	println("Danny下午好！！！")
	// 没有被 recover 的未知错误
	panic("oops")
}

func monitor() {
	const monitorVar = "RUNTIME_DEBUG_MONITOR"
	if os.Getenv(monitorVar) != "" {
		// 实际演示 debug.SetCrashOutput 设置后的逻辑
		log.SetFlags(0)
		log.SetPrefix("monitor: ")

		crash, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("failed to read from input pipe: %v", err)
		}
		if len(crash) == 0 {
			os.Exit(0)
		}

		f, err := os.CreateTemp("", "*.crash")
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write(crash); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		log.Fatalf("saved crash report at %s", f.Name())
	}

	// 模拟应用程序进程，设置 debug.SetCrashOutput 值
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(exe, "-test.run=ExampleSetCrashOutput_monitor")
	cmd.Env = append(os.Environ(), monitorVar+"=1")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stderr
	pipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("StdinPipe: %v", err)
	}
	debug.SetCrashOutput(pipe.(*os.File), debug.CrashOptions{})
	if err := cmd.Start(); err != nil {
		log.Fatalf("can't start monitor: %v", err)
	}

}
