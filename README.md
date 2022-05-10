# diversion

## Overview
Diversion is able to bypasses Amsi and ETW in both 32bit and 64bit processes.


## Requirements
Diversion is written in Go and requires the windows library.

```bash
go get golang.org/x/sys/windows
```

Once installed, you can build the binary using the Makefile or the build command below on Linux or Mac:

```bash
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o diversion.exe
GOOS=windows GOARCH=386 go build -ldflags "-s -w" -trimpath -o diversion32.exe
```

## Examples
Diversion requires the PID of the process to inject into and the bypass method.

```bash
./diversion.exe
  -method string
        Evasion method: amsi, etw (default "amsi")
  -pid int
        PID of the process to inject into
```

### Amsi

Simple check to see that running `AmsiUtil` gets flagged as malicious before bypassing Amsi. After running diversion, the command is no longer flagged as malicious.

![amsi.png](img/amsi.png)

### ETW

Using ProcessHacker to view the loaded .NET Assemblies within a process when using execute-assembly.

![etw1.png](img/etw1.png)

No .NET assemblies will be displayed when bypassing ETW.

![etw2.png](img/etw2.png)

## Future Improvements
- DLL unhooking
- Additional evasions/bypasses

## References / Credit

Amsi bypass technique converted from boku7's BOFF.
- https://github.com/boku7/injectAmsiBypass

ETW bypass technique coverted from mdsec example. Boku7 also created a BOFF file which uses Syscalls.
- https://www.mdsec.co.uk/2020/03/hiding-your-net-etw/
- https://github.com/boku7/injectEtwBypass
