package judger

import (
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
)

type ResultCode int
type ErrorCode int

const UNLIMITED = -1

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
	Args             []string
	Env              []string
	SecCompRuleName  string
	Uid              int
	Gid              int
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

func JudgerRun(config Config) (Result, error) {
	args := []string{}
	// parsing args
	args = append(args, "--max_cpu_time="+strconv.Itoa(config.MaxCpuTime))
	args = append(args, "--max_real_time="+strconv.Itoa(config.MaxRealTime))
	args = append(args, "--max_memory="+strconv.Itoa(config.MaxMemory))
	args = append(args, "--max_process_number="+strconv.Itoa(config.MaxProcessNumber))
	args = append(args, "--max_stack="+strconv.Itoa(config.MaxStack))
	args = append(args, "--max_output_size="+strconv.Itoa(config.MaxOutPutSize))
	args = append(args, "--exe_path="+config.ExePath)
	args = append(args, "--input_path="+config.InputPath)
	args = append(args, "--output_path="+config.OutputPath)
	args = append(args, "--error_path="+config.ErrorPath)
	args = append(args, "--log_path="+config.LogPath)
	// parsing list args
	for _, arg := range config.Args {
		args = append(args, "--args="+arg)
	}

	for _, env := range config.Env {
		args = append(args, "--env="+env)
	}

	args = append(args, "--seccomp_rule_name="+config.SecCompRuleName)
	if config.Uid > 0 {
		args = append(args, "--uid="+strconv.Itoa(config.Uid))
	}
	if config.Gid > 0 {
		args = append(args, "--gid="+strconv.Itoa(config.Gid))
	}

	cmd := exec.Command("/usr/lib/judger/libjudger.so", args...)

	output, err := cmd.Output()

	if err != nil {
		log.Println(err)
	}

	var runResult Result

	if err := json.Unmarshal(output, &runResult); err != nil {
		log.Println(err)
	}

	return runResult, err
}

