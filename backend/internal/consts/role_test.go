package consts

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestRoleConstants(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test role constants are defined
		t.Assert(RoleSuper, "super")
		t.Assert(RoleAdmin, "admin")
		t.Assert(RoleUser, "user")
		t.Assert(RoleGuest, "guest")
	})
}

func TestDefaultRole(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test default role returns RoleSuper
		t.Assert(DefaultRole(), RoleSuper)
	})
}

func TestAllRoles(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		roles := AllRoles()

		// Test all roles are returned
		t.Assert(len(roles), 4)

		// Test each role is included
		t.Assert(contains(roles, RoleSuper), true)
		t.Assert(contains(roles, RoleAdmin), true)
		t.Assert(contains(roles, RoleUser), true)
		t.Assert(contains(roles, RoleGuest), true)
	})
}

func TestIsValidRole(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test valid roles
		t.Assert(IsValidRole(RoleSuper), true)
		t.Assert(IsValidRole(RoleAdmin), true)
		t.Assert(IsValidRole(RoleUser), true)
		t.Assert(IsValidRole(RoleGuest), true)

		// Test invalid roles
		t.Assert(IsValidRole(""), false)
		t.Assert(IsValidRole("invalid"), false)
		t.Assert(IsValidRole("Super"), false) // case sensitive
		t.Assert(IsValidRole("ADMIN"), false) // case sensitive
		t.Assert(IsValidRole("root"), false)
		t.Assert(IsValidRole("moderator"), false)
	})
}

func TestIsValidRole_AllConstants(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test all role constants are valid
		allRoles := AllRoles()
		for _, role := range allRoles {
			t.AssertEQ(IsValidRole(role), true)
		}
	})
}

func TestRoleConstants_Uniqueness(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test all role constants are unique
		roles := []string{RoleSuper, RoleAdmin, RoleUser, RoleGuest}
		uniqueRoles := make(map[string]bool)
		for _, role := range roles {
			t.AssertEQ(uniqueRoles[role], false)
			uniqueRoles[role] = true
		}
	})
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}