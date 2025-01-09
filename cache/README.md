# 说明

Golang本地缓存器， 变量缓存本地进程内存中，超时释放。

# 使用

```shell
go get gitee.com/ivfzhou/study_golang/cache@latest
```

```golang
import cache "gitee.com/ivfzhou/study_golang/cache"

// 创建缓存
c := cache.New()

// 设置缓存
c.Set("key", value, time.Second*10)

// 获取缓存
val, ok := c.Get("key")

// 删除缓存
c.Del("key")

// 设置缓存最大容量
c.SetMaxMemory("256mb")

// 查询缓存是否存在
c.Exists("key")

// 占用内存大小
c.Size()
```
