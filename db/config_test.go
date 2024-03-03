package db

import "testing"

var mockConfigBountyRoles = []struct {
	Name string
}{
	{Name: "admin"},
	{Name: "editor"},
	{Name: "viewer"},
}

func TestRolesCheck(t *testing.T) {
	mockRolesMap := make(map[string]string)
	for _, role := range mockConfigBountyRoles {
		mockRolesMap[role.Name] = role.Name
	}

	userWithRoles := []UserRoles{
		{Role: "admin"},
		{Role: "editor"},
	}

	userWithoutRoles := []UserRoles{
		{Role: "viewer"},
	}

	tests := []struct {
		name      string
		rolesMap  map[string]string
		userRoles []UserRoles
		check     string
		want      bool
	}{
		{"RoleExistsAndUserHasRole", mockRolesMap, userWithRoles, "admin", true},
		{"RoleExistsButUserDoesNotHaveRole", mockRolesMap, userWithoutRoles, "admin", false},
		{"RoleDoesNotExist", mockRolesMap, userWithRoles, "nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RolesCheck(tt.rolesMap, tt.userRoles, tt.check); got != tt.want {
				t.Errorf("RolesCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckUser(t *testing.T) {
	userRoles := []UserRoles{
		{OwnerPubKey: "pubkey1"},
		{OwnerPubKey: "pubkey2"},
		{OwnerPubKey: "pubkey3"},
	}

	tests := []struct {
		name      string
		userRoles []UserRoles
		pubkey    string
		want      bool
	}{
		{"UserExists", userRoles, "pubkey1", true},

		{"UserDoesNotExist", userRoles, "pubkeyNonExistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckUser(tt.userRoles, tt.pubkey); got != tt.want {
				t.Errorf("CheckUser() for pubkey %v = %v, want %v", tt.pubkey, got, tt.want)
			}
		})
	}
}
