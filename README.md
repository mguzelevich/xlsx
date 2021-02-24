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