package model

// SysRolePermissionIn is the input for assigning/removing a single permission to/from a role.
type SysRolePermissionIn struct {
	RoleID       uint `json:"roleId" v:"required"`
	PermissionID uint `json:"permissionId" v:"required"`
}

// SysRolePermissionOut is the output for getting permissions of a role.
type SysRolePermissionOut struct {
	PermissionIDs []uint `json:"permissionIds"`
}

// SysRolePermissionsIn is the input for assigning multiple permissions to a role.
type SysRolePermissionsIn struct {
	RoleID        uint   `json:"roleId" v:"required"`
	PermissionIDs []uint `json:"permissionIds"`
}
