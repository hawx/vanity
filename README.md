# vanity

Server for golang vanity import paths.

``` bash
$ go get hawx.me/code/vanity
$ cat vanity.conf
/example git git://github.com/example/example
/something/else git git://github.com/org/else
$ vanity my.cool.domain ./vanity.conf
...
```

Each configuration line is the data (minus host) that is returned in the meta
tag, see `go help importpath` for more information. Any pages hit by humans
redirect to `godoc.org`.


## Prior art

- [gimpy](github.com/nf/gimpy), which uses DNS to get projects instead of a
  config file.
