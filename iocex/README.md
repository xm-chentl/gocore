# Go依赖注入
## 安装
go get -v github.com/xm-chentl/go-ioc

## 示例
以下用服务多数据库方式做演示
```go
type IDb interface{
    Db() IDbRepository 
}

type IDbRepository interface{
    Creata(entry interface{}) error
}

type mysqlDb struct{}

func (mysqlDb) Db() IDbRepository{
    return nil
}

type mongoDb struct{}

func (mongoDb) Db() IDbRepository {
    return nil
}

```
### 注入
```go
// 结构体注入
// 默认注入方式
type UserService struct{
    Db IDb `inject:""` 
}

goioc.Set(new(IDb), &mysqlDb{})

// 指定名注入
type UserService struct{
    MongoDb IDb `inject:"mongo"`
}

goioc.SetMap(new(IDb), map[string]interface{}{
    "mongo": &mongoDb{},
})

// 多库注入
type UserService struct{
    MySqlDb IDb `inject:""` // 默认时，不需要指定名字, 只需要在注入时使用goioc.TagDefault
    MongoDb IDb `inject:"mongo"`
}

goioc.SetMap(new(IDb), map[string]interface{}{
    goioc.TagDefault: &mysqlDb{}, // 默认数据库
    "mongo": &mongoDb{},
})

// 以上示例的结构注入的方式，在main.go做统一注入如下：
userService := &UserService{}
// 注入
if err := goioc.Inject(userService); err != nil {
    panic("注入失败")
}
// userService.MongoDb....之后就可以正常使用了



```
### 注册
```go
// 单实例注入
goioc.Set(new(IDb), &mysqlDb{})
// 单实例，多实现注入
goioc.SetMap(new(IDb), map[string]interface{}{
    goioc.TagDefault: &mysqlDb{}, // 多实现的注入需要指定默认实现
    "mongo": &mongoDb{},
})
```
### 获取
```go
// 单实现获取
goioc.Get(new(IDb))
// 指定实现的获取
goioc.GetTag(new(IDb), "mongo").(IDb)

```

### 嵌套注入
```go
type InstDemo struct{
    InstA
}

type InstA struct{
    Db IDb `inject:""` 
}

// 支持针对嵌套
type InstDemo2 struct{
    *InstA
}

```