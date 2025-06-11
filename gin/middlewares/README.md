# Gin中间件

本包提供了一组用于Gin框架的中间件，用于处理常见的Web应用功能。

## 中间件列表

### PanicRecover
防止服务因panic而崩溃，捕获所有panic并恢复服务正常运行。

### InjectLoggerToCtx
将全局日志记录器注入到请求上下文中，使日志在整个请求生命周期内可用。

### GenerateRID
为每个请求生成唯一的请求ID（Request ID），并将其添加到上下文和日志中，便于请求追踪。

### LogHttpReqInfo
记录HTTP请求的详细信息，包括请求方法、路径、查询参数、客户端信息、请求头和响应状态等。

## 推荐的中间件注册顺序

中间件的注册顺序很重要，建议按照以下顺序注册：

1. `PanicRecover` - 首先注册，确保捕获所有其他中间件和处理程序中的panic
2. `InjectLoggerToCtx` - 注入日志记录器，为后续中间件提供日志功能
3. `GenerateRID` - 生成请求ID，使日志可以与特定请求关联
4. `LogHttpReqInfo` - 记录请求信息，应在其他功能性中间件之前

## 使用方法

### 单独使用中间件

```go
import (
    "github.com/gw-gong/gwkit-go/gin/middlewares"
    "github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
    r := gin.New()
    
    // 注册单个中间件
    r.Use(middlewares.PanicRecover)
    r.Use(middlewares.InjectLoggerToCtx)
    r.Use(middlewares.GenerateRID)
    r.Use(middlewares.LogHttpReqInfo(true)) // true表示记录请求体内容
    
    return r
}
```

### 使用预设的基础中间件组

```go
import (
    "github.com/gw-gong/gwkit-go/gin/middlewares"
    "github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
    r := gin.New()
    
    // 一次性注册所有基础中间件
    middlewares.BindBasicMiddlewares(r)
    
    return r
}
```

## 请求ID (RID) 使用

请求ID可以在处理程序中获取：

```go
func Handler(c *gin.Context) {
    rid, exists := c.Get(middlewares.ContextKeyRID)
    if exists {
        // 使用请求ID
        requestID := rid.(string)
    }
}
```

## 日志记录定制

可以通过设置LogHttpReqInfo参数来控制请求体内容是否记录到日志中：

```go
// 记录请求体内容（可能包含敏感信息）
r.Use(middlewares.LogHttpReqInfo(true))

// 不记录请求体内容（用于生产环境或处理敏感数据）
r.Use(middlewares.LogHttpReqInfo(false))
```
