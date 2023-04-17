# Contributing to Morty CLI

## Compile from source

You only need to install [Go](https://go.dev/doc/install) 1.19+ to compile the project locally. Also, you need to have access to the [polyxia-org/morty-gateway](https://github.com/polyxia-org/morty-gateway) repository with your account.

You need to configure Go to give it access to private Github repositories, as we have a dependency to the [polyxia-org/morty-gateway](https://github.com/polyxia-org/morty-gateway) module. To do that, you need first to export the `GOPRIVATE` environment variable : 

```bash
# You can add it to your .zshrc or .bashrc to avoid export everytime you want to build
export GOPRIVATE="github.com/polyxia-org/morty-gateway"
```

Next, update your `~/.gitconfig` to use SSH to retrieve modules instead of HTTPS : 

```bash
# ~/.gitconfig
[url "ssh://git@github.com/"]
        insteadOf = https://github.com/
```

Now, you can run the following command to compile the project : 

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



