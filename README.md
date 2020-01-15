# 群接龙API

# 关于 Restful
关联字段，比如拉取一个接龙或者评论列表的时候，如果需要同时拉取对应的用户资料，可以在API的查询字段添加 include 查询，
比如： /api/discussions?includes=user
得到的返回结果，会包含一个 entities 字段包含了 user 字段
0099
{"code":0,"data":[],"entities":{
    1: { name: "123", id: 1} // id => user
}}

所有的 GET 方法（除了 '/users/self' ）都是不需要token可以公开调用的。所有的写方法都需要token才能调用。

# API 文档

**用户**

1. 注册 `POST /api/register`

```
{
	"username":"ntop",
	"password":"123"
}
```
返回:
```
{
    "id": 3,
    "created_at": 1575950976,
    "updated_at": 1575950976,
    "sign_count": 0,
    "exp_count": 0,
    "login": "ntop",
    "nickname": "",
    "email": "",
    "phone": "",
    "avatar": "",
    "summary": "",
    "blocked_at": 0,
    "silenced_at": 0
}
```

2. 登录 `POST /api/login`
```
{
	"username":"ntop",
	"password":"123"
}
```
返回：
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM3MjcxMDMsInRleHQiOiJudG9wIiwidHlwZSI6InNlc3MifQ.0ms9C9piphls05B6rnNWmDKuI93JCNJJAAQbfw-_k8w",
    "expires_in": 7776000
}


3. 登出 `POST /api/logout`
返回：200

4. 授权(微信)登录 `POST /api/auth`
```
{
    code: 12345
}
```
5.  授权(微信)登录 `POST /api/auth/mp`
同上

6. H5 授权登录 `POST /api/auth/h5`
```
{
    code: 12345
}
```

***用户**
1. 获取一个用户 `GET /api/users/:id`
返回：
{
    "id": 1,
    "created_at": 1575438044,
    "updated_at": 1575947015,
    "sign_count": 0,
    "exp_count": 0,
    "login": "wxOO66FTXNWDIGG===",
    "nickname": "🥛",
    "email": "",
    "phone": "",
    "avatar": "https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJymlMOTQ7Oqsb2C7icsemuVdDJxbyjfntw32tkYpfmYPKgT45oOQJLEVhYjk9UXWcDWzWWO1837gA/132",
    "summary": "",
    "blocked_at": 0,
    "silenced_at": 0
}

2. 获取自身信息 `GET /api/users/self`
返回：
{
    "id": 3,
    "created_at": 1575950976,
    "updated_at": 1575951806,
    "sign_count": 0,
    "exp_count": 0,
    "login": "ntop",
    "nickname": "",
    "email": "",
    "phone": "",
    "avatar": "",
    "summary": "",
    "blocked_at": 0,
    "silenced_at": 0
}

3. 更新用户信息（只能更新自己） `PUT /api/users/1`
{
	"avatar":"123"
}
返回：
{
    "id": 3,
    "created_at": 1575950976,
    "updated_at": 1575952106,
    "sign_count": 0,
    "exp_count": 0,
    "login": "ntop",
    "nickname": "",
    "email": "",
    "phone": "",
    "avatar": "123",
    "summary": "",
    "blocked_at": 0,
    "silenced_at": 0
}
4. 获取该用户的接龙 `GET /api/users/1/discussions`
返回:
{
    "code": 0,
    "data": [],
    "entities": {}
}

5. 返回用户的接龙评论 `GET /api/users/1/posts`
返回：
{
    "code": 0,
    "data": [
        {
            "id": 9,
            "created_at": 1575791739,
            "discussion_id": 9,
            "parent_id": 0,
            "author_id": 1,
            "reply_id": 0,
            "content": "0"
        },

        {
            "id": 3,
            "created_at": 1575791060,
            "discussion_id": 3,
            "parent_id": 0,
            "author_id": 1,
            "reply_id": 0,
            "content": "2"
        }
    ],
    "entities": {}
}

6. 返回用户的点赞 `GET /api/users/1/likes`
返回：
[]


**接龙**

1. 返回所有的接龙 `GET /api/discussions?includes=user`
{
    "code": 0,
    "data": [
        {
            "id": 23,
            "created_at": 1575860661,
            "updated_at": 0,
            "title": "发个接龙来了",
            "content": "这个不是可省的嘛？",
            "author_id": 2,
            "first_post": 0,
            "last_post": 0,
            "comment_count": 0
        },
    ],
    "entities": {
        "users": {
            "2": {
                "id": 2,
                "created_at": 1575860138,
                "updated_at": 0,
                "sign_count": 0,
                "exp_count": 0,
                "login": "wxEJSU4VB4L23UW===",
                "nickname": "",
                "email": "",
                "phone": "",
                "avatar": "",
                "summary": "",
                "blocked_at": 0,
                "silenced_at": 0
            }
        }
    }
}

