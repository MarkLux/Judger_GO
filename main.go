package main

/*
#cgo LDFLAGS: -ldl
#include <stdio.h>
#include <dlfcn.h>
#include <stdlib.h>
#include "runner.h"

void judger_run(struct config * _config,struct result * _result) {

    //load dynamic library from /usr/lib/judger
    void * handler =  dlopen("/usr/lib/judger/libjudger.so",RTLD_LAZY);

    if (!handler) {
        _result->error = LOAD_JUDGER_FAILED;
        _result->result = SYSTEM_ERROR;
        return;
    }

    int (*judger_run)(struct config *,struct result *);

    judger_run = dlsym(handler,"run");

    judger_run(_config,_result);

    return;
}
*/
import "C"
import "fmt"

const ARGS_MAX_NUMBER int = 256
const ENV_MAX_NUMBER int = 256

type ResultCode int
type ErrorCode int

const (
	SUCCESS             = 0
	INVALID_CONFIG      = -1
	FORK_FAILED         = -2
	PTHREAD_FAILED      = -3
	WAIT_FAILED         = -4
	ROOT_REQUIRED       = -5
	LOAD_SECCOMP_FAILED = -6
	SETRLIMIT_FAILED    = -7
	DUP2_FAILED         = -8
	SETUID_FAILED       = -9
	EXECVE_FAILED       = -10
	SPJ_ERROR           = -11
)

const (
	WRONG_ANSWER             = -1
	CPU_TIME_LIMIT_EXCEEDED  = 1
	REAL_TIME_LIMIT_EXCEEDED = 2
	MEMORY_LIMIT_EXCEEDED    = 3
	RUNTIME_ERROR            = 4
	SYSTEM_ERROR             = 5
)

type Config struct {
	maxCpuTime       int
	maxRealTime      int
	maxMemory        int
	maxStack         int
	maxProcessNumber int
	maxOutPutSize    int
	exePath          string
	inputPath        string
	outputPath       string
	errorPath        string
	logPath          string
	args             [ARGS_MAX_NUMBER]string
	env              [ENV_MAX_NUMBER]string
	secCompRuleName  string
	uid              uint
	gid              uint
}

type Result struct {
	cpuTime  int
	realTime int
	memory   int
	signal   int
	exitCode int
	error    ErrorCode
	result   ResultCode
}

func main() {

	var args [ARGS_MAX_NUMBER]string
	var env [ENV_MAX_NUMBER]string
	args[0] = "HelloWorld"

	testConfig := Config{
		maxCpuTime:       1000,
		maxRealTime:      2000,
		maxMemory:        128 * 1024 * 1024,
		maxProcessNumber: 200,
		maxOutPutSize:    10000,
		maxStack:         32 * 1024 * 1024,
		exePath:          "/bin/bash",
		inputPath:        "echo.in",
		outputPath:       "echo.out",
		args:             args,
		env:              env,
		logPath:          "echo.log",
		secCompRuleName:  "c_cpp",
		uid:              0,
		gid:              0,
	}

	var testResult Result

	var _config C.struct_config = parseConfig(testConfig)
	var _result C.struct_result = parseResult(testResult)
	C.judger_run(&_config, &_result)
	fmt.Println(_result.result)
}

func parseResult(r Result) C.struct_result {
	var p C.struct_result
	p.cpu_time = C.int(r.cpuTime)
	p.real_time = C.int(r.realTime)
	p.memory = C.long(r.memory)
	p.signal = C.int(r.signal)
	p.exit_code = C.int(r.exitCode)
	p.error = C.int(r.error)
	p.result = C.int(r.result)
	return p
}

func parseConfig(c Config) C.struct_config {
	var p C.struct_config

	p.max_cpu_time = C.int(c.maxCpuTime)
	p.max_real_time = C.int(c.maxRealTime)
	p.max_memory = C.long(c.maxMemory)
	p.max_stack = C.long(c.maxStack)
	p.max_process_number = C.int(c.maxProcessNumber)
	p.max_output_size = C.long(c.maxOutPutSize)
	p.exe_path = C.CString(c.exePath)
	p.input_path = C.CString(c.inputPath)
	p.output_path = C.CString(c.outputPath)
	p.error_path = C.CString(c.errorPath)
	p.log_path = C.CString(c.logPath)
	p.args = parseArgs(c.args)
	p.env = parseEnv(c.env)
	p.log_path = C.CString(c.logPath)
	p.seccomp_rule_name = C.CString(c.secCompRuleName)

	return p
}

func parseArgs(goArray [ARGS_MAX_NUMBER]string) [ARGS_MAX_NUMBER]*C.char {
	var p [ARGS_MAX_NUMBER]*C.char
	for i := 0; i < len(goArray); i++ {
		p[i] = C.CString(goArray[i])
	}
	return p
}

func parseEnv(goArray [ENV_MAX_NUMBER]string) [ENV_MAX_NUMBER]*C.char {
	var p [ARGS_MAX_NUMBER]*C.char
	for i := 0; i < len(goArray); i++ {
		p[i] = C.CString(goArray[i])
	}
	return p
}
