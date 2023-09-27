# Project makefile
OUTPUT = bin/horinezumi

# Target commands
BUILDCMD = go build -o $(OUTPUT) main.go
CLEANCMD = rm -f $(OUTPUT)

ifeq ($(OS), Windows_NT) # Windows
	OUTPUT = bin|horinezumi.exe
	CLEANCMD = del $(OUTPUT)
endif

all: clean build

build:
	$(BUILDCMD)
clean:
	$(CLEANCMD)
run:
	go run main.go
test:
	go test -v
