# CLI morty

Morty cli is a tool to help you to create, build and deploy your function.

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
export MINIO_ENDPOINT="<MINIO_HOST>" MINIO_USER="<MINIO_USER>" MINIO_PASSWORD="<MINIO_PASSWORD>"
morty build --name <name> --runtime <runtime> PATH
```

Example:

```bash
export MINIO_ENDPOINT="localhost:9000" MINIO_USER="minioadmin" MINIO_PASSWORD="minioadmin"
morty build --name test --runtime node-19 --build-arg ADDITIONAL_PACKAGE="iputils curl" --build-arg TARGETPLATFORM="linux/amd64" ./function
```
