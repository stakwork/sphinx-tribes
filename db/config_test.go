package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRolesCheck(t *testing.T) {
	userRoles := []UserRoles{
		{Role: "ADD BOUNTY"},
	}

	// test returns true when a user has a role
	assert.True(t, RolesCheck(userRoles, "ADD BOUNTY"))

	// returns false when a use does not have a role
	assert.False(t, RolesCheck(userRoles, "DELETE BOUNTY"))
	assert.False(t, RolesCheck(userRoles, "DELETE BOUNTY2"))
}

func TestCheckUser(t *testing.T) {
	userRoles := []UserRoles{
		{OwnerPubKey: "userPublicKey"},
	}
	// if in the user roles, one of the owner_pubkey belongs to the user return true else return false
	assert.True(t, CheckUser(userRoles, "userPublicKey"))
	assert.False(t, CheckUser(userRoles, "anotherPublicKey"))
}
