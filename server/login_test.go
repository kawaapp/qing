package server

import "testing"

func TestGenerateUserId(t *testing.T) {
	t.Error(GenerateUserId("wx"))
}
