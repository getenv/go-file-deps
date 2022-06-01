This tool prints the list of files required to build a Go package.

    $ go install github.com/getenv/go-file-deps@latest
    $ go-file-deps ./cmd/hello

## Makefile integration

Pass `-r rule` to print `rule: $(wildcard files...)\n` instead.

```makefile
hello:
    go build -o $@ ./cmd/hello
    go-file-deps -r $@ ./cmd/hello > $@.d

-include hello.d
```
