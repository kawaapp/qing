# 设计文档

依赖：
1. Echo - Web 框架
2. Meddler - 简单的半自动的ORM框架

主体参考 [Drone](https://github.com/drone/drone) 的架构，分为以下主要模块：
1. Router - 路由逻辑
2. Server - 业务逻辑，处理具体的HTTP请求
3. Store  - 存储相关逻辑，封装了数据库相关的操作
4. Remote - 处理第三方授权相关逻辑

## 中间件

Store/Remote/AttachUser 为三个自定义中间件，主要负责连接存储/第三方登录/用户鉴权相关的模块，
通过 Context 将各个模块耦合在一起。

用户登录授权采用的是 Token 方案，使用中间件在请求到达的时候对用户身份进行检查，并把合法用户写入
Context 中，可以通过 `Context.get(".user")` 或 `session.User(Context)` 取得当前用户。
Token 的签名沿用 Drone 的方案，安全性更高些，它会是对每个用户使用不同的秘钥，它存储在用户表中。

## API设计

API 采用 RESTful 风格API

## 分页方案

采用 Page + Size/Limit 的方案。

## 点赞

## 标签

## 通知

## 数据库

采用类似 Drone 的数据库接入方式，可以支持切换不同类似的数据库，目前底层支持 MySQL/SQLite/PostgreSQL

## 编译部署
