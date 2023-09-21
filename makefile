# Project makefile
# Make sure to run some sort of bash
OUTPUT = bin/horinezumi
BUILDCMD = go build -o $(OUTPUT) main.go

ifeq ($(UNAME_S), Linux) # Linux
	BUILDCMD = echo Currently not supported on Linux
endif

ifeq ($(UNAME_S), Darwin) # Mac
	BUILDCMD = echo Currently not supported on Mac
endif

ifeq ($(OS), Windows_NT) # Windows
	OUTPUT = bin/horinezumi.exe
endif

all: build
run:
	go run main.go
build:
	$(BUILDCMD)