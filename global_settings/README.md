# 全局设置 (global_settings)

`global_settings`包提供了应用程序级别的全局设置和上下文管理功能。

## 功能组件

### 环境管理

环境管理模块提供了应用程序运行环境的设置和获取方法，支持开发、测试、预发布和生产环境。

```go
// 设置环境
global_settings.SetEnv(global_settings.ENV_PROD)

// 获取当前环境
env := global_settings.GetEnv()

// 根据环境执行不同操作
switch env {
case global_settings.ENV_DEV:
    // 开发环境处理
case global_settings.ENV_TEST:
    // 测试环境处理
case global_settings.ENV_STAGING:
    // 预发布环境处理
case global_settings.ENV_PROD, global_settings.ENV_LIVE:
    // 生产环境处理
}
```

预定义的环境常量:
- `ENV_DEV`: 开发环境
- `ENV_TEST`: 测试环境
- `ENV_STAGING`: 预发布环境
- `ENV_PROD`, `ENV_LIVE`: 生产环境

### 服务上下文

服务上下文模块提供了应用程序级别的全局上下文管理，用于存储跨请求的共享数据或服务级别的配置。

```go
// 重置服务上下文
ctx := context.Background()
ctx = context.WithValue(ctx, "app_version", "1.0.0")
global_settings.ResetServiceContext(ctx)

// 获取服务上下文
serviceCtx := global_settings.GetServiceContext()
appVersion := serviceCtx.Value("app_version").(string)
```

## 最佳实践

### 初始化全局设置

在应用程序启动时初始化全局设置：

```go
package main

import (
    "context"
    "github.com/gw-gong/gwkit-go/global_settings"
    "github.com/gw-gong/gwkit-go/log"
    "os"
)

func main() {
    // 设置环境
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = global_settings.ENV_DEV
    }
    global_settings.SetEnv(env)
    
    // 初始化日志
    var logConfig log.Config
    if global_settings.GetEnv() == global_settings.ENV_DEV {
        logConfig = log.NewDevelopmentConfig()
    } else {
        logConfig = log.NewProductionConfig()
    }
    log.InitGlobalLogger(logConfig)
    
    // 初始化服务上下文
    ctx := context.Background()
    ctx = context.WithValue(ctx, "app_start_time", time.Now())
    global_settings.ResetServiceContext(ctx)
    
    // 启动应用
    // ...
}
```

### 在中间件中使用

可以在HTTP中间件中使用全局设置：

```go
func EnvCheckMiddleware(c *gin.Context) {
    if global_settings.GetEnv() == global_settings.ENV_PROD {
        // 生产环境特定逻辑
    }
    c.Next()
}
```

## 注意事项

- 服务上下文应谨慎使用，避免存储请求级别的数据
- 尽量在应用程序启动时设置环境，避免运行时更改
- 全局设置应限制在配置类信息，避免用于存储业务数据 