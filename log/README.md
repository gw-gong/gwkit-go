# 日志组件 (Logger)

这个日志组件基于 [zap](https://github.com/uber-go/zap) 和 [lumberjack](https://github.com/natefinch/lumberjack) 实现，提供了高性能的结构化日志记录和日志文件切分功能。

## 快速开始

### 初始化全局日志器

最简单的方式是使用默认配置初始化全局日志器：

```go
import "github.com/gw-gong/gwkit-go/log"

func main() {
    // 使用默认配置初始化全局日志器
    config := log.NewDefaultLoggerConfig()
    syncFn, err := log.InitGlobalLogger(*config)
    if err != nil {
        panic(err)
    }
    defer syncFn()

    // 使用全局日志器
    logger := log.GlobalLogger()
    logger.Info("应用启动成功")
}
```

### 自定义配置

你可以根据需要自定义日志配置：

```go
config := &log.LoggerConfig{
    Level: log.LoggerLevelInfo,
    OutputToFile: log.OutputToFileConfig{
        Enable:     true,
        FilePath:   "logs/app.log",
        MaxSize:    100,
        MaxBackups: 5,
        MaxAge:     7,
        Compress:   true,
    },
    OutputToConsole: log.OutputToConsoleConfig{
        Enable:   true,
        Encoding: log.OutputEncodingConsole,
    },
    AddCaller:  true,
    StackTrace: log.StackTraceConfig{
        Enable:     true,
        TraceLevel: log.LoggerLevelError,
    },
}

syncFn, err := log.InitGlobalLogger(*config)
if err != nil {
    panic(err)
}
defer syncFn()
```

## 与Context结合使用

该日志组件支持将日志器存放到 Context 中，并在需要时取出使用。这种方式特别适合在请求处理过程中追踪特定请求的日志信息。

### 存入Context

```go
import (
    "context"
    "github.com/gw-gong/gwkit-go/log"
    "go.uber.org/zap"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
    // 创建包含请求ID的日志器
    requestID := getRequestID(r)
    logger := log.GlobalLogger().With(zap.String("request_id", requestID))
  
    // 将日志器存入context
    ctx := log.SetLoggerToCtx(r.Context(), logger)
  
    // 使用更新后的context继续处理请求
    processRequest(ctx, w, r)
}
```

### 从Context获取日志器

```go
func processRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    // 从context获取日志器
    logger := log.GetLoggerFromCtx(ctx)
  
    // 使用日志器记录信息
    logger.Info("处理请求",
        zap.String("path", r.URL.Path),
        zap.String("method", r.Method),
    )
  
    // 在子函数中继续使用同一个context
    subProcess(ctx)
}

func subProcess(ctx context.Context) {
    logger := log.GetLoggerFromCtx(ctx)
  
    // 添加更多字段
    logger = logger.With(zap.String("sub_process", "data_processing"))
  
    logger.Debug("子流程执行中")
    // 处理逻辑...
}
```

如果context中没有日志器，`GetLoggerFromCtx`将返回全局日志器，确保总是能获取到有效的日志器。

## 配置详解

### LoggerConfig

| 字段            | 类型                  | 描述               | 默认值                               |
| --------------- | --------------------- | ------------------ | ------------------------------------ |
| Level           | string                | 日志级别           | "debug"                              |
| OutputToFile    | OutputToFileConfig    | 文件输出配置       | {Enable: false, ...}                 |
| OutputToConsole | OutputToConsoleConfig | 控制台输出配置     | {Enable: true, Encoding: "console"}  |
| AddCaller       | bool                  | 是否添加调用者信息 | true                                 |
| StackTrace      | StackTraceConfig      | 堆栈跟踪配置       | {Enable: false, TraceLevel: "error"} |

### OutputToFileConfig

| 字段       | 类型   | 描述                       | 默认值         |
| ---------- | ------ | -------------------------- | -------------- |
| Enable     | bool   | 是否启用文件输出           | false          |
| FilePath   | string | 日志文件路径               | "logs/app.log" |
| MaxSize    | int    | 单个日志文件最大大小（MB） | 500            |
| MaxBackups | int    | 保留的旧文件最大数量       | 10             |
| MaxAge     | int    | 保留旧文件的最大天数       | 30             |
| Compress   | bool   | 是否压缩                   | true           |

### OutputToConsoleConfig

| 字段     | 类型   | 描述                          | 默认值    |
| -------- | ------ | ----------------------------- | --------- |
| Enable   | bool   | 是否启用控制台输出            | true      |
| Encoding | string | 编码方式（"json"或"console"） | "console" |

### StackTraceConfig

| 字段       | 类型   | 描述                   | 默认值  |
| ---------- | ------ | ---------------------- | ------- |
| Enable     | bool   | 是否启用堆栈跟踪       | false   |
| TraceLevel | string | 触发堆栈跟踪的日志级别 | "error" |

## 默认行为

当不提供配置或提供的配置不完整时，系统将采用以下默认行为：

1. 默认输出到控制台，不输出到文件
2. 默认日志级别为 DEBUG
3. 默认添加调用者信息，便于追踪日志来源
4. 默认不启用堆栈跟踪
5. 默认使用控制台友好的编码方式（非JSON）

如果同时启用了文件输出和控制台输出，日志会被写入到指定的文件中，同时也会输出到控制台。

如果配置中未启用任何输出（文件输出和控制台输出都禁用），初始化时将返回错误。

## 后续优化计划

计划对 zap 进行完整封装，通过定义清晰的接口来隐藏底层实现细节。这样可以让使用者专注于日志功能的使用，而不需要了解 zap 的具体实现。主要目标包括：

1. 提供简单直观的 API 接口
2. 简化日志器的初始化和配置过程
3. 保持与现有功能的兼容性
4. 提供详细的接口文档和使用示例
