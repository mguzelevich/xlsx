# xlsx-tools

## build

```
export $(egrep -v ‘^#’ .env | xargs) && docker build --build-arg GO_IMPORT_TOKEN=${GITHUB_TOKEN} -t go-private-example .
```

## modes

- `one2many` - 1 csv for every xlsx sheet
- `one2one` - 1 csv for all xlsx sheets

## examples

```
$ xlsxcli -mode one2many -mapping mapping.csv input.xlsx
```

## storage

key-value
key - tupple of strings

access modes
- key - value
- key - table

## Mapping file

```
ID,#,№,Numer,N
Name,name,username
Surname,SecondName,second
```

## Makefile

```
.PHONY: run test clean

exec_root = $(shell pwd)
ts = $(shell date +%Y%m%d_%H%M%S)
short_ts = $(shell date +%Y%m%d)

XLSXCLI_IMAGE = mguzelevich/xlsxcli

prepare:
	echo "prepare"

run: input_file = input.xlsx
run: output_file = output.csv
run: mapping_file = mapping.csv
run: prepare
	docker run -it --rm \
      -v "${exec_root}":/project \
	  ${XLSXCLI_IMAGE} \
	    sh -c '/xlsx --mapping "/project/${mapping_file}" --mode=one2one "/project/${input_file}" > /project/${output_file}'
```

## usage

```
$ make run
```