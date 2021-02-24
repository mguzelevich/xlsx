.PHONY: run test clean

exec_root = $(shell pwd)
ts = $(shell date +%Y%m%d_%H%M%S)

dummy:
	@echo "xlsxcli"

build-builder:
	docker build -t mguzelevich/golang -f Dockerfile-golang .

fmt:
	gofmt -w src/

clean:
	rm -rf .build/*

# build: clean fmt
# #	go build -o build/staff-tools .
# 	time docker run -it --rm \
# 		-v "$(exec_root)/src:/src" \
# 		-v "$(exec_root)/.build:/.build" \
# 		-v "$(exec_root)/.build/.cache:/root/go" \
# 		mguzelevich/golang \
# 		sh -c "cd /src; go build -v -o ../.build/xlsxcli ./cmd/xlsxcli/"

# CGO_ENABLED = CGO_ENABLED=1
# build_tags = -tags 'sqlite3 json1 sqlite_json'
# ld_flags = -ldflags '-extldflags -static'
# build-docker:
# 	time docker run -it --rm \
# 		-v "$(exec_root)/src:/src" \
# 		-v "$(exec_root)/.build:/.build" \
# 		-v "$(exec_root)/.build/.cache:/root/go" \
# 		mguzelevich/golang \
# 		sh -c "cd /src; ${CGO_ENABLED} ${GOOS} ${GOARCH} go build -v ${build_tags} ${ld_flags} -o ../.build/$(output) ./cmd/$(src)/"

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
	go run cmd/xlsxcli/main.go --log-level trace samples/test-00.xlsx
