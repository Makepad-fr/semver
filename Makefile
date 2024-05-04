OUTPUT_DIR?=./out

.PHONY: build
build: create-output-dir
	go build -o ${OUTPUT_DIR}/semver ./cli/main.go

.PHONY: create-output-dir
create-output-dir:
	mkdir -p ${OUTPUT_DIR}

.PHONY: clean
clean:
	rm -rf ${OUTPUT_DIR}