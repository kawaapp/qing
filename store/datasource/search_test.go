package datasource

import (
	"testing"
	"github.com/kawaapp/kawaqing/model"
	"sort"
)

func TestSQLUserQuery(t *testing.T) {
	params := make(map[string]string, 0)
	page, size := 0, 20

	// search nothing...
	query, args := sqlUserQuery("SELECT * ", params, page, size)
	if query != "SELECT *  FROM users ORDER BY id DESC LIMIT 20 OFFSET 0" {
		t.Fatal("sqlUserQuery query err:", query)
	}
	if len(args) != 0 {
		t.Fatal("sqlUserQuery args err:", len(args))
	}

	// search login, nickname, silence, block
	params["login"] = "foobar"
	params["nickname"] = "fb"
	params["silence"] = ""
	params["block"] = ""

	query, args = sqlUserQuery("SELECT *", params, page, size)
	if query != "SELECT * FROM users WHERE 1=1 AND login=? AND nickname LIKE ? AND silenced_at > ? AND blocked_at > 0 ORDER BY id DESC LIMIT 20 OFFSET 0" {
		t.Fatal("sqlUserQuery query err:", query)
	}
	if len(args) != 2 {
		t.Fatal("sqlUserQuery args err:", args)
	}

	// count
	query, args = sqlUserQuery("SELECT COUNT(id)", params, 0, 0)
	if query != "SELECT COUNT(id) FROM users WHERE 1=1 AND login=? AND nickname LIKE ? AND silenced_at > ? AND blocked_at > 0" {
		t.Fatal("sqlUserQuery count(id) err:", query)
	}
}

