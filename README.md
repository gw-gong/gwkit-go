# gwkit-go

`gwkit-go`是一个Go语言工具包集合，提供常用功能组件，方便Go应用程序开发。

## 功能组件

### HTTP相关

- [HTTP响应](http/response/README.md) - 统一的HTTP错误码定义和处理

### Gin框架扩展

- [响应处理](gin/response/README.md) - Gin框架的统一响应格式和处理
- [中间件](gin/middlewares/README.md) - 包含请求ID生成、日志记录、异常恢复等中间件

### gRPC相关

- [gRPC拦截器](grpc/interceptors/README.md) - gRPC客户端和服务端拦截器，支持请求追踪、异常恢复、元数据传递
- [Consul服务治理](grpc/consul_agent/README.md) - 基于Consul的服务注册与发现，支持健康检查和连接管理

### 配置管理

- [热配置](hot_cfg/README.md) - 支持本地文件和Consul的热配置更新，无需重启应用

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

更多详细的使用示例请查看各模块的文档。

## 组件列表

| 组件 | 描述 | 文档 |
|------|------|------|
| `http/response` | HTTP响应处理 | [文档](http/response/README.md) |
| `gin/response` | Gin响应格式化 | [文档](gin/response/README.md) |
| `gin/middlewares` | Gin中间件集合 | [文档](gin/middlewares/README.md) |
| `grpc/interceptors` | gRPC拦截器 | [文档](grpc/interceptors/README.md) |
| `grpc/consul_agent` | Consul服务治理 | [文档](grpc/consul_agent/README.md) |
| `hot_cfg` | 热配置管理 | [文档](hot_cfg/README.md) |
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
- [google.golang.org/grpc](https://github.com/grpc/grpc-go) - gRPC框架
- [hashicorp/consul/api](https://github.com/hashicorp/consul/api) - Consul API客户端
- [mbobakov/grpc-consul-resolver](https://github.com/mbobakov/grpc-consul-resolver) - gRPC Consul解析器
- [spf13/viper](https://github.com/spf13/viper) - 配置管理库
- [fsnotify/fsnotify](https://github.com/fsnotify/fsnotify) - 文件系统监控

## 贡献

欢迎提交Issue和Pull Request