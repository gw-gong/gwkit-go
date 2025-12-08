# gRPC Interceptors

这个目录提供了gRPC客户端和服务端的拦截器功能，用于增强gRPC通信的可靠性和可观测性。

## 目录结构

```
interceptors/
├── client/           # 客户端拦截器
│   └── unary/       # 一元调用拦截器
├── server/           # 服务端拦截器
│   └── unary/       # 一元调用拦截器
└── meta_data/       # 元数据定义
```

## 功能特性

- **请求追踪**: 自动传递RequestID和TraceID
- **异常恢复**: 服务端panic自动恢复
- **元数据管理**: 统一的元数据键定义
- **上下文传递**: 客户端到服务端的上下文信息传递

## 文件说明

### meta_data/keys.go
定义gRPC元数据中使用的键名：

- `MetaKeyRequestID`: 请求ID键名
- `MetaKeyTraceID`: 追踪ID键名

### client/unary/inject_meta_from_ctx.go
客户端拦截器，从上下文注入元数据：

- `InjectMetaFromCtx()`: 创建客户端拦截器
- 自动从上下文提取RequestID和TraceID
- 将元数据注入到gRPC请求中

### server/unary/parse_meta_to_ctx.go
服务端拦截器，解析元数据到上下文：

- `ParseMetaToCtx()`: 创建服务端拦截器
- 从gRPC元数据中提取RequestID和TraceID
- 将元数据设置到请求上下文中
- 自动添加日志字段

### server/unary/panic_recover.go
服务端异常恢复拦截器：

- `PanicRecoverInterceptor()`: 创建panic恢复拦截器
- 自动捕获并处理panic
- 返回标准的gRPC错误响应
- 记录异常信息到日志

## 使用示例

### 客户端使用
```go
import (
    "github.com/gw-gong/gwkit-go/grpc/interceptors/client/unary"
    "google.golang.org/grpc"
)

// 创建带有拦截器的gRPC客户端连接
conn, err := grpc.Dial(
    "localhost:8080",
    grpc.WithUnaryInterceptor(unary.InjectMetaFromCtx()),
)
```

### 服务端使用
```go
import (
    "github.com/gw-gong/gwkit-go/grpc/interceptors/server/unary"
    "google.golang.org/grpc"
)

// 创建带有拦截器的gRPC服务器
server := grpc.NewServer(
    grpc.UnaryInterceptor(
        grpc.ChainUnaryInterceptor(
            unary.ParseMetaToCtx(),
            unary.PanicRecoverInterceptor(),
        ),
    ),
)
```

### 拦截器链使用
```go
// 客户端拦截器链
clientInterceptors := grpc.WithUnaryInterceptor(
    unary.InjectMetaFromCtx(),
)

// 服务端拦截器链
serverInterceptors := grpc.ChainUnaryInterceptor(
    unary.ParseMetaToCtx(),
    unary.PanicRecoverInterceptor(),
)
```

## 工作流程

1. **客户端发起请求**:
   - 从上下文提取RequestID和TraceID
   - 注入到gRPC元数据中
   - 发送请求到服务端

2. **服务端处理请求**:
   - 从gRPC元数据中提取RequestID和TraceID
   - 设置到请求上下文中
   - 添加日志字段
   - 执行业务逻辑
   - 异常时自动恢复

3. **响应返回**:
   - 服务端返回处理结果
   - 客户端接收响应

## 配置说明

- **元数据键**: 统一使用`request_id`和`trace_id`
- **异常处理**: 自动将panic转换为gRPC内部错误
- **日志集成**: 自动添加RequestID和TraceID到日志字段

## 依赖

- `google.golang.org/grpc`: gRPC框架
- `github.com/gw-gong/gwkit-go/util/common`: 通用工具
- `github.com/gw-gong/gwkit-go/log`: 日志模块 