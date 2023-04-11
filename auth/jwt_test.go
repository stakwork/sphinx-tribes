package auth

import (
	"testing"

	"github.com/stakwork/sphinx-tribes/config"
)

func TestInitJwt(t *testing.T) {

	config.InitConfig()
	InitJwt()

	if TokenAuth == nil {
		t.Error("Could not init JWT")
	} else {
		t.Log("JWT inited successfully")
	}
}