func TestSearchUser(t *testing.T) {
	s := beforeSearch()
	defer s.Close()

	// create user
	users := []model.User{
		{
			Login: "ntop",
			Nickname:"ntop",
		}, {
			Login: "tom",
			Nickname: "tom",
		}, {
			Login: "lily",
			Nickname: "lily",
		},
	}
	if err := createUserList(s, users); err != nil {
		t.Fatal(err)
	}

	// search all
	params := make(map[string]string, 0)
	gets, err := s.SearchUser(params, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(gets) != len(users) {
		t.Fatal("user search total, expect:", len(users), "get:", len(gets))
	}

	// search login
	params = map[string]string {
		"login": "tom",
	}
	gets, err = s.SearchUser(params, 0, 100)
	if err != nil {
		t.Fatal(err)
	}

	if sz := len(gets); sz != 1 {
		t.Fatal("user search result, expect:", 1, "get:", sz)
	} else if tom := gets[0]; tom.Login != "tom" {
		t.Fatal("user search result, expect:", "tom", "get:", tom.Login)
	}

	// search keyword, get ntop, tom
	params = map[string]string {
		"nickname": "to",
	}
	gets, err = s.SearchUser(params, 0, 100)
	if err != nil {
		t.Error(err)
	}
	if sz := len(gets); sz != 2 {
		t.Error("user search keyword, expect:", 2, "get:", sz)
	}
}


func TestSearchUserCount(t *testing.T) {
	s := beforeSearch()
	defer s.Close()

	// create user
	users := []model.User{
		{
			Login: "ntop",
			Nickname:"ntop",
		}, {
			Login: "tom",
			Nickname: "tom",
		}, {
			Login: "lily",
			Nickname: "lily",
		},
	}
	if err := createUserList(s, users); err != nil {
		t.Error(err)
	}

	params := make(map[string]string, 0)
	num, err := s.SearchUserCount(params)
	if err != nil {
		t.Fatal(err)
	}
	if num != len(users) {
		t.Fatal("count user size, expect:", len(users), "get:", num)
	}
}

func TestSQLDiscussionQuery(t *testing.T) {
	params := make(map[string]string, 0)
	page, size := 0, 20

	// search nothing...
	query, args := sqlDiscussionQuery("SELECT *", params, page, size)
	if query != "SELECT * FROM discussions ORDER BY id DESC LIMIT 20 OFFSET 0" {
		t.Error("sqlDiscussionQuery query err:", query)
	}
	if len(args) != 0 {
		t.Fatal("sqlDiscussionQuery args err:", len(args))
	}

	// search login, nickname, silence, block
	params["content"] = "foobar"
	params["author"] = "fb"

	query, args = sqlDiscussionQuery("SELECT *", params, page, size)
	if query != "SELECT * FROM discussions WHERE 1=1 AND content LIKE ? AND author_id IN (SELECT id FROM users WHERE nickname=?) ORDER BY id DESC LIMIT 20 OFFSET 0" {
		t.Fatal("sqlDiscussionQuery query err:", query)
	}
	if len(args) != 2 {
		t.Fatal("sqlDiscussionQuery args err:", args)
	}

	// count
	query, args = sqlDiscussionQuery("SELECT COUNT(id)", params, 0, 0)
	if query != "SELECT COUNT(id) FROM discussions WHERE 1=1 AND content LIKE ? AND author_id IN (SELECT id FROM users WHERE nickname=?)" {
		t.Fatal("sqlDiscussionQuery count(id) err:", query)
	}
}

func TestSearchDiscussion(t *testing.T) {
	s := beforeSearch()
	defer s.Close()

	// create user
	user := model.User{
		Login:"ntop",
		Nickname:"ntoooop",
	}
	if err := s.CreateUser(&user); err != nil {
		t.Error(err)
	}

	// create discussion
	discussions := []model.Discussion {
		{
			Content: "Hello",
			AuthorID: user.ID,
		}, {
			Content: "World",
			AuthorID: user.ID,
		}, {
			Content: "FooBar",
			AuthorID: 1234567,
		},
	}
	if err := createDiscussionList(s, discussions); err != nil {
		t.Error(err)
	}

	// Search ALL
	params := make(map[string]string, 0)
	gets, err := s.SearchDiscussion(params, 0, 0)
	if err != nil {
		t.Error(err)
	}
	if len(gets) != len(discussions) {
		t.Error("post search all, expect:", len(discussions), "get:", len(gets))
	}

	// Search Login
	params = map[string]string {
		"author": user.Nickname,
	}
	gets, err = s.SearchDiscussion(params, 0, 100)
	if err != nil {
		t.Error(err)
	}
	if len(gets) != 2 {
		t.Error("post search by login, expect:", 2, "get:", len(gets))
	}

	// Search keyword
	params = map[string]string {
		"content": "foo",
	}
	gets, err = s.SearchDiscussion(params, 0, 100)
	if err != nil {
		t.Error(err)
	}
	if len(gets) != 1 {
		t.Error("post search by keyword, expect:", 1, "get:", len(gets))
	}
}

func TestSearchDiscussionCount(t *testing.T) {
	s := beforeSearch()
	defer s.Close()

	// create user
	user := model.User{
		Login:"ntop",
		Nickname:"ntoooop",
	}
	if err := s.CreateUser(&user); err != nil {
		t.Error(err)
	}

	// create post
	discussions := []model.Discussion {
		{
			Content: "Hello",
			AuthorID: user.ID,
		}, {
			Content: "World",
			AuthorID: user.ID,
		}, {
			Content: "FooBar",
			AuthorID: 1234567,
		},
	}
	if err := createDiscussionList(s, discussions); err != nil {
		t.Error(err)
	}

	params := make(map[string]string, 0)

	// Count
	if num, err := s.SearchDiscussionCount(params); err != nil {
		t.Error(err)
	} else if num != len(discussions) {
		t.Fatal("count discussion size, expect:", len(discussions), "get:", num)
	}
}


func TestSQLPostQuery(t *testing.T) {
	params := make(map[string]string, 0)
	page, size := 0, 20

	// search nothing...
	query, args := sqlPostQuery("SELECT *", params, page, size)
	if query != "SELECT * FROM posts ORDER BY id DESC LIMIT 20 OFFSET 0" {
		t.Fatal("sqlPostQuery query err:", query)
	}
	if len(args) != 0 {
		t.Fatal("sqlPostQuery args err:", len(args))
	}

	// search login, nickname, silence, block
	params["content"] = "foobar"
	params["author"] = "fb"

	query, args = sqlPostQuery("SELECT *", params, page, size)
	if query != "SELECT * FROM posts WHERE 1=1 AND content LIKE ? AND author_id IN (SELECT id FROM users WHERE nickname=?) ORDER BY id DESC LIMIT 20 OFFSET 0" {
		t.Fatal("sqlPostQuery query err:", query)
	}
	if len(args) != 2 {
		t.Fatal("sqlPostQuery args err:", args)
	}

	// count
	query, args = sqlPostQuery("SELECT COUNT(id)", params, 0, 0)
	if query != "SELECT COUNT(id) FROM posts WHERE 1=1 AND content LIKE ? AND author_id IN (SELECT id FROM users WHERE nickname=?)" {
		t.Fatal("sqlPostQuery count(id) err:", query)
	}
}


func TestSearchPost(t *testing.T) {
	s := beforeSearch()
	defer s.Close()

	// create user
	user := model.User{
		Login:"ntop",
		Nickname:"ntoooop",
	}
	if err := s.CreateUser(&user); err != nil {
		t.Error(err)
	}

	// create post
	posts := []model.Post{
		{
			Content: "Hello",
			AuthorID: user.ID,
		}, {
			Content: "World",
			AuthorID: user.ID,
		}, {
			Content: "FooBar",
			AuthorID: 1234567,
		},
	}
	if err := createPostList(s, posts); err != nil {
		t.Error(err)
	}

	params := make(map[string]string, 0)

	// Search ALL
	gets, err := s.SearchPost(params, 0, 100)
	if err != nil {
		t.Error(err)
	}
	if len(gets) != len(posts) {
		t.Error("posts search all, expect:", len(posts), "get:", len(gets))
	}

	// Search Login
	params = map[string]string {
		"author": user.Nickname,
	}
	gets, err = s.SearchPost(params, 0, 100)
	if err != nil {
		t.Error(err)
	}
	if len(gets) != 2 {
		t.Error("posts search by login, expect:", 2, "get:", len(gets))
	}

	// Search keyword
	params = map[string]string {
		"content": "foo",
	}
	gets, err = s.SearchPost(params, 0, 100)
	if err != nil {
		t.Error(err)
	}
	if len(gets) != 1 {
		t.Error("posts search by keyword, expect:", 1, "get:", len(gets))
	}
}

func TestSearchPostCount(t *testing.T) {
	s := beforeSearch()
	defer s.Close()

	// create user
	user := model.User{
		Login:"ntop",
		Nickname:"ntoooop",
	}
	if err := s.CreateUser(&user); err != nil {
		t.Error(err)
	}

	// create post
	posts := []model.Post{
		{
			Content: "Hello",
			AuthorID: user.ID,
		}, {
			Content: "World",
			AuthorID: user.ID,
		}, {
			Content: "FooBar",
			AuthorID: 1234567,
		},
	}
	if err := createPostList(s, posts); err != nil {
		t.Error(err)
	}

	params := make(map[string]string, 0)
	num, err := s.SearchPostCount(params)
	if err != nil {
		t.Fatal(err)
	}
	if num != len(posts) {
		t.Fatal("count discussion size, expect:", len(posts), "get:", num)
	}
}

func TestSQLReportQuery(t *testing.T) {
	params := make(map[string]string, 0)
	page, size := 0, 20

	// search nothing...
	query, args := sqlReportQuery("SELECT *", params, page, size)
	expect :=
		"SELECT * FROM reports c " +
		"LEFT JOIN posts a ON a.id = c.entity_id " +  "LEFT JOIN users b ON b.id = c.user_id " +
		"ORDER BY id DESC LIMIT 20 OFFSET 0"
	if query != expect {
		t.Fatal("sqlReportQuery query err:", query)
	}
	if len(args) != 0 {
		t.Fatal("sqlReportQuery args err:", len(args))
	}

	// search login, nickname, silence, block
	params["post"] = "foobar"
	params["user"] = "fb"
	expect =
		"SELECT * FROM reports c " +
		"LEFT JOIN posts a ON a.id = c.entity_id LEFT JOIN users b ON b.id = c.user_id " +
		"WHERE 1=1 AND a.content LIKE ? AND b.nickname LIKE ? " +
		"ORDER BY id DESC LIMIT 20 OFFSET 0"

	query, args = sqlReportQuery("SELECT *", params, page, size)
	if query != expect	{
		t.Fatal("sqlReportQuery query err:", query)
	}
	if len(args) != 2 {
		t.Fatal("sqlReportQuery args err:", args)
	}

	// count
	query, args = sqlReportQuery("SELECT COUNT(id)", params, 0, 0)
	expect =
		"SELECT COUNT(id) FROM reports c " +
		"LEFT JOIN posts a ON a.id = c.entity_id LEFT JOIN users b ON b.id = c.user_id " +
		"WHERE 1=1 AND a.content LIKE ? AND b.nickname LIKE ?"
	if query != expect {
		t.Fatal("sqlReportQuery count(id) err:", query)
	}
}


func TestSearchReport(t *testing.T) {
	s := beforeSearch()
	defer s.Close()

	// create user
	users := []model.User{
		{
			Login: "ntop",
			Nickname:"ntop",
		}, {
			Login: "tom",
			Nickname: "tom",
		}, {
			Login: "lily",
			Nickname: "lily",
		},
	}
	if err := createUserList(s, users); err != nil {
		t.Error(err)
	}

	// create post
	posts := []model.Post {
		{
			Content: "Hello",
			AuthorID: 1,
		}, {
			Content: "World",
			AuthorID: 1,
		}, {
			Content: "FooBar",
			AuthorID: 1234567,
		},
	}
	if err := createPostList(s, posts); err != nil {
		t.Error(err)
	}

	// create reports
	reports := []model.Report {
		{
			EntityId: posts[0].ID,
			UserId: users[0].ID,
		},
		{
			EntityId: posts[0].ID,
			UserId: users[1].ID,
		},
		{
			EntityId: posts[1].ID,
			UserId: users[2].ID,
			Status: 1,
		},
		{
			EntityId: posts[2].ID,
			UserId: users[2].ID,
		},
	}
	if err := createReportList(s, reports); err != nil {
		t.Fatal(err)
	}

	// 0. Search ALL, status=0,  expected: report[0], report[2], report[3]
	params := make(map[string]string, 0)
	gets, err := s.SearchReport(params, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(gets) != len(reports) {
		t.Fatal("report search all, expect:", len(reports), "get:", len(gets))
	}

	// 1. search keyword: 'hello' && status=0, get reports[0]
	params = map[string]string {
		"post": "World",
	}

	gets, err = s.SearchReport(params, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(gets) != 1 {
		t.Fatal("report search keyword 'hello', expect:", 1, "get:", len(gets))
	}

	// 2. search user: lily, expected: report[2], report[3]
	params = map[string]string {
		"user": "lily",
	}
	gets, err = s.SearchReport(params, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(gets) != 2 {
		t.Fatal("report search user 'lily', expect:", 2, "get:", len(gets))
	}

	// 3. search user: lily and keyword: foo, expect: report[2]
	params = map[string]string {
		"user": "lily",
		"post": "foo",

	}
	gets, err = s.SearchReport(params, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(gets) != 1 {
		t.Fatal("report search user 'lily' keyword 'foo', expect:", 1, "get:", len(gets))
	}
	if gets[0].ID != reports[3].ID {
		t.Fatal("report search user 'lily' keyword 'foo', expect:", reports[2], gets[0])
	}
}


func TestSearchSignUser(t *testing.T) {
	s := beforeSearch()
	defer s.Close()

	// create user
	users := []model.User{
		{
			Login: "ntop",
			SignCount: 2,
		}, {
			Login: "tom",
			SignCount: 1,
		}, {
			Login: "lily",
			SignCount: 3,
		},
	}
	if err := createUserList(s, users); err != nil {
		t.Error(err)
	}

	gets, err := s.SearchSignUser(1, 100)
	if err != nil {
		t.Error(err)
	}
	sorted := sort.SliceIsSorted(gets, func(i, j int) bool {
		return gets[i].SignCount > gets[j].SignCount
	})
	if !sorted {
		t.Error("search sign users is not sorted!")
	}
}

func createUserList(s *datasource, users []model.User) error {
	for i := range users {
		if err := s.CreateUser(&users[i]); err != nil {
			return err
		}
	}
	return nil
}

func createDiscussionList(s *datasource, discussions []model.Discussion) error {
	for i := range discussions {
		if err := s.CreateDiscussion(&discussions[i]); err != nil {
			return  err
		}
	}
	return nil
}

func createPostList(s *datasource, posts []model.Post) error {
	for i := range posts {
		if err := s.CreatePost(&posts[i]); err != nil {
			return err
		}
	}
	return nil
}

func createLikeList(s *datasource, likes []model.Like) error {
	for i := range likes {
		if err := s.CreateLike(&likes[i]); err != nil {
			return err
		}
	}
	return nil
}

func createReportList(s *datasource, reports []model.Report) error {
	for i := range reports {
		if err := s.CreateReport(&reports[i]); err != nil {
			return err
		}
	}
	return nil
}

func beforeSearch() *datasource {
	s := newTest()
	s.Exec("DELETE FROM users")
	s.Exec("DELETE FROM discussions")
	s.Exec("DELETE FROM posts")
	s.Exec("DELETE FROM likes")
	s.Exec("DELETE FROM reports")
	return s
}