2. 返回一条接龙 `GET /api/discussions/:id`
{
    "code": 0,
    "data": {
        "id": 22,
        "created_at": 1575860490,
        "updated_at": 1575860490,
        "title": "流量",
        "content": "经济",
        "author_id": 2,
        "first_post": 0,
        "last_post": 0,
        "comment_count": 0
    },
    "entities": {}
}

3. 发起一个接龙 `POST /api/discussions`
```
{
	"content":"如何看待小米股价暴跌!!",
	"title":"小米股价腰斩喽...目前9块"
}
```
返回：
{
    "id": 27,
    "created_at": 1575952903,
    "updated_at": 1575952903,
    "title": "小米股价腰斩喽...目前9块",
    "content": "如何看待小米股价暴跌!!",
    "author_id": 3,
    "first_post": 0,
    "last_post": 0,
    "comment_count": 0
}

2. 删除一个接龙个 `DELETE /api/discussions/:id`
返回 200

3. 更新接龙 `PUT /api/discussions/:id`
```
{
	"title": "标题更新啦",
	"content": "内容也更新啦"
}
```
返回：
```
{
    "id": 7,
    "created_at": 1579070277,
    "updated_at": 1579070277,
    "title": "标题更新啦",
    "content": "内容也更新啦",
    "author_id": 1,
    "first_post": 0,
    "last_post": 0,
    "comment_count": 0
}
```


**接龙评论**
1. 获取接龙下的评论列表 `Get /api/discussions/26/posts?includes=user`
返回：
{
    "code": 0,
    "data": [
        {
            "id": 30,
            "created_at": 1575942274,
            "discussion_id": 26,
            "parent_id": 0,
            "author_id": 1,
            "reply_id": 0,
            "content": "222\n"
        }
    ],
    "entities": {
        "users": {
            "1": {
                "id": 1,
                "created_at": 1575438044,
                "updated_at": 0,
                "sign_count": 0,
                "exp_count": 0,
                "login": "wxOO66FTXNWDIGG===",
                "nickname": "🥛",
                "email": "",
                "phone": "",
                "avatar": "https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJymlMOTQ7Oqsb2C7icsemuVdDJxbyjfntw32tkYpfmYPKgT45oOQJLEVhYjk9UXWcDWzWWO1837gA/132",
                "summary": "",
                "blocked_at": 0,
                "silenced_at": 0
            }
        }
    }
}
2. 获取一条评论 `GET /api/discussions/posts/1`
返回：
{
    "code": 0,
    "data": {
        "id": 1,
        "created_at": 1575790683,
        "discussion_id": 1,
        "parent_id": 0,
        "author_id": 1,
        "reply_id": 0,
        "content": "2"
    },
    "entities": {}
}

3. 发表评论 `POST /api/discussions/posts`
```
{
	"discussion_id":1,
	"author_id": 1,
	"content": "快来评论一下下"，
}
```
返回：
{
    "id": 32,
    "created_at": 1575953472,
    "discussion_id": 1,
    "parent_id": 0,
    "author_id": 3,
    "reply_id": 0,
    "content": "快来评论一下下"
}
4. 删除评论 `DEL  /api/discussions/posts/:id`
返回 200

5. 更新评论 `PUT /api/discussions/posts/:id`
只能更新 content 字段
```
{
	"content": "内容更新啦"
}
```
返回：
{
    "id": 13,
    "created_at": 1579069965,
    "discussion_id": 1,
    "parent_id": 0,
    "author_id": 1,
    "reply_id": 0,
    "content": "内容更新啦"
}

** 标签 **
1. 获取标签列表 `GET /api/tags`
返回：
[]

2. 返回某个标签下的接龙 `GET /api/tags/hello/discussions`
返回:
[]

** 点赞 **
1. 获取点赞列表  `GET /api/posts/:id/likes`
返回:
[]

2.点赞 `POST /api/posts/likes`
返回 200

3. 取消点赞 `DEL /api/posts/:id/likes`
返回 200

4. 获取帖子点赞数 `GET /api/posts/:id/likes/count`
返回
{
    "num": 10
}


** 图片上传API **
固定地址：http://kawaapp.com/x/api/images
采用 multi-part form 格式上传，字段名 file

** 媒体管理 **





