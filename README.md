# 卡哇轻论坛

该论坛系统源自于 [卡哇微社区](https://kawaapp.com) 的后端核心系统，剥离了卡哇业务业务相关的模块，仅留下最核心的社区模块。
论坛采用前后端分离架构，前端采用 React 打造，后端服务用 Go.

当前开源的部分是后端模块，前端模块还在准备中，择良辰吉日开放。


# 关于 Restful
关联字段，比如拉取一个接龙或者评论列表的时候，如果需要同时拉取对应的用户资料，可以在API的查询字段添加 include 查询，
比如： /api/discussions?includes=user
得到的返回结果，会包含一个 entities 字段包含了 user 字段

{"code":0,"data":[],"entities":{
    1: { name: "123", id: 1} // id => user
}}

# API 文档

**用户**

1. 注册 `POST /api/register`

```
{
	"username":"ntop",
	"password":"123"
}
```

2. 登录 `POST /api/login`
```
{
	"username":"ntop",
	"password":"123"
}
{
    "id": 1,
    "status": 0,
    "role": 0,
    "name": "ntop",
    "email": "",
    "title": "",
    "summary": "",
    "text": "",
    "image_id": 0,
    "password": ""
}
```
3. 登出 `POST /api/logout`
4. 授权(微信)登录 `POST /api/auth`
```
{
    code: 12345
}
{
    "id": 1,
    "status": 0,
    "role": 0,
    "name": "ntop",
    "email": "",
    "title": "",
    "summary": "",
    "text": "",
    "image_id": 0,
    "password": ""
}
```

**接龙**
1. 发起一个接龙 `POST /api/discussions`
```
{
	"content":"如何看待小米股价暴跌!!",
	"title":"小米股价腰斩喽...目前9块"
}
```
2. 删除一个接龙个 `DELETE /api/discussions/:id`
3. 获取接龙下的评论列表 `Get /api/discussions/:id/posts/`
```
[
    {
        "id": 1,
        "post": 1,
        "other": 0,
        "author": 0,
        "text": "快来评论一下下",
        "status": 0,
        "name": ""
    }
]
```
4. 发表评论 `POST /api/discussions/posts`
```
{
	"discussion_id":1,
	"author_id": 1,
	"content": "快来评论一下下"，
}
```
5. 删除评论 `DEL  /api/discussions/posts/:id`