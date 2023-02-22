# CLI morty

# Usage

## Create a new function

```bash
morty new --runtime <runtime> FUNCTION_NAME
```

_This command create a new workspace and default configuration for you're function._
The function workspace repository should be available at [./workspace](./workspace)

Example:

```bash
morty new --runtime node-19 myFirstFuntion
```

## Package the function

```bash
morty build --name <name> --runtime <runtime> PATH
```

Example:

```bash
morty build --name test --runtime node-19 --build-arg ADDITIONAL_PACKAGE="iputils curl" --build-arg TARGETPLATFORM="linux/amd64" ./function
```
