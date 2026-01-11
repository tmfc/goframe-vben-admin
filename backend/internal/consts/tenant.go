package consts

// DefaultTenantID is used when no tenant is resolved from context/token.
const DefaultTenantID = "1"

// Context key to bypass tenant scoping in DAO (only for explicit cross-tenant ops).
const CtxKeySkipTenant = "ctxSkipTenant"
