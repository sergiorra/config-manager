<div align="center">
  <h1>Config Manager</h1>
  <blockquote>This config manager is intended to manage an application configuration via a variety of providers
    like defaults, global config files, environment config files, environment variables and flags. It is based
    on the Golang Viper library.</blockquote>
</div>

## 📜 Documentation

### How it works
It splits the configuration in two parts, the Host and the Application.
The Host configuration structure is managed by the Config Manager itself, and contains
all the data related with the machine where the application will run.
It currently only contains the Environment variable.
The application configuration is managed by each client and is passed to the Config Manager.
This configuration contains all the data that the app itself will use,
like logging, database connections, other services information and so on.

It applies a specific descending priority for the different providers.

For the Host configuration:
1. Defaults
2. Environment Variables
3. Flags

For the Application configuration:
1. Defaults
2. Global config file
3. Environment config file
4. Environment Variables
5. Flags

### Prerequisites

The Config Manager has some prerequisites in order to work as expected:

- Folder `config` on the root of the project
- The global config file will be named `config.json` and will be inside the `config` folder. It is a mandatory file.
- The environment config files will follow the name structure `config.[Env].json` and will be inside the `config` folder. It is not a mandatory file.
- Only one level of nested structs will be supported for the Application configuration
- There can not be a root key variable named `Env` in the Application config. If found, it will be
  a readonly variable for the `Env` Host config variable.
- Only support for the types `Struct`, `Int`, `String` and `Bool` in the Application config.
- The Application config passed to the Config Manager must be a pointer value.
- The keys for the default values defined in the constructor will follow the structure (all in lowercase):
    - `[field_name]` for the root variables
    - `[field_name].[inner_field_name]` for the nested variables
- The environment variables will follow the structure (all in uppercase):
    - `[field_name]` for the root variables
    - `[field_name]_[inner_field_name]` for the nested variables
    - In case to add the environment variable prefix in the constructor will follow the structure:
        - `[envVariablesPrefix]_[field_name]` for the root variables
        - `[envVariablesPrefix]_[field_name]_[inner_field_name]` for the nested variables
- The flags passed will follow the structure (all in lowercase):
    - `--[field_name] [value]` for the root variables
    - `--[field_name].[inner_field_name] [value]` for the nested variables

### More information

You can check how Viper works and all its documentation in the [Viper](https://github.com/spf13/viper) repository

## ⚙️ Usage and examples

### Simple load config

```go
cfg := &Configuration{}

mgr := config.NewManager()

err := mgr.Load(cfg)
if err != nil {
    log.Fatal(err)
}
```

### Load config with default values

```go
cfg := &Configuration{}

mgr := config.NewManager(config.WithDefault("rpchost", 13000), config.WithDefault("logger.loglevel", "WARN"))

err := mgr.Load(cfg)
if err != nil {
     log.Fatal(err)
}
```

### Load config with default values and environment variables prefix

```go
cfg := &Configuration{}

mgr := config.NewManager(config.WithEnvPrefix("OFFCTRL"), config.WithDefault("rpchost", 13000))

err := mgr.Load(cfg)
if err != nil {
     log.Fatal(err)
}
```

### Json config file example
```json
{
  "Logger": {
    "UseSysLog": false,
    "LogLevel": "INFO",
    "LogFormat": "plain"
  },
  "RpcHost": "localhost",
  "RpcPort": 13001,
  "RpcTimeoutConn": 15
}
```

### Environment variables declaration example
```shell
// without prefix
$ export ENV="PROD"
$ export LOGGER_LOGLEVEL="WARN"
$ export RPCPORT="8080"

// with prefix
$ export OFFCTRL_ENV="PROD"
$ export OFFCTRL_LOGGER_LOGLEVEL="WARN"
$ export OFFCTRL_RPCPORT="8080"
```

### Flags usage example
```shell
$ go run main.go --env PROD --logger.loglevel WARN --rpcport 8080
```

## 📦 Installation

```
$ go get -u github.com/sergiorra/config-manager
```
