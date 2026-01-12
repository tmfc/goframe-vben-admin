package service

import (
	"testing"

	"backend/internal/consts"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestParseLoginIdentity(t *testing.T) {
	cases := []struct {
		name       string
		input      string
		multi      bool
		wantUser   string
		wantCode   string
		wantTenant string
		wantSuffix bool
	}{
		{
			name:       "single tenant without suffix",
			input:      "admin",
			multi:      false,
			wantUser:   "admin",
			wantTenant: consts.DefaultTenantID,
			wantSuffix: false,
		},
		{
			name:       "single tenant with suffix",
			input:      "admin@system",
			multi:      false,
			wantUser:   "admin",
			wantTenant: consts.DefaultTenantID,
			wantSuffix: true,
		},
		{
			name:       "multi tenant without suffix",
			input:      "admin",
			multi:      true,
			wantUser:   "admin",
			wantTenant: consts.DefaultTenantID,
			wantSuffix: false,
		},
		{
			name:       "multi tenant with suffix",
			input:      "admin@system",
			multi:      true,
			wantUser:   "admin",
			wantCode:   "system",
			wantTenant: "",
			wantSuffix: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				got := parseLoginIdentity(tc.input, tc.multi)
				t.Assert(got.Username, tc.wantUser)
				t.Assert(got.TenantCode, tc.wantCode)
				t.Assert(got.TenantID, tc.wantTenant)
				t.Assert(got.HasSuffix, tc.wantSuffix)
			})
		})
	}
}
