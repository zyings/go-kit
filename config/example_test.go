package config_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/hpifu/go-kit/config"
)

// 场景一：直接使用 config 的全局 Get 方法
// 可读性：代码最简单，无需提前声明，写起来最方便，可读性也还可以
// 复用性：配置项是写死的，这样的模块很难复用
// 可测试性：需要调用 config 的 Set 方法，mock 使用到的配置项，mock 的代价较高
// 可维护性：配置项散落在代码中，新增和修改都不太方便
// 安全性：可使用 GetD，GetE 之类的方法来做一些错误处理，关联的多个配置项无法保证原子性（这种场景触发几率较低，目前还未碰到过）
func TestExample1(t *testing.T) {
	// package main
	if err := config.Init("testfile/base.json"); err != nil {
		panic(err)
	}
	if err := config.Watch(); err != nil {
		panic(err)
	}
	defer config.Stop()

	// package module
	fmt.Println(config.GetString("OSS.AccessKeyID"))
	fmt.Println(config.GetString("OSS.AccessKeySecret"))
	fmt.Println(config.GetString("OSS.Endpoint"))

	// package test
	_ = config.UnsafeSet("OSS.AccessKeyID", "test-ak")
	_ = config.UnsafeSet("OSS.AccessKeyID", "test-sk")
	_ = config.UnsafeSet("OSS.AccessKeyID", "endpoint")
}

// 场景二：使用全局 config 的 bind 类型方法，类 flag 的使用方式
// 可读性：代码比较简单，类似 flag 的使用方法，需提前将变量绑定一个 key 上，使用时直接使用变量，可读性较好
// 复用性：配置项依然是写死的，复用性较差
// 可测试性：可重新赋值配置项，相比较调用 config 的 Set 方法，稍微好一点点
// 可维护性：配置项可集中在模块中声明，维护起来稍微方便一些
// 安全性：可以在检查声明中增加失败回调保证安全，关联的多个配置项同样无法保证原子性
func TestExample2(t *testing.T) {
	// package main
	if err := config.Init("testfile/base.json"); err != nil {
		panic(err)
	}
	if err := config.Watch(); err != nil {
		panic(err)
	}
	defer config.Stop()

	// package module
	var AccessKeyID = config.String("OSS.AccessKeyID", config.OnFail(func(err error) {
		fmt.Println(err)
	}))
	var AccessKeySecret = config.String("OSS.AccessKeySecret")
	var Endpoint = config.String("OSS.Endpoint")

	fmt.Println(AccessKeyID.Get())
	fmt.Println(AccessKeySecret.Get())
	fmt.Println(Endpoint.Get())

	// package test
	AccessKeyID = config.NewAtomicString("test-ak")
	AccessKeySecret = config.NewAtomicString("test-sk")
	Endpoint = config.NewAtomicString("endpoint")
}

// 场景三: 将配置项定义成一个结构体，使用全局 config 的 bind 对象方法
// 可读性：可读性较好
// 复用性：配置项依然是写死的，复用性较差
// 可测试性：需要手动设置 option 中的各个配置项，和场景二差不多
// 可维护性：配置项可集中在一个结构体中，和场景二差不多
// 安全性：可保证关联配置项的原子性
func TestExample3(t *testing.T) {
	// package main
	if err := config.Init("testfile/base.json"); err != nil {
		panic(err)
	}
	if err := config.Watch(); err != nil {
		panic(err)
	}
	defer config.Stop()

	// package module
	type ModOption struct {
		AccessKeyID     string
		AccessKeySecret string
		Endpoint        string
	}

	var option = config.Bind("OSS", ModOption{})
	fmt.Println(option.Load().(ModOption).AccessKeyID)
	fmt.Println(option.Load().(ModOption).AccessKeySecret)
	fmt.Println(option.Load().(ModOption).Endpoint)

	// package test
	option = &atomic.Value{}
	option.Store(ModOption{
		AccessKeyID:     "test-ak",
		AccessKeySecret: "test-sk",
		Endpoint:        "endpoint",
	})
}

