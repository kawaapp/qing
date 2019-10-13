package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hash := HashPassword("123")
	if err := CheckPassword("123", hash); err != nil {
		t.Error(err)
	}
}

func TestCheckPassword(t *testing.T)  {
	err := CheckPassword("1234", "$2a$10$4FqTyAB6YFSrWw4g/E8am.376PVezhzi/fjSB.Obd1OVZJq8Sj2b2")
	if err != nil {
		t.Error(err)
	}
	err = CheckPassword("1234", "$2a$10$DzPP.4r7fFyeXyKiY0t8Oev.gz8FJd7kxIRSh.2Kr5Tu6dN0jj39.")
	if err != nil {
		t.Error(err)
	}
}