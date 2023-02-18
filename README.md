# CLI morty

# Usage 
```bash
morty build --name <name> --runtime <runtime> PATH
```

Example:
```bash
morty build --name test --runtime node-19 --build-arg ADDITIONAL_PACKAGE="iputils curl" --build-arg TARGETPLATFORM="linux/amd64" ./function
```