// 场景四: 在 OnChangeHandler 中初始化变量，模块中仅声明变量，不绑定 key
// 可读性：模块内部的可读性较好，main 会随着配置项增多而变差
// 复用性：配置项不再和 key 绑定，模块是可复用的
// 可测试性：和场景二一样
// 可维护性：和场景二一样
// 安全性：和场景二一样
func TestExample4(t *testing.T) {
	// package module
	var AccessKeyID config.AtomicString
	var AccessKeySecret config.AtomicString
	var Endpoint config.AtomicString

	// package main
	conf, err := config.NewConfigWithBaseFile("testfile/base.json")
	if err != nil {
		panic(err)
	}
	conf.AddOnChangeHandler(func(conf *config.Config) {
		AccessKeyID.Set(conf.GetString("OSS.AccessKeyID"))
		AccessKeySecret.Set(conf.GetString("OSS.AccessKeySecret"))
		Endpoint.Set(conf.GetString("OSS.Endpoint"))
	})
	if err := conf.Watch(); err != nil {
		panic(err)
	}
	defer conf.Stop()

	// package module
	fmt.Println(AccessKeyID.Get())
	fmt.Println(AccessKeySecret.Get())
	fmt.Println(Endpoint.Get())

	// package test
	AccessKeyID = *config.NewAtomicString("test-ak")
	AccessKeySecret = *config.NewAtomicString("test-sk")
	Endpoint = *config.NewAtomicString("endpoint")
}

// 场景五: 在 OnChangeHandler 中初始化变量，模块中仅声明变量，不绑定 key
// 可读性：模块内部的可读性较好，main 中每个模块都只有一次 Unmarshal 操作，不会随着配置项增多而增加复杂性，可读性较好
// 复用性：配置项不和 key 绑定，模块是可复用的
// 可测试性：和场景三一样
// 可维护性：和场景三一样
// 安全性：和场景三一样
func TestExample5(t *testing.T) {
	// package module
	type ModOption struct {
		AccessKeyID     string
		AccessKeySecret string
		Endpoint        string
	}
	var option atomic.Value

	// package main
	conf, err := config.NewConfigWithBaseFile("testfile/base.json")
	if err != nil {
		panic(err)
	}
	conf.AddOnChangeHandler(func(conf *config.Config) {
		var opt ModOption
		if err := conf.Sub("OSS").Unmarshal(&opt); err != nil {
			return
		}
		option.Store(opt)
	})
	if err := conf.Watch(); err != nil {
		panic(err)
	}
	defer conf.Stop()

	// package module
	fmt.Println(option.Load().(ModOption).AccessKeyID)
	fmt.Println(option.Load().(ModOption).AccessKeySecret)
	fmt.Println(option.Load().(ModOption).Endpoint)

	// package test
	option.Store(ModOption{
		AccessKeyID:     "test-ak",
		AccessKeySecret: "test-sk",
		Endpoint:        "endpoint",
	})
}

// 场景六: 使用绑定变量作为参数传递给构造函数
// 可读性：模块内部的可读性较好，对象在 main 中构造，main 会随着对象增多而变复杂
// 复用性：可以做到对象级别的复用，复用性最好
// 可测试性：直接通过构造函数构造测试对象，测试很方便
// 可维护性：每个对象维护自己的动态参数，维护比较方便
// 安全性：无法保证关联配置的原子性
func TestExample6(t *testing.T) {
	// package main
	conf, err := config.NewConfigWithBaseFile("testfile/base.json")
	if err != nil {
		panic(err)
	}
	if err := conf.Watch(); err != nil {
		panic(err)
	}
	defer conf.Stop()

	myType := NewMyType6(
		conf.String("OSS.AccessKeyID"),
		conf.String("OSS.AccessKeySecret"),
		conf.String("OSS.Endpoint"),
	)

	// package module
	myType.DoSomething()

	// package test
	testType := NewMyType6(config.NewAtomicString("test-ak"), config.NewAtomicString("test-sk"), config.NewAtomicString("endpoint"))
	testType.DoSomething()
}

// package module
type MyType6 struct {
	AccessKeyID     *config.AtomicString
	AccessKeySecret *config.AtomicString
	Endpoint        *config.AtomicString
}

