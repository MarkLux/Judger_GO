package judger

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

import "unsafe"

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

func JudgerRun(config Config) Result {

	var _config C.struct_config = parseConfig(config)
	var _result C.struct_result

	defer freeConfig(&_config)
	defer freeArgs(_config.args, len(config.args))
	defer freeEnv(_config.env, len(config.env))

	C.judger_run(&_config, &_result)

	return parseResult(_result)
}

func parseResult(r C.struct_result) Result {
	var p Result
	p.cpuTime = int(r.cpu_time)
	p.realTime = int(r.real_time)
	p.memory = int(r.memory)
	p.signal = int(r.signal)
	p.exitCode = int(r.exit_code)
	p.error = ErrorCode(r.error)
	p.result = ResultCode(r.result)
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

// memory manage

func freeConfig(c *C.struct_config) {

	// free space that out of go's gc

	C.free(unsafe.Pointer(c.exe_path))
	C.free(unsafe.Pointer(c.input_path))
	C.free(unsafe.Pointer(c.output_path))
	C.free(unsafe.Pointer(c.error_path))
	C.free(unsafe.Pointer(c.log_path))
	C.free(unsafe.Pointer(c.seccomp_rule_name))
}

func freeArgs(args [ARGS_MAX_NUMBER]*C.char, length int) {
	for i := 0; i < length; i++ {
		C.free(unsafe.Pointer(args[i]))
	}
}

func freeEnv(env [ENV_MAX_NUMBER]*C.char, length int) {
	for i := 0; i < length; i++ {
		C.free(unsafe.Pointer(env[i]))
	}
}
