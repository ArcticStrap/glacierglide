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

# General commands
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "	gencert		Generate SSL certificate"
	@echo "	build		Build the project"
	@echo "	clean		Clean the project"
	@echo "	dversion	Display databse migration version"
	@echo "	run		Run the project"
	@echo "	getdeps		Get the project dependencies"
	@echo "	test		Run unit tests"
	@echo "	benchmark	Run benchmarks"
	@echo ""
	@echo "Variables:"
	@echo "	DEV		Run in development mode"

# Setup commands
gencert:
	mkdir -p certs
	openssl ecparam -genkey -name secp384r1 -out certs/key.pem
	openssl req -new -x509 -sha256 -key certs/key.pem -out certs/cert.pem -days 3650 
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
