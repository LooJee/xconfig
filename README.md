# what

xconfig 是基于 agollo 开发的一个配置库，可以将 apollo 配置中心的配置解析到对应的结构中。

# 使用

## 结构定义

```go
type Db struct {
    Postgres struct {
        Host  string
        Port  uint
        Debug bool
    }
}

type App struct {
    Grpc struct {
        Name string
        Port string
    }
}

type Config struct {
    Db  Db  `xconfig:"appId:db-config"`
    App App `xconfig:"appId:commodity-srv"`
}
```

上面这个结构是一个常见的可以被 xconfig 解析的格式，在 tag 中指定 xconfig 可以传递自定义的参数，
目前支持以下几个参数。

 - name 用于指定字段名在 apollo 配置中心的映射
 - default 用于指定字段的默认值，当配置中心中未配置相应字段时，会使用该值，未指定时会采用类型 0 值。
 - appId 用于指定结构中的配置是对应哪个 appId 的，通常只需要最外层的结构需要标记

当指定为 `xconfig:"-"` 时，表示该字段不需要解析

### 字段名映射

结构体中的字段和 apollo 配置中心的字段映射基本规则是这样的 : 
在 tag 中指定 name ，程序将会根据该值来解析结构中的字段;若未指定该值，则会取字段名，将首字母转为小写后作为 apollo 中的字段名处理。

根据这个规则有两种编排配置字段的方式 :

1. 定义一个大结构，所有的配置在结构中平铺，然后每个字段定义一个包含 name 的 xconfig tag，name 是 apollo 配置中心的字段名全称。
```go
type Config struct {
    Host        string `xconfig:"name:db.postgres.host"`
    User        string `xconfig:"name:db.postgres.user"`
    Password    string `xconfig:"name:db.postgres.password"`
    TimeZone    string `xconfig:"name:db.postgres.timeZone"`
}
```
2. 定义一个嵌套结构，层级关系和 apollo 配置中心的字段名点分层级关系相同，这时候可以不指定 name。
```go
type Db struct {
    Postgres struct {
        Host  string
        Port  uint
        Debug bool
    }
}

type App struct {
    Grpc struct {
        Name string
        Port string
    }
}

type Config struct {
    Db  Db  `xconfig:"appId:db-config"`
    App App `xconfig:"appId:commodity-srv"`
}
```
> Config 结构里的是配置中心字段第一级的名称

### 默认值

在 tag 中可以通过 default 指定一个字段的默认值，如果未指定，则会使用零值。
```go
type Config struct {
    Host        string `xconfig:"name:db.postgres.host"` //默认值为""
    User        string `xconfig:"name:db.postgres.user;default:hehe"` //默认值为hehe
}
```

## 获取配置

配置在本库不是一个权益的值，每个配置都需要一个client来初始化 : 
```go
cfg := Config{}
if err := NewConfig(&cfg, 
	ApolloIp("apollo.api.thingyouwe.com"), 
	LocalConfig("local.yaml")); err != nil {
	    return err
	}
```

## Options

本库支持通过 Option 来自定义 apollo 集群的配置信息，具体内容查看 config.go 中的 ConfigOption。

## 本地配置

指定 LocalConfig option可以加载本地的配置。本地配置优先级最高，会覆盖从 apollo 配置中心获取的配置。
本地文件仅支持 yaml 文件，内容格式如下所示 :

```yaml
db-config: // appId
  db.postgres.host: postgres2 //格式1，一个key的所有分段在一起
  db:
    postgres:
      host: postgres2
      port: 1234
      debug: false            //格式2，一个可以的分段根据yaml规则分层
commodity-srv:
  app:
    grpc:
      port: :1111
      name: commodity-srv1
```

## 热更新

目前只有内置的一个热更新，如果有需求再开发。
