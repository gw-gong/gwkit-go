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
    syncFn, err := log.InitGlobalLogger(config)
    if err != nil {
        panic(err)
    }
    defer syncFn()

    // 使用日志器
    log.Info("应用启动成功")
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

syncFn, err := log.InitGlobalLogger(config)
if err != nil {
    panic(err)
}
defer syncFn()
```

## 基本使用

### 日志级别

```go
// 不同级别的日志
log.Debug("调试信息")
log.Info("普通信息")
log.Warn("警告信息")
log.Error("错误信息")
```

### 结构化日志

```go
// 添加字段
log.Info("用户登录",
    log.String("user_id", "123"),
    log.String("ip", "192.168.1.1"),
    log.Int("attempt", 1),
)

// 错误日志
if err != nil {
    log.Error("操作失败",
        log.String("operation", "create_user"),
        log.Error(err),
    )
}
```

## 与Context结合使用

### 基本用法

```go
// 将全局日志器存入context
ctx := log.SetGlobalLoggerToCtx(context.Background())

// 添加字段
ctx = log.WithFields(ctx, 
    log.String("request_id", "req-123"),
    log.String("user_id", "user-456"),
)

// 使用带上下文的日志方法
log.Infoc(ctx, "处理请求")
log.Errorc(ctx, "请求处理失败", log.Error(err))
```

## 配置说明

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

## 注意事项

1. 默认输出到控制台，使用控制台友好的格式
2. 默认日志级别为 DEBUG
3. 默认添加调用者信息
4. 默认不启用堆栈跟踪
5. 必须至少启用一种输出方式（文件或控制台）
6. 使用 `defer syncFn()` 确保日志被正确刷新
