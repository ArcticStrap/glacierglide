# Project makefile

# Locations
OUTPUT = bin/glacierglide
MAIN = ./cmd/server/main.go
MIGRATE = ./cmd/dbmigrate/migrate

# Target commands
BUILDCMD = go build -o $(OUTPUT) $(MAIN)
CLEANCMD = rm -f $(OUTPUT)

ifeq ($(OS), Windows_NT) # Windows
	OUTPUT = bin\glacierglide.exe
	CLEANCMD = del $(OUTPUT)
endif

# Targets
.PHONY: all help gencert genschema getdeps build clean migrateup migratedown migrateforce dversion run test benchmark tidy

# Default target
all: help

# General commands
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "	gencert		Generate SSL certificate"
	@echo "	getdeps		Get the project dependencies"
	@echo "	genschema	Generate database schemas"
	@echo "	build		Build the project"
	@echo "	clean		Clean the project"
	@echo "	dversion	Display databse migration version"
	@echo "	migrateup	Migrate up one version"
	@echo "	migratedown	Migrate down one version"
	@echo "	migrateforce	Migrate to a specific version"
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
	$(MIGRATE) -database $(URL) -path migrations version

migrateup:
	go run ./cmd/dbmigrate/dbmigrate.go up
migratedown:
	go run ./cmd/dbmigrate/dbmigrate.go down
migrateforce:
	go run ./cmd/dbmigrate/dbmigrate.go force $(VERSION)

run:
	go run $(MAIN)

test:
	go test -v ./...

benchmark:
	go test -bench=. ./...

tidy:
	go mod tidy

