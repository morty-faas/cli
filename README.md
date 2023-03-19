# CLI morty

Morty CLI is a tool that helps you to create, build and deploy your function.

# Usage

## Configure your context

Morty CLI has the ability to work with multiple contexts. By default, Morty will run with the following context : 

```
Name         : localhost
Gateway URL  : http://localhost:8080
Registry URL : http://localhost:8081
```

To add your own context, use the following command : 

```bash
morty config add-context $CONTEXT_NAME --gateway=$GATEWAY --registry=$REGISTRY
```

Replace `$CONTEXT_NAME`, `$GATEWAY`, `$REGISTRY` with your own values. 

To see all the contexts available in your configuration, use the following command: 

```bash
morty config contexts
```

You can view the information about your current context by running the following command : 

```bash
morty config current-context
```

By default, Morty CLI configuration is stored in `~/.morty/config.yaml`. If you want, you can use a configuration file at another location by exporting the following environment variable : `export MORTYCONFIG=/path/to/your/config.yaml`.

To switch to a different context, run the following command : 

```bash
morty config use-context $CONTEXT_NAME
```

## Create a new function

```bash
morty new --runtime <runtime> FUNCTION_NAME
```

_This command creates a new workspace and default configuration for your function._

Example:

```bash
morty new --runtime node-19 myFirstFuntion
```

## Package the function

_This command will package your function and create a lz4 file._

If MORTY_REGISTRY_URL is set, the function will be uploaded to the registry.

```bash
export MORTY_REGISTRY_URL=<registry_url>
sudo morty build <name> --runtime <runtime> --path <path_to_function> --build-arg [build_arg]
```

Example:

```bash
export MORTY_REGISTRY_URL="http://localhost:8080"
sudo morty build test --runtime node-19 --build-arg ADDITIONAL_PACKAGE="iputils curl" --build-arg TARGETPLATFORM="linux/amd64" --path ./function
```

**Care about the `http(s)://` prefix in the registry URL**

**This command should be run with root privileges**

or in an more intuitive way in spcifying just the name of the function, and he will go in the workspace folder and build it:

```bash
export MORTY_REGISTRY_URL="http://localhost:8080"
sudo morty build test
```
