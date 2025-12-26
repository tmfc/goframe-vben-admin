package consts

// Role constants define the system role names.
const (
	RoleSuper = "super"
	RoleAdmin = "admin"
	RoleUser  = "user"
	RoleGuest = "guest"
)

// DefaultRole returns the default role for users without assigned roles.
func DefaultRole() string {
	return RoleSuper
}

// AllRoles returns a list of all valid role names.
func AllRoles() []string {
	return []string{RoleSuper, RoleAdmin, RoleUser, RoleGuest}
}

// IsValidRole checks if the given role name is valid.
func IsValidRole(role string) bool {
	switch role {
	case RoleSuper, RoleAdmin, RoleUser, RoleGuest:
		return true
	}
	return false
}