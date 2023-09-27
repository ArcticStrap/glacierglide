# Project makefile
# Make sure to run some sort of bash
OUTPUT = bin/horinezumi
BUILDCMD = go build -o $(OUTPUT) main.go

ifeq ($(OS), Windows_NT) # Windows
	OUTPUT = bin\horinezumi.exe
endif

all: clean build

build:
	$(BUILDCMD)
clean:
	del $(OUTPUT)
run:
	go run main.go
test:
	go test -v
