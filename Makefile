NAME := diversion
BUILD := go build -ldflags "-s -w" -trimpath

default: windows

clean:
	rm -f $(NAME)*.exe


windows:
	echo "Compiling for Windows x64"
	GOOS=windows GOARCH=amd64 $(BUILD) -o $(NAME).exe
