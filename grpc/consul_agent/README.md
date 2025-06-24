# Consul Agent

这个目录提供了基于Consul的服务注册与发现功能，用于gRPC服务的服务治理。

## 功能特性

- **服务注册**: 支持普通和TLS两种方式注册gRPC服务到Consul
- **服务发现**: 从Consul获取健康的服务实例端点
- **连接管理**: 提供基于Consul的gRPC客户端连接创建
- **健康检查**: 自动进行gRPC健康检查

## 文件说明

### registry.go
提供Consul服务注册与发现的核心功能：

- `ConsulRegistry` 接口：定义服务注册与发现的基本操作
- `NewConsulRegistry()`: 创建Consul注册表实例
- `Service()`: 获取健康的服务实例端点
- `Register()`: 注册服务实例（普通模式）
- `RegisterTLS()`: 注册服务实例（TLS模式）
- `Deregister()`: 注销服务实例

### grpc_conn.go
提供基于Consul的gRPC连接管理：

- `NewHealthyGrpcConn()`: 创建到健康服务的gRPC连接
- `CheckServiceExists()`: 检查服务是否存在
- `formatGrpcConnTarget()`: 格式化gRPC连接目标地址

## 使用示例

### 服务注册
```go
registry, err := NewConsulRegistry("my-service")
if err != nil {
    log.Fatal(err)
}

// 注册服务
err = registry.Register("my-service-1", 8080, []string{"v1", "prod"})
if err != nil {
    log.Fatal(err)
}
```

### 服务发现
```go
registry, err := NewConsulRegistry("my-service")
if err != nil {
    log.Fatal(err)
}

// 获取服务端点
endpoints, err := registry.Service([]string{"v1", "prod"})
if err != nil {
    log.Fatal(err)
}
```

### 创建gRPC连接
```go
conn, err := NewHealthyGrpcConn("127.0.0.1:8500", "my-service", "v1", "")
if err != nil {
    log.Fatal(err)
}
defer conn.Close()
```

## 配置说明

- **默认Consul地址**: `127.0.0.1:8500`
- **健康检查间隔**: 10秒
- **健康检查超时**: 3秒
- **支持TLS**: 可配置TLS健康检查

## 依赖

- `github.com/hashicorp/consul/api`: Consul API客户端
- `github.com/mbobakov/grpc-consul-resolver`: gRPC Consul解析器
- `google.golang.org/grpc`: gRPC框架 