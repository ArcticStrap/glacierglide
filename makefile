# Project makefile
OUTPUT = bin/horinezumi
MAIN = ./cmd/server/main.go

# Target commands
BUILDCMD = go build -o $(OUTPUT) $(MAIN)
CLEANCMD = rm -f $(OUTPUT)

ifeq ($(OS), Windows_NT) # Windows
	OUTPUT = bin\horinezumi.exe
	CLEANCMD = del $(OUTPUT)
endif

all: clean build

build:
	$(BUILDCMD)
clean:
	$(CLEANCMD)
dversion:
	migrate -database $(URL) -path data/migrations version
run:
	go run $(MAIN)
test:
	go test -v ./...
