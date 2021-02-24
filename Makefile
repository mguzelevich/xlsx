.PHONY: run test clean

exec_root = $(shell pwd)
ts = $(shell date +%Y%m%d_%H%M%S)
tag = mgu/xlsxcli:latest

dummy:
	@echo "xlsxcli"

build:
	docker build -t $(tag) .

fmt:
	gofmt -w src/

clean:
	rm -rf .build/*

build-fast: clean
	go build -o ./.build/xlsxcli ./cmd/xlsxcli

# build-all-fast: clean fmt
# 	make src=xlsx output=xlsx build-fast

# build-all: clean fmt
# 	make src=converter output=converter build-docker

input=samples/*.xlsx
mapping=samples/mapping.csv
test-1: build-fast
	.build/xlsxcli --out-prefix=/tmp/xlsxcli-tst- --mapping ${mapping} --mode=one2one ${input}
	# .build/xlsxcli ${input} | jq '.' > /tmp/xlsxcli.json
	# cat /tmp/xlsxcli.out | jq '.'

input=/home/humans.net/git.humans-it.net/_users/mgu/maps/sites/20201006/sites.202009.xlsx
mapping=/home/humans.net/git.humans-it.net/_users/mgu/maps/sites/mapping.csv
output_process=/tmp/process-${ts}
test-sites: build-fast
	go build -o ${output_process} /home/humans.net/git.humans-it.net/_users/mgu/maps/sites/process.go
	.build/xlsxcli --out-prefix=/tmp/xlsxcli-tst- --mapping ${mapping} --mode=one2one ${input} | ${output_process}

test:
	go run cmd/xlsxcli/main.go --mode=one2one --log-level debug samples/test-00.xlsx
