package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"testing"
)

func TestUsers(t *testing.T) {
	s := beforeUser()
	defer s.Close()

	user := model.User{
		Login: "ntop",
		Email: "ntop.liu@gmail.com",
	}
	if err := s.CreateUser(&user); err != nil {
		t.Error(err)
	}
	if err := s.UpdateUser(&user); err != nil {
		t.Error(err)
	}
	getuser, err := s.GetUser(user.ID)
	if err != nil {
		t.Error(err)
	}
	if user.ID != getuser.ID {
		t.Error("user id not equal")
	}

	if err := s.DeleteUser(user.ID); err != nil {
		t.Error(err)
	}
	if _, err := s.GetUser(user.ID); err == nil {
		t.Error("user should be deleted")
	}
}

func TestGetUserByLogin(t *testing.T) {
	s := beforeUser()
	defer s.Close()

	user := model.User{
		Login: "ntop",
		Email: "ntop.liu@gmail.com",
	}
	if err := s.CreateUser(&user); err != nil {
		t.Error(err)
	}

	if _, err := s.GetUserByLogin("ntop"); err != nil {
		t.Error(err)
	}
}

func TestGetUserByPhone(t *testing.T) {
	s := beforeUser()
	defer s.Close()

	user := model.User{
		Login: "ntop",
		Phone: "18610342257",
	}
	if err := s.CreateUser(&user); err != nil {
		t.Error(err)
	}

	get, err := s.GetUserByPhone("18610342257");
	if  err != nil {
		t.Error(err)
	}
	if get.Login != user.Login {
		t.Error("get user-phone, expect:", user, "get:", get)
	}
}

func TestBind(t *testing.T) {
	s := beforeUser()
	defer s.Close()

	user := model.User{
		Login: "ntop",
		Email: "ntop.liu@gmail.com",
	}
	if err := s.CreateUser(&user); err != nil {
		t.Error(err)
	}
	bind := model.UserBind{
		BindId: "openid12345",
		UserId: user.ID,
	}
	if err := s.CreateBind(&bind); err != nil {
		t.Error(err)
	}
	if tmp, err := s.GetUserByBind(bind.BindId); err != nil {
		t.Error(err)
	} else {
		if tmp.ID != user.ID || tmp.Login != user.Login {
			t.Error("did not find the correct user!")
		}
	}
	if err := s.DeleteBind(bind.BindId); err != nil {
		t.Error(err)
	}
	if _, err := s.GetUserByBind(bind.BindId); err == nil {
		t.Error("shit, why bind still there??")
	}
}


func TestUnion(t *testing.T) {
	s := beforeUser()
	defer s.Close()

	user := model.User{
		Login: "ntop",
		Email: "ntop.liu@gmail.com",
	}
	bind := model.UserBind{
		BindId: "openid12345",
		UnionId: "unionid123",
		UserId: user.ID,
	}
	if err := s.CreateBindUser(&user, &bind); err != nil {
		t.Error(err)
	}

	if tmp, err := s.GetUserByUnion(bind.UnionId); err != nil {
		t.Error(err)
	} else {
		if tmp.ID != user.ID || tmp.Login != user.Login {
			t.Error("did not find the correct user!")
		}
	}
}


func TestCreateBindUser(t *testing.T) {
	s := beforeUser()
	defer s.Close()

	user := model.User{
		Login: "ntop",
		Email: "ntop.liu@gmail.com",
	}
	bind := &model.UserBind{
		BindId: "openid12345",
	}

	if err := s.CreateBindUser(&user, bind); err != nil {
		t.Error(err)
	}

	get, err := s.GetUserByBind(bind.BindId)
	if err != nil {
		t.Fatal(err)
	}
	if get.Login != user.Login {
		t.Fatal("CreateBindUser err, expect:", user, "get:", get)
	}

	if err := s.DeleteBindUser(user.ID); err != nil {
		t.Fatal(err)
	}

	if _, err := s.GetUserByBind(bind.BindId); err == nil {
		t.Fatal("shit, why bind still there??")
	}

	if _, err := s.GetUser(user.ID); err == nil {
		t.Fatal("shit, why user still there??")
	}
}

// test DeleteBindUser without bind
func TestDeleteBindUser(t *testing.T) {
	s := beforeUser()
	defer s.Close()

	user := model.User{
		Login: "ntop",
		Email: "ntop.liu@gmail.com",
	}
	if err := s.CreateUser(&user); err != nil {
		t.Fatal(err)
	}

	err := s.DeleteBindUser(user.ID)
	if err != nil {
		t.Fatal(err)
	}
}


func TestGetUserByIds(t *testing.T) {
	s := beforeUser()
	defer s.Close()

	users := []model.User{
		{
			Login: "ntop-1",
			Email: "ntop.liu@gmail.com",
		},
		{
			Login: "ntop-2",
			Email: "ntop.liu@gmail.com",
		},
	}
	for i := range users {
		if err := s.CreateUser(&users[i]); err != nil {
			t.Error(err)
		}
	}
	result, err := s.getUserByIds([]int64{users[0].ID, users[1].ID})
	if err != nil {
		t.Error(err)
	}
	for _, user := range users {
		if _, ok := result[user.ID]; !ok {
			t.Error("read != write")
		}
	}
}

func TestGetUserByIds_Empty(t *testing.T) {
	s := beforeUser()
	defer s.Close()

	_, err := s.getUserByIds([]int64{})
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateUserSign(t *testing.T) {
	s := beforeUser()
	defer s.Close()

	user := model.User{
		Login: "Hello",
	}
	if err := s.CreateUser(&user); err != nil {
		t.Error(err)
	}

	if err := s.UpdateUserSign(user.ID, 12); err != nil {
		t.Error(err)
	}

	get, err := s.GetUser(user.ID)
	if err != nil {
		t.Error(err)
	}
	if get.SignCount != 12 {
		t.Error("update user sign_count err, get:", get.SignCount, "expect:", 12)
	}
}


func TestUpdateUserExp(t *testing.T) {
	s := beforeUser()
	defer s.Close()

	user := model.User{
		Login: "Hello",
	}
	if err := s.CreateUser(&user); err != nil {
		t.Error(err)
	}

	if err := s.UpdateUserExp(user.ID, 12); err != nil {
		t.Error(err)
	}

	get, err := s.GetUser(user.ID)
	if err != nil {
		t.Error(err)
	}
	if get.ExpCount != 12 {
		t.Error("update user exp_count err, get:", get.ExpCount, "expect:", 12)
	}
}

func beforeUser() *datasource {
	s := newTest()
	s.Exec("DELETE FROM users;")
	return s
}
