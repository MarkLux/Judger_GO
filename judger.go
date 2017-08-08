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
	MaxCpuTime       int
	MaxRealTime      int
	MaxMemory        int
	MaxStack         int
	MaxProcessNumber int
	MaxOutPutSize    int
	ExePath          string
	InputPath        string
	OutputPath       string
	ErrorPath        string
	LogPath          string
	Args             [ARGS_MAX_NUMBER]string
	Env              [ENV_MAX_NUMBER]string
	SecCompRuleName  string
	Uid              uint
	Gid              uint
}

type Result struct {
	CpuTime  int
	RealTime int
	Memory   int
	Signal   int
	ExitCode int
	Error    ErrorCode
	Result   ResultCode
}

func JudgerRun(config Config) Result {

	var _config C.struct_config = parseConfig(config)
	var _result C.struct_result

	defer freeConfig(&_config)
	defer freeArgs(_config.args, len(config.Args))
	defer freeEnv(_config.env, len(config.Env))

	C.judger_run(&_config, &_result)

	return parseResult(_result)
}

func parseResult(r C.struct_result) Result {
	var p Result
	p.CpuTime = int(r.cpu_time)
	p.RealTime = int(r.real_time)
	p.Memory = int(r.memory)
	p.Signal = int(r.signal)
	p.ExitCode = int(r.exit_code)
	p.Error = ErrorCode(r.error)
	p.Result = ResultCode(r.result)
	return p
}

func parseConfig(c Config) C.struct_config {
	var p C.struct_config

	p.max_cpu_time = C.int(c.MaxCpuTime)
	p.max_real_time = C.int(c.MaxRealTime)
	p.max_memory = C.long(c.MaxMemory)
	p.max_stack = C.long(c.MaxStack)
	p.max_process_number = C.int(c.MaxProcessNumber)
	p.max_output_size = C.long(c.MaxOutPutSize)
	p.exe_path = C.CString(c.ExePath)
	p.input_path = C.CString(c.InputPath)
	p.output_path = C.CString(c.OutputPath)
	p.error_path = C.CString(c.ErrorPath)
	p.log_path = C.CString(c.LogPath)
	p.args = parseArgs(c.Args)
	p.env = parseEnv(c.Env)
	p.log_path = C.CString(c.LogPath)
	p.seccomp_rule_name = C.CString(c.SecCompRuleName)

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
