# `iocgo`

English | [简体中文](./README.zh-CN.md)

An annotation-first IoC container and static code generator for Go.

## Design Summary

`iocgo` is built around four ideas:

- `+ioc:` comments are first-class annotations
- `ioctl gen` turns annotations into static registration code
- generated `init()` code registers package or starter registrars
- the runtime container resolves dependencies, manages scope, and runs lifecycle hooks

This project is inspired by [IOC-Golang](https://github.com/alibaba/IOC-Golang) and [controller-tools](https://github.com/kubernetes-sigs/controller-tools), but it does not rely on runtime package scanning.

## Why Static Generation Instead of Runtime Scanning

This is a deliberate core design choice of `iocgo`, and it is not expected to change.

`iocgo` intentionally pushes container assembly work into the code generator instead of keeping that complexity inside the runtime container.

The reasoning is simple:

- the generator can do the expensive structural work earlier
- the runtime can stay smaller, cleaner, and more deterministic
- startup behavior becomes easier to explain and easier to test
- many structural mistakes can be surfaced earlier instead of being deferred to runtime

This also means `iocgo` accepts a tradeoff:

- it may generate many files

But those generated files are not the primary human-facing programming model.

The intended developer experience is:

- developers understand the `+ioc:` annotation system
- developers run the generator
- the runtime consumes generated metadata and registration code

In other words, `iocgo` is designed so that:

- humans understand the annotations
- machines understand the generated files

The framework does not optimize for "minimal generated output". It optimizes for:

- keeping runtime behavior clean
- keeping container rules explicit
- reducing runtime magic
- preserving deterministic boot and resolution behavior

That is why `iocgo` is not designed as a runtime-scanning IoC framework even though that would reduce generated files on disk.

## Install CLI

Install the `ioctl` generator locally:

```bash
go install github.com/photowey/iocgo/cmd/ioctl@latest
```

Or build it from the current workspace:

```bash
go build -o ioctl ./cmd/ioctl
```

Then verify:

```bash
ioctl version
```

## Supported Annotations

Package level:

- `+ioc:package`
- `+ioc:init-register`
- `+ioc:starter`

Type level:

- `+ioc:component`
- `+ioc:component:name=<beanName>`
- `+ioc:configuration`
- `+ioc:scope=<singleton|prototype>`
- `+ioc:primary`
- `+ioc:lazy`
- `+ioc:expose={TypeA;TypeB}`
- `+ioc:configuration-properties`
- `+ioc:configuration-properties:prefix=<config.path>`
- `+ioc:constructor=<FuncName>`

Field level:

- `+ioc:inject`
- `+ioc:inject:name=<beanName>`

Function and method level:

- `+ioc:bean`
- `+ioc:bean:name=<beanName>`
- `+ioc:scope=<singleton|prototype>`
- `+ioc:primary`
- `+ioc:lazy`
- `+ioc:expose={TypeA;TypeB}`
- `+ioc:init=<MethodName>`
- `+ioc:destroy=<MethodName>`
- `+ioc:inject:named=paramA:beanA;paramB:beanB`

## Component Style Example

```go
// +ioc:package
// +ioc:init-register
package user

// +ioc:component
// +ioc:component:name=userRepo
// +ioc:scope=singleton
type UserRepository struct{}

// +ioc:component
// +ioc:component:name=userService
// +ioc:scope=singleton
type UserService struct {
    // +ioc:inject
    Repo *UserRepository
}
```

## Bound Configuration Component Example

```go
// +ioc:package
// +ioc:init-register
package app

// +ioc:component
// +ioc:component:name=featureConfig
// +ioc:configuration-properties
// +ioc:configuration-properties:prefix=app.feature
type FeatureConfig struct {
    Enabled bool `binder:"enabled" required:"true"`
    Port    int  `binder:"port" default:"8080"`
}
```

When the `nemo` starter is active, generated code will:

- resolve the starter binder bean
- bind `app.feature` into `*FeatureConfig`
- register the bound config object as a normal singleton bean

## Configuration Style Example

```go
package person

import "context"

// +ioc:configuration
type PersonConfiguration struct{}

type Person struct {
    ID   int64
    Name string
    Age  uint8
}

// +ioc:bean
// +ioc:bean:name=person
// +ioc:scope=singleton
func (c *PersonConfiguration) CreatePerson(ctx context.Context) *Person {
    return &Person{
        ID:   9527,
        Name: "photowey",
        Age:  18,
    }
}
```

## Nemo Starter Example

Use the first-party `nemo` starter when you want a Spring Environment style configuration runtime inside `iocgo`.

```go
package main

import (
    "context"

    "github.com/photowey/iocgo"
    starterNemo "github.com/photowey/iocgo/pkg/starters/nemo"
)

func main() {
    starterNemo.Configure(
        starterNemo.WithSearchPaths("configs"),
        starterNemo.WithProfiles("dev"),
    )

    app := iocgo.New()
    if err := app.Boot(context.Background()); err != nil {
        panic(err)
    }

    env := iocgo.MustGet[starterNemo.Environment](context.Background(), app)
    _, _ = env.Get("nemo.application.name")

    cfg := starterNemo.MustBind[struct {
        Name string `binder:"name" default:"demo"`
    }](env, "app.info")
    _ = cfg
}
```

You can also depend on `nemo.Environment` directly in a configuration bean method:

```go
// +ioc:package
// +ioc:init-register
package app

import (
    "context"

    starterNemo "github.com/photowey/iocgo/pkg/starters/nemo"
)

// +ioc:configuration
// +ioc:configuration-properties
// +ioc:configuration-properties:prefix=app.info
type AppConfiguration struct{}

type AppInfo struct {
    Name string
}

// +ioc:bean
// +ioc:bean:name=appInfo
func (c *AppConfiguration) CreateAppInfo(ctx context.Context, env starterNemo.Environment) *AppInfo {
    value, _ := env.Get("nemo.application.name")
    name, _ := value.(string)
    return &AppInfo{Name: name}
}
```

This is the recommended convention for configuration-aware bean factories:

- use `starterNemo.Environment` as the parameter type
- let the generated code resolve it from the starter-registered singleton
- keep configuration access inside `+ioc:configuration` bean methods
- use `+ioc:configuration-properties:prefix=...` on config structs or configuration classes when you want automatic binding

`+ioc:bind:prefix=...` is still supported as a compatibility alias, but `+ioc:configuration-properties:prefix=...` is now the preferred public syntax.

`nemo` binder tags can also participate in configuration binding:

- ``binder:"field"`` maps the property key
- ``required:"true"`` makes the property mandatory
- ``default:"value"`` provides a fallback value when the property is missing

The current binder also supports richer target types such as:

- `time.Duration`
- slices like `[]string` and `[]int`
- pointer fields like `*int`
- nested pointer structs

For typed binding, you can use the helper API directly:

```go
cfg, err := starterNemo.Bind[FeatureConfig](env, "app.feature")
if err != nil {
    panic(err)
}
```

Or resolve the starter binder bean:

```go
binder := iocgo.MustGet[*starterNemo.Binder](context.Background(), app, starterNemo.BinderBeanName)
cfg := FeatureConfig{}
if err := binder.Bind("app.feature", &cfg); err != nil {
    panic(err)
}
```

For diagnostics, you can inspect wrapped configuration bind errors:

```go
if cfgErr, ok := starterNemo.AsConfigurationBindError(err); ok {
    _ = cfgErr.BeanName
    _ = cfgErr.Prefix
}
```

Or render a terminal-friendly diagnostic string:

```go
fmt.Println(starterNemo.FormatConfigurationBindDiagnostic(err))
```

## First-Party Starter Demo

A runnable demo package is available at:

- [examples/firstpartystarterdemo](/D:/workws/gopath/src/github.com/photowey/iocgo/examples/firstpartystarterdemo/demo.go)

It demonstrates:

- `nemo` starter bootstrap
- `configuration-properties` binding
- default and required binder semantics
- richer binding targets such as defaults, pointers, and nested config objects
- profile-specific file override with `application-demo.yml`
- configuration bean + bound component + runtime bean retrieval

## Generate Registrars

```bash
go run ./cmd/ioctl gen register paths=./...
```

Generated files emit `RegisterIocgoBeans(reg iocgo.Registry) error` and, when `+ioc:init-register` is enabled, a package `init()` that calls `iocgo.RegisterBeans(...)`.

## Public Codegen API

`iocgo` now exposes a public codegen package for external callers and higher-level toolchains:

```go
import publicioc "github.com/photowey/iocgo/codegen/ioc"

func main() {
    if err := publicioc.Run([]string{"register", "paths=./..."}); err != nil {
        panic(err)
    }
}
```

This is the intended integration point for:

- `ioctl`
- `longctl codegen ioc`
- external build or code generation workflows

## Boot the Container

```go
package main

import (
    "context"

    "github.com/photowey/iocgo"
    _ "github.com/your/module/person"
)

func main() {
    app := iocgo.New()
    if err := app.Boot(context.Background()); err != nil {
        panic(err)
    }

    person := iocgo.MustGet[*Person](context.Background(), app, "person")
    _ = person
}
```

## Runtime Features

- singleton and prototype scope
- named and primary bean resolution
- configuration class support
- generated field injection
- generated package `init()` registration
- starter-style registration support
- lifecycle callbacks through interfaces and generated hooks
- shared `Ordered` / `PriorityOrdered` ordering model
- container-owned event core with deterministic listener ordering
- deterministic cycle detection and lookup errors

## License

[Apache-2.0](./LICENSE)