# Contributing to Morty CLI

## Compile from source

You only need to install [Go](https://go.dev/doc/install) 1.19+ to compile the project locally.

Once you have Go ready, you can run the following command to compile the project : 

```bash
go build -o morty main.go
```

You should now be able to use your CLI :

```bash
./morty config current-context

# Output
Name         : localhost
Gateway URL  : http://localhost:8080
Registry URL : http://localhost:8081
```



