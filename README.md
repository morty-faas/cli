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

## Logger

A logger is embedded in the code. You need to activate it to be able to see the logs.

To enable it, simply export the following environment variable with the level you want :

```bash
export MORTY_LOG=debug
```

You should now see logs :

```bash
morty config current-context

# Output without MORTY_LOG set
Name         : thomas-dev
Gateway URL  : http://162.38.112.57:8080
Registry URL : http://162.38.112.57:8081

# Output with MORTY_LOG
INFO[0000] Loading configuration from path: /Users/thomas/.morty/config.yaml
INFO[0000] Active context : thomas-dev
Name         : thomas-dev
Gateway URL  : http://162.38.112.57:8080
Registry URL : http://162.38.112.57:8081
```

## Create a new function

To create a new function, also called a **workspace**, you can use the following command:

```bash
morty function init $FUNCTION_NAME <opt:$DIRECTORY>
# or alias
morty fn init $FUNCTION_NAME <opt:$DIRECTORY>
```

> By default, if the flag `--runtime` is not specified, the function will be created using the `node-19` runtime.

Replace `$FUNCTION_NAME` with your own value. You can also specify an additional argument `$DIRECTORY` to control where the function will be created.

This command will create for you the skeleton of the function to enable fast development. Please note that the **morty.yaml** file inside the function directory is mandatory and it contains function metadata used during the build process.

## Build the function

_This command will send your function to the registry, and create it inside the system. It will rely on the `morty.yaml` file in the function directory to define function name and runtime._

```bash
morty function build $FUNCTION_DIRECTORY
```
