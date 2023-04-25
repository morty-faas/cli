# Contributing to Morty CLI

## Compile from source

You only need to install [Go](https://go.dev/doc/install) 1.19+ to compile the project locally. Also, you need to have access to the [morty-faas/controller](https://github.com/morty-faas/controller) repository with your account.

Now, you can run the following command to compile the project :

```bash
go build -o morty main.go
```

You should now be able to use your CLI :

```bash
./morty config current-context

# Output
Name         : localhost
Controller URL  : http://localhost:8080
Registry URL : http://localhost:8081
```
