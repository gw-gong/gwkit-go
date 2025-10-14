# HTTP响应错误码

本包提供了统一的HTTP错误码定义和处理机制，包括错误码格式规范和预定义的通用错误码。

## 设计理念

[设计参考](https://learn.lianglianglee.com/%E4%B8%93%E6%A0%8F/Go%20%E8%AF%AD%E8%A8%80%E9%A1%B9%E7%9B%AE%E5%BC%80%E5%8F%91%E5%AE%9E%E6%88%98/18%20%E9%94%99%E8%AF%AF%E5%A4%84%E7%90%86%EF%BC%88%E4%B8%8A%EF%BC%89%EF%BC%9A%E5%A6%82%E4%BD%95%E8%AE%BE%E8%AE%A1%E4%B8%80%E5%A5%97%E7%A7%91%E5%AD%A6%E7%9A%84%E9%94%99%E8%AF%AF%E7%A0%81%EF%BC%9F.md)

HTTP状态码使用简化策略，主要使用以下三种HTTP状态码：

- `200` - 表示请求成功执行
- `400` - 表示客户端出问题
- `500` - 表示服务端出问题

## 错误码格式

错误码采用9位数字格式，便于分类和管理：

```
[3位:服务][3位:服务模块][3位:错误类型]
```

通用错误码前缀：
```
[100][000][具体错误]
```

## 预定义错误码

### 成功响应

| 错误码 | 消息 | HTTP状态码 | 使用场景 |
|-------|------|-----------|---------|
| 0 | success | 200 | 请求成功处理并返回 |

### 客户端错误 (4xx)

| 错误码 | 消息 | HTTP状态码 | 使用场景 |
|-------|------|-----------|---------|
| 100000000 | param error | 400 | 通用参数错误，当参数不符合要求时使用 |
| 100000001 | invalid json format | 400 | 请求体JSON格式错误 |
| 100000002 | invalid query parameter | 400 | 查询参数格式或值不符合要求 |
| 100000003 | missing required parameter | 400 | 缺少必需的参数 |
| 100000004 | unauthorized | 401 | 用户未认证或未登录 |
| 100000005 | token expired | 401 | 用户令牌已过期，需要重新登录 |
| 100000006 | invalid token | 401 | 用户令牌无效或被篡改 |
| 100000007 | forbidden | 403 | 用户无权访问该资源 |
| 100000008 | permission denied | 403 | 用户权限不足，无法执行操作 |
| 100000009 | resource not found | 404 | 请求的资源不存在 |
| 100000010 | method not allowed | 405 | 请求方法不被允许 |
| 100000011 | resource conflict | 409 | 资源冲突，如唯一键冲突 |
| 100000012 | too many requests | 429 | 请求频率超过限制，触发限流 |
| 100000013 | request entity too large | 413 | 请求体过大 |

### 服务端错误 (5xx)

| 错误码 | 消息 | HTTP状态码 | 使用场景 |
|-------|------|-----------|---------|
| 100000100 | internal server error | 500 | 服务器内部错误，未处理的异常 |
| 100000101 | database error | 500 | 数据库操作失败 |
| 100000102 | cache service error | 500 | 缓存服务异常 |
| 100000103 | third-party service error | 500 | 第三方服务调用失败 |
| 100000104 | bad gateway | 502 | 网关错误，上游服务返回无效响应 |
| 100000105 | service unavailable | 503 | 服务不可用，可能是过载或维护 |
| 100000106 | gateway timeout | 504 | 网关超时，上游服务响应超时 |
| 100000199 | unknown error | 500 | 未知错误，无法分类的服务端错误 |

## 使用方法

### 定义自定义错误码

```go
// 定义业务模块特定的错误码
var (
    // 用户模块错误码 - 前缀100001
    ErrUserNotFound = NewErrorCode(100001001, "user not found", http.StatusNotFound)
    ErrInvalidPassword = NewErrorCode(100001002, "invalid password", http.StatusBadRequest)
    
    // 订单模块错误码 - 前缀100002
    ErrOrderNotFound = NewErrorCode(100002001, "order not found", http.StatusNotFound)
    ErrOrderCancelled = NewErrorCode(100002002, "order already cancelled", http.StatusBadRequest)
)
```

### 在处理程序中使用

```go
func GetUser(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        response.ResponseError(c, ErrMissingRequiredParam)
        return
    }
    
    user, err := userService.GetUserByID(id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            response.ResponseError(c, ErrUserNotFound)
            return
        }
        response.ResponseError(c, ErrDatabase)
        return
    }
    
    response.ResponseSuccess(c, user)
}
```

## 最佳实践

1. 对于通用错误场景，优先使用预定义的通用错误码
2. 为每个业务模块定义特定前缀的错误码，方便问题定位
3. 在日志中记录详细错误信息，但对客户端返回友好提示
4. 考虑在错误响应中添加requestID，方便追踪问题
5. 客户端错误应返回足够信息帮助用户修正错误
