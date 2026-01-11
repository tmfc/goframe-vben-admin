package model

// UserPermissionsIn is the input for getting permissions of a user.
type UserPermissionsIn struct {
	UserID int64 `json:"userId" v:"required"`
}

// UserPermissionsOut is the output for getting permissions of a user.
type UserPermissionsOut struct {
	PermissionIDs []uint `json:"permissionIds"`
}
