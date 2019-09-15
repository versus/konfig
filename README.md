[![Go Doc][godoc-image]][godoc-url]
[![Build Status][workflow-image]][workflow-url]
[![Go Report Card][goreport-image]][goreport-url]
[![Test Coverage][coverage-image]][coverage-url]
[![Maintainability][maintainability-image]][maintainability-url]

# konfig

This is a very minimal and unopinionated utility for reading configuration values
in Go applications based on [The 12-Factor App](https://12factor.net/config).

The idea is that you just define a `struct` with _fields_, _types_, and _defaults_,
then you just use this library to read and populate the values for fields of your struct
from either **command-line flags**, **environment variables**, or **configuration files**.
It can also watch for new values read from _configuration files_ and notify subscribers.

This library does not use `flag` package for parsing flags, so you can still parse your flags separately.

## Quick Start

```go
package main

import (
  "fmt"
  "net/url"
  "time"

  "github.com/moorara/konfig"
)

var Config = struct {
  Enabled   bool
  LogLevel  string
  Timeout   time.Duration
  Address   url.URL
  Endpoints []string
} {
  Enabled:  true,   // default
  LogLevel: "info", // default
}

func main() {
  konfig.Pick(&Config)
  fmt.Printf("%+v\n", Config)
}
```

## Documentation

The precedence of sources for reading values is as follows:

  1. command-line flags
  2. environment variables
  3. configuration files
  4. default values (set when creating the instance)

You can pass the configuration values with **flags** using any of the syntaxes below:

```bash
main  -enabled  -log.level info  -timeout 30s  -address http://localhost:8080  -endpoints url1,url2,url3
main  -enabled  -log.level=info  -timeout=30s  -address=http://localhost:8080  -endpoints=url1,url2,url3
main --enabled --log.level info --timeout 30s --address http://localhost:8080 --endpoints url1,url2,url3
main --enabled --log.level=info --timeout=30s --address=http://localhost:8080 --endpoints=url1,url2,url3
```

You can pass the configuration values using **environment variables** as follows:

```bash
export ENABLED=true
export LOG_LEVEL=info
export TIMEOUT=30s
export ADDRESS=http://localhost:8080
export ENDPOINTS=url1,url2,url3
```

You can also write the configuration values in **files** (or mount your configuration values and secrets as files)
and pass the paths to the files using environment variables:

```bash
export ENABLED_FILE=...
export LOG_LEVEL_FILE=...
export TIMEOUT_FILE=...
export ADDRESS_FILE=...
export ENDPOINTS_FILE=...
```

### Skipping

If you want to skip a source for reading values, use `-` as follows:

```go
type Config struct {
  GithubToken string `env:"-" fileenv:"-"`
}
```

In the example above, `GithubToken` can only be set using `github.token` command-line flag.

### Customization

You can use Go _struct tags_ to customize the name of expected command-line flags or environment variables.

```go
type Config struct {
  Database string `flag:"config.database" env:"CONFIG_DATABASE" fileenv:"CONFIG_DATABASE_FILE_PATH"`
}
```

In the example above, `Database` will be read from either:

  1. The command-line flag `config.databas`
  2. The environment variable `CONFIG_DATABASE`
  3. The file specified by environment variable `CONFIG_DATABASE_FILE_PATH`
  4. The default value set on struct instance

### Using flag Package

`konfig` plays nice with `flag` package since it does NOT use `flag` package for parsing command-line flags.
That means you can define, parse, and use your own flags using built-in `flag` package.

If you use `flag` package, `konfig` will also add the command-line flags it is expecting.
Here is an example:

```go
package main

import (
  "flag"
  "time"

  "github.com/moorara/konfig"
)

var Config = struct {
  Enabled   bool
  LogLevel  string
} {
  Enabled:  true,   // default
  LogLevel: "info", // default
}

func main() {
  konfig.Pick(&Config)
  flag.Parse()
}
```

If you run this example with `-help` or `--help` flag,
you will see `-enabled` and `-log.level` flags are also added with descriptions!

### Options

You can pass a list of options to `Pick`.
These options are helpers for specific situations and setups.

| Option | Description | Example |
|--------|-------------|---------|
| `konfig.Debug()` | Printing debugging information | [example](./examples/2-debug), [blog](https://milad.dev/projects/konfig/#debugging) |
| `konfig.PrefixFlag()` | Prefixing all flag names with a string | |
| `konfig.PrefixEnv()` | Prefixing all environment variable names with a string | |
| `konfig.PrefixFileEnv()` | Prefixing all file environment variable names with a string | |
| `konfig.ListSep()` | Specifying list separator for all fields with slice type | |
| `konfig.SkipFlag()` | Skipping command-line flags as a source for all fields | |
| `konfig.SkipEnv()` | Skipping environment variables as a source for all fields | |
| `konfig.SkipFileEnv()` | Skipping file environment variables (and configuration files) as a source for all fields | |
| `konfig.Telepresence()` | Reading configuration files in a _Telepresence_ environment | [blog](https://milad.dev/posts/telepresence-with-konfig) |
| `konfig.WatchInterval()` | Overriding default interval for `Watch()` method | [example](./examples/3-watch), [blog](https://milad.dev/posts/dynamic-config-secret) |

### Debugging

If for any reason the configuration values are not read as you expected, you can view the debugging logs.
You can enable debugging logs either by using `Debug` option or by setting `KONFIG_DEBUG` environment variable.
In both cases you need to specify a verbosity level for logs.

### Watching Changes

konfig allows you to watch _configuration files_ and dynamically update your configurations as your application is running.

When using `Watch()` method, your struct should have a `sync.Mutex` field on it for synchronization and preventing data races.
You can find an example of using `Watch()` method [here](./examples/3-watch).

[Here](https://milad.dev/posts/dynamic-config-secret) you will find a real-world example of using `konfig.Watch()`
for **dynamic configuration management** and **secret injection** for Go applications running in Kubernetes.


[godoc-url]: https://godoc.org/github.com/moorara/konfig
[godoc-image]: https://godoc.org/github.com/moorara/konfig?status.svg
[workflow-url]: https://github.com/moorara/konfig/actions
[workflow-image]: https://github.com/moorara/konfig/workflows/Main/badge.svg
[goreport-url]: https://goreportcard.com/report/github.com/moorara/konfig
[goreport-image]: https://goreportcard.com/badge/github.com/moorara/konfig
[coverage-url]: https://codeclimate.com/github/moorara/konfig/test_coverage
[coverage-image]: https://api.codeclimate.com/v1/badges/207d8984d6121c11d9d6/test_coverage
[maintainability-url]: https://codeclimate.com/github/moorara/konfig/maintainability
[maintainability-image]: https://api.codeclimate.com/v1/badges/207d8984d6121c11d9d6/maintainability
