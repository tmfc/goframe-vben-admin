package model

import "github.com/gogf/gf/v2/os/gtime"

// TenantCreateIn is the input for creating a tenant.
type TenantCreateIn struct {
	Name           string      `json:"name" v:"required"`
	Code           string      `json:"code" v:"required"`
	ContactName    string      `json:"contactName"`
	ContactMobile  string      `json:"contactMobile"`
	StartAt        *gtime.Time `json:"startAt"`
	ExpireAt       *gtime.Time `json:"expireAt"`
	PackageVersion string      `json:"packageVersion"`
	Domain         string      `json:"domain"`
	Status         int         `json:"status" d:"1"`
}

// TenantUpdateIn is the input for updating a tenant.
type TenantUpdateIn struct {
	ID             int64       `json:"id" v:"required"`
	Name           string      `json:"name" v:"required"`
	Code           string      `json:"code" v:"required"`
	ContactName    string      `json:"contactName"`
	ContactMobile  string      `json:"contactMobile"`
	StartAt        *gtime.Time `json:"startAt"`
	ExpireAt       *gtime.Time `json:"expireAt"`
	PackageVersion string      `json:"packageVersion"`
	Domain         string      `json:"domain"`
	Status         int         `json:"status"`
}

// TenantListIn is the input for listing tenants.
type TenantListIn struct {
	Page     int    `json:"page" d:"1"`
	PageSize int    `json:"pageSize" d:"10"`
	Name     string `json:"name"`
	Code     string `json:"code"`
}

// TenantListItem is the representation of a tenant in list.
type TenantListItem struct {
	Id             int64       `json:"id"`
	Name           string      `json:"name"`
	Code           string      `json:"code"`
	Status         int         `json:"status"`
	ContactName    string      `json:"contactName"`
	ContactMobile  string      `json:"contactMobile"`
	StartAt        *gtime.Time `json:"startAt"`
	ExpireAt       *gtime.Time `json:"expireAt"`
	PackageVersion string      `json:"packageVersion"`
	Domain         string      `json:"domain"`
	CreatedAt      *gtime.Time `json:"createdAt"`
	UpdatedAt      *gtime.Time `json:"updatedAt"`
}

// TenantListOut is the output for listing tenants.
type TenantListOut struct {
	Items []*TenantListItem `json:"items"`
	Total int               `json:"total"`
}
