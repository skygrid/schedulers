gofmt -s -r '(a) -> a' -w *.go;
go test -v .;
rm -rf temp;