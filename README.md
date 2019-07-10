# konfig

This is a very minimal and unopinionated utility for reading configuration values
in Go applications based on [The 12-Factor App](https://12factor.net/config).

The idea is that you just define a `struct` with _fields_, _types_, and _defaults_,
then you just use this library to read and populate the values for fields of your struct
from either **command-line flags**, **environment variables**, or **configuration files**.

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

### Debugging

If you are not sure how your configuration values are going to be read,
or they are not read as expected, you can use `PickAndLog` for debugging purposes.