func NewMyType6(accessKeyID *config.AtomicString, accessKeySecret *config.AtomicString, endpoint *config.AtomicString) *MyType6 {
	return &MyType6{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		Endpoint:        endpoint,
	}
}

func (t MyType6) DoSomething() {
	fmt.Println(t.AccessKeyID.Get())
	fmt.Println(t.AccessKeySecret.Get())
	fmt.Println(t.Endpoint.Get())
}

// 场景七: 使用绑定的结构体作为参数传递给构造函数
// 可读性：模块内部的可读性较好，对象在 main 中构造，main 会随着对象增多而变复杂
// 复用性：可以做到对象级别的复用，复用性好
// 可测试性：直接通过构造函数构造测试对象，测试很方便
// 可维护性：每个对象维护自己的动态参数，维护比较方便
// 安全性：可以保证关联配置项的原子性
func TestExample7(t *testing.T) {
	// package main
	conf, err := config.NewConfigWithBaseFile("testfile/base.json")
	if err != nil {
		panic(err)
	}
	if err := conf.Watch(); err != nil {
		panic(err)
	}
	defer conf.Stop()

	myType := NewMyType7(conf.Bind("OSS", MyType2Option{}))
	myType.DoSomething()

	// package test
	var option atomic.Value
	option.Store(MyType2Option{
		AccessKeyID:     "test-ak",
		AccessKeySecret: "test-sk",
		Endpoint:        "endpoint",
	})
	testType := NewMyType7(&option)
	testType.DoSomething()
}

// package module
type MyType7 struct {
	option *atomic.Value
}

type MyType2Option struct {
	AccessKeyID     string
	AccessKeySecret string
	Endpoint        string
}

func NewMyType7(option *atomic.Value) *MyType7 {
	return &MyType7{
		option: option,
	}
}

func (t MyType7) DoSomething() {
	fmt.Println(t.option.Load().(MyType2Option).AccessKeyID)
	fmt.Println(t.option.Load().(MyType2Option).AccessKeySecret)
	fmt.Println(t.option.Load().(MyType2Option).Endpoint)
}

// 场景八: 类似缓存的使用方式，用配置的 key 来初始化对象
// 可读性：模块内部的可读性较好，对象在 main 中构造，main 会随着对象增多而变复杂
// 复用性：不依赖固定的配置项，对象是可复用的
// 可测试性：需要 mock 全局的配置对象，mock 代价较高，和场景一类似
// 可维护性：和场景一类似
// 安全性：和场景一类似
func TestExample8(t *testing.T) {
	// package main
	if err := config.Init("testfile/base.json"); err != nil {
		panic(err)
	}
	if err := config.Watch(); err != nil {
		panic(err)
	}
	defer config.Stop()

	myType := NewMyType8("OSS.AccessKeyID", "OSS.AccessKeySecret", "OSS.Endpoint")

	// package module
	myType.DoSomething()

	// package test
	testType := NewMyType8("OSS.AccessKeyID", "OSS.AccessKeySecret", "OSS.Endpoint")
	_ = config.UnsafeSet("OSS.AccessKeyID", "test-ak")
	_ = config.UnsafeSet("OSS.AccessKeyID", "test-sk")
	_ = config.UnsafeSet("OSS.AccessKeyID", "endpoint")
	testType.DoSomething()
}

// package module
type MyType8 struct {
	AccessKeyIDConfigKey     string
	AccessKeySecretConfigKey string
	EndpointConfigKey        string
}

func NewMyType8(AccessKeyIDConfigKey string, AccessKeySecretConfigKey string, EndpointConfigKey string) *MyType8 {
	return &MyType8{
		AccessKeyIDConfigKey:     AccessKeyIDConfigKey,
		AccessKeySecretConfigKey: AccessKeySecretConfigKey,
		EndpointConfigKey:        EndpointConfigKey,
	}
}

func (t MyType8) DoSomething() {
	fmt.Println(config.GetString(t.AccessKeyIDConfigKey))
	fmt.Println(config.GetString(t.AccessKeySecretConfigKey))
	fmt.Println(config.GetString(t.EndpointConfigKey))
}
