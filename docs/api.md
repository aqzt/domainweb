# 域名估价系统 API 文档

本文档详细说明了域名估价系统提供的API接口，包括请求参数、响应格式和示例。

## API 概述

域名估价系统提供了RESTful风格的API，支持JSON格式的请求和响应。所有API都以`/api`为前缀。

## 认证

当前版本不需要认证。未来版本可能会添加API密钥认证机制。

## API 端点

### 1. 域名估价

#### 请求

- **URL**: `/api/estimate`
- **方法**: POST
- **Content-Type**: `application/json`
- **请求体**:

```json
{
  "domain": "example.com"
}
```

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| domain | string | 是 | 要估价的域名，如example.com |

#### 响应

- **Content-Type**: `application/json`
- **状态码**: 200 OK
- **响应体**:

```json
{
  "domain": "example.com",
  "grade": 3.5,
  "price": 11917,
  "baseAttributes": [
    {
      "name": "com后缀",
      "value": "com",
      "description": "com后缀",
      "priceFactor": 9.55,
      "gradeFactor": 0.5
    },
    {
      "name": "7位长度",
      "value": "7",
      "description": "7位长度",
      "priceFactor": 1.2,
      "gradeFactor": 0.1
    },
    {
      "name": "纯字母结构",
      "value": "纯字母",
      "description": "纯字母结构",
      "priceFactor": 1.26,
      "gradeFactor": 0.31
    }
  ],
  "otherAttributes": [
    {
      "name": "Alexa排名良好",
      "value": "85462",
      "description": "Alexa 排名 85462",
      "priceFactor": 1.8,
      "gradeFactor": 0.5
    },
    {
      "name": "搜索量较高",
      "value": "3480",
      "description": "搜索量 3480",
      "priceFactor": 1.8,
      "gradeFactor": 0.6
    }
  ],
  "estimationDate": "2023-05-10T14:30:45Z"
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| domain | string | 估价的域名 |
| grade | number | 品相等级，范围通常在0-10之间 |
| price | number | 保守估价，单位为人民币元 |
| baseAttributes | array | 基础属性列表，包含影响估价的基础因素 |
| otherAttributes | array | 其他属性列表，包含影响估价的动态因素 |
| estimationDate | string | 估价时间，ISO 8601格式 |

#### 错误响应

- **状态码**: 400 Bad Request
- **响应体**:

```json
{
  "error": "无效的请求参数"
}
```

- **状态码**: 500 Internal Server Error
- **响应体**:

```json
{
  "error": "估价失败: 无法解析域名"
}
```

### 2. 查询历史记录

#### 请求

- **URL**: `/api/history`
- **方法**: GET
- **参数**:

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| domain | string | 否 | 按域名筛选，支持部分匹配 |
| limit | integer | 否 | 限制返回记录数量，默认50 |

#### 响应

- **Content-Type**: `application/json`
- **状态码**: 200 OK
- **响应体**:

```json
[
  {
    "id": 1,
    "domain": "example.com",
    "grade": 3.5,
    "price": 11917,
    "estimationDate": "2023-05-10T14:30:45Z"
  },
  {
    "id": 2,
    "domain": "domain.com",
    "grade": 4.2,
    "price": 15680,
    "estimationDate": "2023-05-10T14:35:12Z"
  }
]
```

#### 错误响应

- **状态码**: 500 Internal Server Error
- **响应体**:

```json
{
  "error": "获取历史记录失败: 数据库连接错误"
}
```

## 状态码

| 状态码 | 描述 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 数据模型

### EstimationResult

| 字段 | 类型 | 描述 |
|------|------|------|
| domain | string | 域名 |
| grade | number | 品相等级 |
| price | number | 保守估价 |
| baseAttributes | AttributeDetail[] | 基础属性详情 |
| otherAttributes | AttributeDetail[] | 其他属性详情 |
| estimationDate | string | 估价时间 |

### AttributeDetail

| 字段 | 类型 | 描述 |
|------|------|------|
| name | string | 属性名称 |
| value | string | 属性值 |
| description | string | 属性描述 |
| priceFactor | number | 估价倍数 |
| gradeFactor | number | 等级增量 |

### HistoryRecord

| 字段 | 类型 | 描述 |
|------|------|------|
| id | integer | 记录ID |
| domain | string | 域名 |
| grade | number | 品相等级 |
| price | number | 估价结果 |
| estimationDate | string | 查询时间 |
