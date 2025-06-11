参考设计：https://help.aliyun.com/document_detail/53414.html

# 统一响应结构

## 响应格式

所有API请求的响应都遵循统一的JSON格式：

```json
{
    "code": 0,        // 业务状态码，0表示成功，非0表示失败
    "msg": "success", // 状态描述信息
    "request_id": "550e8400-e29b-41d4-a716-446655440000", // 请求唯一标识
    "data": {},       // 业务数据，成功时返回，可以是任意JSON结构
    "err_details": {} // 错误详情，仅在错误时返回，可选字段
}
```

## 成功响应

成功响应的`code`固定为0，`data`字段包含返回的业务数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "data": {
        "user_id": 12345,
        "username": "example"
    }
}
```

## 错误响应

错误响应的`code`为非0值，表示特定的错误类型：

```json
{
    "code": 10001,
    "msg": "invalid json format",
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "data": null,
    "err_details": {
        "field": "username",
        "reason": "required field missing"
    }
}
```

`err_details`字段为可选，用于提供更详细的错误信息。

## 使用方法

### 返回成功响应

```go
// 返回带数据的成功响应
func GetUserInfo(c *gin.Context) {
    user := getUserFromDB()
    response.ResponseSuccess(c, user)
}

// 返回不带数据的成功响应
func DeleteUser(c *gin.Context) {
    deleteUserFromDB()
    response.ResponseSuccess(c, nil)
}
```

### 返回错误响应

```go
// 返回简单错误
func GetUser(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        response.ResponseError(c, gwkit_res.ErrParam)
        return
    }
    // ...
}

// 返回带详情的错误
func CreateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        details := map[string]interface{}{
            "errors": []string{"用户名不能为空", "密码长度必须大于6位"},
            "fields": []string{"username", "password"}
        }
        response.ResponseErrorWithDetails(c, gwkit_res.ErrInvalidJSON, details)
        return
    }
    // ...
}
```

## 错误码设计

错误码按HTTP状态码分类：
- 客户端错误（4xx）：10000-19999
- 服务端错误（5xx）：20000-29999

具体错误码定义见 `http/response/error_code.go`
