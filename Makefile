########################################
# Commands
########################################

# Services
# make build    - Build binary
# make run      - Run service
# make clean    - Remove binary output
# make db       - Run docker database
# make clean-db - Clean DB volume

########################################
# Config
########################################

OUT_DIR=out
BIN_NAME=glower

ifeq ($(OS), Windows_NT)
	BIN_EXT=.exe
	MKDIR=if not exist $(OUT_DIR) mkdir $(OUT_DIR)
	RM=del /Q
else
	BIN_EXT=
	MKDIR=mkdir -p $(OUT_DIR)
	RM=rm -f
endif

########################################

# Build

build:
	$(MKDIR)
	go build -o $(OUT_DIR)/$(BIN_NAME)$(BIN_EXT) main.go

clean:
	$(RM) $(OUT_DIR)/$(BINARY_NAME)$(EXE)

run:
	go run main.go

# Tests

test-l1:
	go test -tags=L1 ./tests/...

test-l2:
	go test -tags=L2 ./tests/...

# Docker

db:
	cd docker && docker-compose -f glower-db.yaml up

clean-db:
	cd docker && docker-compose -f glower-db.yaml down -v

# Pipeline

test-l1-json:
	mkdir -p $(OUT_DIR)
	go test -json -tags=L1 ./tests/... > $(OUT_DIR)/l1-results.json || true

test-l2-json:
	mkdir -p $(OUT_DIR)
	go test -json -tags=L2 ./tests/... > $(OUT_DIR)/l2-results.json || true