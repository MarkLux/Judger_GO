# Judger Binder for GO

> a go binding for qduoj judger

need dynamic library from qduoj judger and seccomp,only for linux x64 platform

* [reference of judger](https://github.com/QingdaoU/Judger)

## install&usage

(tested on ubuntu:14.04)

* insatll

```
# install depencies
sudo apt-get install libseccomp-dev
go get -u -v github.com/MarkLux/Judger_GO
# enter your go path to find this pkg src
cd $GOPATH/src/github.com/MarkLux/Judger_GO
# move the dynamic library
sudo mkdir /usr/lib/judger
sudo mv libjudger.so /usr/lib/judger/libjudger.so
```

* usage(example)

```
package main

import (
	"fmt"

	"github.com/MarkLux/Judger_GO"
)

func main() {
	testConfig := judger.Config{
		MaxCpuTime:       1000,
		MaxRealTime:      2000,
		MaxMemory:        128 * 1024 * 1024,
		MaxProcessNumber: 200,
		MaxStack:         32 * 1024 * 1024,
		MaxOutPutSize:    10000,
		ExePath:          "bin/bash",
		ErrorPath:        "echo.err",
		OutputPath:       "echo.out",
		InputPath:        "echo.in",
		LogPath:          "echo.log",
		Args:             [judger.ARGS_MAX_NUMBER]string{"hello World"},
		Env:              [judger.ENV_MAX_NUMBER]string{"foo=bar"},
		SecCompRuleName:  "c_cpp",
		Uid:              0,
		Gid:              0,
	}

	testResult := judger.JudgerRun(testConfig)

	fmt.Println("result = ", testResult.Result)
}
```

## notice
due to the usage of cgo,you may need to execute `go build` then run the programm instead of execute `go run` directly.