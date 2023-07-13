# SQL-powered Web-server Demo

You can run the server as follows:

```
PORT=8080 go run .
```

Alternately, build it and run it as separate steps:

```
go build .
export PORT=8080
./unboxed
```

Setting the PORT env. var. is optional; the default is 8080.

You can cross-compile the binary for Windows (on x64) like this:

```
env GOOS=windows GOARCH=amd64 go build
```

That will produce an unboxed.exe file.
