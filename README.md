# CLI morty

Morty CLI is a tool that helps you to create, build and deploy your function.

# Usage

## Create a new function

```bash
morty new --runtime <runtime> FUNCTION_NAME
```

_This command creates a new workspace and default configuration for your function._

The function workspace repository should be available at [./workspace](./workspace)

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