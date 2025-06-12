# gwkit-go

`gwkit-go`是一个Go语言工具包集合，提供常用功能组件，方便Go应用程序开发。

## 功能组件

### HTTP相关

- [HTTP响应](http/response/README.md) - 统一的HTTP错误码定义和处理

### Gin框架扩展

- [响应处理](gin/response/README.md) - Gin框架的统一响应格式和处理
- [中间件](gin/middlewares/README.md) - 包含请求ID生成、日志记录、异常恢复等中间件

### 日志

- [日志模块](log/README.md) - 基于zap的日志组件，支持日志轮转和上下文

### 全局设置

- [全局设置](global_settings/README.md) - 提供环境管理和服务级别上下文管理

### 工具函数

- `utils/str` - 字符串处理工具，如UUID生成
- `utils/http` - HTTP请求处理工具
- `utils/time` - 时间处理工具
- `utils/common` - 通用工具函数

## 快速开始

### 安装

```bash
go get github.com/gw-gong/gwkit-go
```

### 使用Gin中间件和响应

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gw-gong/gwkit-go/gin/middlewares"
	"github.com/gw-gong/gwkit-go/gin/response"
	"github.com/gw-gong/gwkit-go/http/response"
)

func main() {
	r := gin.New()
	
	// 注册中间件
	middlewares.BindBasicMiddlewares(r)
	
	// 路由处理
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response.ResponseError(c, gwkit_res.ErrParam)
			return
		}
		
		// 处理成功情况
		user := map[string]interface{}{
			"id": id,
			"name": "用户名",
		}
		response.ResponseSuccess(c, user)
	})
	
	r.Run(":8080")
}
```

### 使用日志组件

```go
package main

import (
	"context"
	"github.com/gw-gong/gwkit-go/log"
)

func main() {
	// 初始化日志配置
	config := log.NewDefaultLoggerConfig()
	log.InitGlobalLogger(config)
	
	// 使用全局日志
	log.Info("应用启动")
	
	// 带上下文的日志
	ctx := context.Background()
	ctx = log.SetLoggerToCtx(ctx, log.GlobalLogger().With("requestID", "123456"))
	logger := log.GetLoggerFromCtx(ctx)
	logger.Info("处理请求")
}
```

### 使用全局设置

```go
package main

import (
	"context"
	"github.com/gw-gong/gwkit-go/global_settings"
	"github.com/gw-gong/gwkit-go/log"
	"os"
)

func main() {
	// 设置运行环境
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = global_settings.ENV_DEV
	}
	global_settings.SetEnv(env)
	
	// 根据环境初始化日志
	var logConfig log.Config
	if global_settings.GetEnv() == global_settings.ENV_DEV {
		logConfig = log.NewDevelopmentConfig()
	} else {
		logConfig = log.NewProductionConfig()
	}
	log.InitGlobalLogger(logConfig)
	
	// 初始化服务上下文
	ctx := context.Background()
	global_settings.ResetServiceContext(ctx)
	
	// ...
}
```

## 组件列表

| 组件 | 描述 | 文档 |
|------|------|------|
| `http/response` | HTTP响应处理 | [文档](http/response/README.md) |
| `gin/response` | Gin响应格式化 | [文档](gin/response/README.md) |
| `gin/middlewares` | Gin中间件集合 | [文档](gin/middlewares/README.md) |
| `log` | 日志处理组件 | [文档](log/README.md) |
| `global_settings` | 全局设置和环境管理 | [文档](global_settings/README.md) |
| `utils/str` | 字符串工具 | - |
| `utils/http` | HTTP工具 | - |
| `utils/time` | 时间工具 | - |
| `utils/common` | 通用工具 | - |

## 依赖

- [gin-gonic/gin](https://github.com/gin-gonic/gin) - HTTP Web框架
- [uber-go/zap](https://github.com/uber-go/zap) - 高性能日志库
- [google/uuid](https://github.com/google/uuid) - UUID生成

## 贡献

欢迎提交Issue和Pull Request