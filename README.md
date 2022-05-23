# `iocgo`

An `IOC ` container in `Golang`

> This project is inspired by the [IOC-Golang](https://github.com/alibaba/IOC-Golang)
> and `K8S` [controller-tools](https://github.com/kubernetes-sigs/controller-tools) projects

## `Design ideas`

1.`BeanDefinition`

- `bean name`

- `scope`

- `autowire enabled`

- `configuration struct`

- `factory func`

- `init func`

- `destroy func`

- `description`

- ```go
  // +ioc:autowired:beandefinition  // mark a beandefinition
  
  type DumpyBeanDefinition struct {
  	Name              string      // bean name
  	Scope             scope.Scope // scope of bean
  	AutowireCandidate bool        // autowire enabled
  	Configuration     string      // configuration struct
  	FactoryFunc       string      // factory func
  	InitFunc          string      // init func
  	DestroyFunc       string      // destroy func
  	Description       string      // description of bean
  }
  ```

-

2.`BeanFactory`

- `Register`
- `GetBean`
- `Destroy`

3.`LifeCycle`

- `InitializingBean`
  - `AfterPropertiesSet()`
- `DisposableBean`
  - `Destroy()`

4.`Configuration`

```go
// +ioc:autowired:factorybean // mark a bean is a factory bean

type Person struct {
ID   int64  `json:"id"`
Name string `json:"name"`
Age  uint8  `json:"age"`
}

// +ioc:autowired:configuration // mark a configuration struct -> @Configuration
type PersonConfiguration struct {
}

// +ioc:autowired:factoryfunc // mark a func is configuration factory func -> @Bean
// +ioc:autowired:scope=singleton // mark a bean scope is singleton

func (config *PersonConfiguration) CreatePerson(ctx context.Context) Person {
return Person{
ID:   9527,
Name: "photowey",
Age:  18,
}
}

```

