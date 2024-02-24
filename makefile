# Project makefile

# Locations
OUTPUT = bin/horinezumi
MAIN = ./cmd/server/main.go

# Target commands
BUILDCMD = go build -o $(OUTPUT) $(MAIN)
CLEANCMD = rm -f $(OUTPUT)

ifeq ($(OS), Windows_NT) # Windows
	OUTPUT = bin\horinezumi.exe
	CLEANCMD = del $(OUTPUT)
endif

# Targets
.PHONY: all help gencert getdeps build clean dversion run test benchmark tidy

# Default target
all: help

# General commands
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "	gencert		Generate SSL certificate"
	@echo "	getdeps		Get the project dependencies"
	@echo "	genschema		Generate database schemas"
	@echo "	build		Build the project"
	@echo "	clean		Clean the project"
	@echo "	dversion	Display databse migration version"
	@echo "	run		Run the project"
	@echo "	test		Run unit tests"
	@echo "	benchmark	Run benchmarks"
	@echo "	tidy		Tidy up mod file"
	@echo ""
	@echo "Variables:"
	@echo "	DEV		Run in development mode"

# Setup commands
gencert:
	mkdir -p certs
	openssl ecparam -genkey -name secp384r1 -out certs/key.pem
	openssl req -new -x509 -sha256 -key certs/key.pem -out certs/cert.pem -days 3650

genschema:
	go run ./cmd/genschema/main.go

getdeps:
	go mod download

# Go commands
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

benchmark:
	go test -bench=. ./...

tidy:
	go mod tidy

