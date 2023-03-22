# Morty CLI

The Morty CLI is an interface allowing developers or users to easily interact with one or more Morty serverless instances.

## Quick start

First, you need to install the Morty CLI using one of the following methods:

- Download a pre-compiled binary from the [releases](https://github.com/polyxia-org/morty-cli/releases) page
- Build it from source, please see the [CONTRIBUTING.md](./CONTRIBUTING.md#compile-from-source)

Once you have your Morty CLI installed locally (for the rest of the commands, we assume that `morty` is available in your `$PATH`), you can list the available runtimes and choose the more appropriate for your first function : 

```bash
morty runtime ls
# Or the short hand syntax
morty rt ls
# Output (can be different depending on when you execute it)
Available runtimes:
- go-1.19
- node-19
- python-3
- rust-1.67
```

Here we will use the `node-19` runtime to create our first function : 

```bash
morty function init --runtime node-19 hello-morty
# Or the shorthand syntax
morty fn init -r node-19 hello-morty
```

You should now have a directory `hello-morty` in your current directory that contains a function code example. 

Once your function is ready, you will need to build it before invocating it. To do that, you need to have access to a Morty FaaS instance. Please refer to the [General Documentation(TODO)](#) to learn how to run a local Morty FaaS instance. For the sake of simplicity here, we will assume that we have a Morty Gateway running on `https://faas.morty.io` and a Morty Registry running on `https://registry.morty.io`, but you can use your own values.

You need first to configure your context: 

```bash
morty config add-context morty-faas --gateway https://faas.morty.io --registry https://registry.morty.io
# Output
Success ! Your context 'morty-faas' has been saved and it is now the active context.
```

Every requests will now be performed against this Morty FaaS instance.

You can now build your function (it can take some time to complete) :

```bash
morty fn build hello-morty
# Once done, you should see something like : 
Function hello-morty has been created !
```

Finally, you can invoke your function using the following command : 
```bash
morty fn invoke hello-morty
```

Et voilà ! You've created your first function in Morty !

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
