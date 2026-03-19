# `iocgo`

[English](./README.md) | 简体中文

`iocgo` 是一个面向 Go 的**注解优先、静态代码生成驱动**的 IoC 容器。

它的核心目标不是把复杂度堆进运行时，而是：

- 用 `+ioc:` 注解表达装配意图
- 用生成器把结构性工作前移
- 让运行时容器保持干净、可预测、易测试

## 设计总览

`iocgo` 围绕四个核心理念构建：

- `+ioc:` 注释是一等编程模型
- `ioctl gen` 把注解转成静态注册代码
- 生成的 `init()` 只注册 package/starter registrar，不创建 bean 实例
- 运行时容器只负责依赖解析、作用域管理、生命周期和错误报告

它参考了：

- [IOC-Golang](https://github.com/alibaba/IOC-Golang)
- [controller-tools](https://github.com/kubernetes-sigs/controller-tools)

但它**不会**走运行时包扫描路线。

## 为什么选择静态生成，而不是运行时扫描

这是 `iocgo` 最重要、也不会改变的设计理念之一。

`iocgo` 有意把 IoC 装配中的结构性复杂度前移到代码生成阶段，而不是把这部分复杂度留在运行时容器中。

这样设计的原因很明确：

- 生成器可以更早完成注解解析和装配结构整理
- 运行时容器可以更小、更干净、更确定
- 启动行为更容易解释，也更容易测试
- 很多结构性错误可以更早暴露，而不是拖到运行时

这当然带来一个现实取舍：

- 生成文件可能会比较多

但 `iocgo` 的设计并不要求开发者去阅读和维护这些生成文件。

`iocgo` 期望的开发体验是：

- 开发者理解 `+ioc:` 注解体系
- 生成器负责把注解转成静态注册代码
- 运行时只消费已经生成好的元数据和注册结果

也就是说：

- **人理解注解**
- **机器理解生成文件**

`iocgo` 追求的不是“生成文件最少”，而是：

- 保持运行时干净
- 降低运行时魔法
- 保持规则显式
- 让 Boot 和依赖解析更稳定、更确定

如果把这些能力都放在运行时扫描和反射装配里，短期看起来文件会少一点，但长期会让：

- 容器行为更隐式
- 启动成本更难控制
- 排错成本更高
- 框架边界更模糊

这不是 `iocgo` 想走的路线。

## 安装 CLI

安装 `ioctl`：

```bash
go install github.com/photowey/iocgo/cmd/ioctl@latest
```

或者在当前工作区构建：

```bash
go build -o ioctl ./cmd/ioctl
# 或者
make build
```

验证：

```bash
ioctl version
```

## 支持的注解

### 包级

- `+ioc:package`
- `+ioc:init-register`
- `+ioc:starter`

### 类型级

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

### 字段级

- `+ioc:inject`
- `+ioc:inject:name=<beanName>`

### 函数 / 方法级

- `+ioc:bean`
- `+ioc:bean:name=<beanName>`
- `+ioc:scope=<singleton|prototype>`
- `+ioc:primary`
- `+ioc:lazy`
- `+ioc:expose={TypeA;TypeB}`
- `+ioc:init=<MethodName>`
- `+ioc:destroy=<MethodName>`
- `+ioc:inject:named=paramA:beanA;paramB:beanB`

## 组件风格示例

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

## 配置绑定组件示例

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

当启用 `nemo` starter 时，生成代码会自动：

- 解析 starter binder bean
- 把 `app.feature` 绑定到 `*FeatureConfig`
- 将绑定后的对象注册成普通 singleton bean

## 配置类风格示例

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

## Nemo Starter 示例

如果你希望在 `iocgo` 中使用类似 Spring Environment 的配置运行时，可以启用 first-party `nemo` starter。

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
}
```

也可以在配置类方法里直接依赖 `starterNemo.Environment`：

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

推荐约定：

- 配置感知 bean 方法参数优先使用 `starterNemo.Environment`
- 配置对象绑定优先使用 `+ioc:configuration-properties:prefix=...`
- `+ioc:bind:prefix=...` 仍兼容，但不再是首选公开语义

`nemo` binder tag 也支持配置绑定：

- ``binder:"field"``：字段映射
- ``required:"true"``：缺失时报错
- ``default:"value"``：提供默认值

当前 binder 还支持更丰富的目标类型：

- `time.Duration`
- `[]string`、`[]int`
- 指针字段
- 嵌套指针结构体

## First-Party Starter Demo

完整示例位于：

- [examples/firstpartystarterdemo](./examples/firstpartystarterdemo/demo.go)

它演示了：

- `nemo` starter 启动
- `configuration-properties` 绑定
- `required/default` 语义
- 更丰富的绑定目标
- `application-demo.yml` profile 覆盖
- 配置类 + 绑定组件 + 运行时 bean 获取

## 生成注册代码

```bash
go run ./cmd/ioctl gen register paths=./...
```

生成文件会导出：

- `RegisterIocgoBeans(reg iocgo.Registry) error`

如果启用 `+ioc:init-register`，还会生成包级 `init()`，调用：

- `iocgo.RegisterBeans(...)`

## 公共 Codegen API

`iocgo` 现在也暴露了公共 codegen 包，供外部工具或更高层工具链直接调用：

```go
import publicioc "github.com/photowey/iocgo/codegen/ioc"

func main() {
    if err := publicioc.Run([]string{"register", "paths=./..."}); err != nil {
        panic(err)
    }
}
```

它是以下场景的正式集成点：

- `ioctl`
- `longctl codegen ioc`
- 外部构建或代码生成工作流

## 启动容器

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

## 运行时特性

- singleton / prototype scope
- named / primary bean resolution
- configuration class support
- generated field injection
- generated package `init()` registration
- starter-style registration support
- interface / hook lifecycle support
- shared `Ordered` / `PriorityOrdered` ordering model
- container-owned event core with deterministic listener ordering
- deterministic cycle detection and lookup errors

## License

[Apache-2.0](./LICENSE)
