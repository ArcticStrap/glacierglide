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

# Setup commands
gencert:
	mkdir -p certs
	openssl ecparam -genkey -name secp384r1 -out certs/key.pem
	openssl req -new -x509 -sha256 -key certs/key.pem -out certs/cert.pem -days 3650 

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
