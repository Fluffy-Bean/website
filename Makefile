TARGET = www
BUILD_DIR = dist
GENERATED_DIR = static/generated

BUILD_TIME = $(shell date +"%Y-%m-%d.%H:%M:%S")

clean:
	rm -Rf "${BUILD_DIR}"
	rm -Rf "${GENERATED_DIR}"

deps:
	go mod tidy

generate:
	go run cmd/generate/main.go -input=static/images/art -output="${GENERATED_DIR}/256" -size=256
	go run cmd/generate/main.go -input=static/images/art -output="${GENERATED_DIR}/512" -size=512
	go run cmd/generate/main.go -input=static/images/art -output="${GENERATED_DIR}/1024" -size=1024

build:
	go build -o ${BUILD_DIR}/${TARGET} -ldflags "-X main.BuildTime=${BUILD_TIME}" cmd/website/main.go

all: clean deps generate build

.PHONY: clean deps generate build all
