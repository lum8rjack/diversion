NAME := diversion
BUILD := go build -ldflags "-s -w" -trimpath

default: windows

clean:
	rm -f $(NAME)*.exe


windows:
	echo "Compiling for Windows x64"
	GOOS=windows GOARCH=amd64 $(BUILD) -o $(NAME).exe

windows32:
	echo "Compiling for Windows x86"
	GOOS=windows GOARCH=386 $(BUILD) -o $(NAME)32.exe
