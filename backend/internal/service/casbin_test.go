package service

import (
	"testing"

	"backend/internal/consts"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestNormalizeDomain(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(NormalizeDomain(""), "default")
		t.AssertEQ(NormalizeDomain("  "), "default")
		t.AssertEQ(NormalizeDomain("tenant1"), "tenant1")
		// NormalizeDomain doesn't trim, it only checks if empty
		t.AssertEQ(NormalizeDomain("  tenant2  "), "  tenant2  ")
	})
}

func TestBuildAccessCodes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test with empty roles (should use default role)
		codes := buildAccessCodes([]string{})
		t.AssertNE(len(codes), 0)
		t.AssertIN("System:Menu:List", codes)

		// Test with super role
		codes = buildAccessCodes([]string{consts.RoleSuper})
		t.AssertIN("System:Menu:List", codes)
		t.AssertIN("System:Menu:Create", codes)
		t.AssertIN("System:Menu:Edit", codes)
		t.AssertIN("System:Menu:Delete", codes)

		// Test with admin role
		codes = buildAccessCodes([]string{consts.RoleAdmin})
		t.AssertIN("System:Menu:List", codes)
		t.AssertIN("System:Menu:Edit", codes)
		t.AssertNE("System:Menu:Create", codes)

		// Test with user role
		codes = buildAccessCodes([]string{consts.RoleUser})
		t.AssertIN("System:Menu:List", codes)
		t.AssertIN("System:Dept:List", codes)
		t.AssertNE("System:Menu:Edit", codes)

		// Test with guest role
		codes = buildAccessCodes([]string{consts.RoleGuest})
		t.AssertIN("System:Menu:List", codes)
		t.AssertNE("System:Dept:List", codes)

		// Test with multiple roles (should merge access codes)
		codes = buildAccessCodes([]string{consts.RoleUser, consts.RoleAdmin})
		t.AssertIN("System:Menu:List", codes)
		t.AssertIN("System:Menu:Edit", codes)
		t.AssertIN("System:Dept:List", codes)

		// Test with invalid role (should use default)
		codes = buildAccessCodes([]string{"invalid_role"})
		t.AssertIN("System:Menu:List", codes)
	})
}

func TestBuildAccessCodesSorting(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		codes := buildAccessCodes([]string{consts.RoleSuper, consts.RoleAdmin})
		// Verify codes are sorted
		for i := 1; i < len(codes); i++ {
			t.AssertLT(codes[i-1], codes[i])
		}
	})
}