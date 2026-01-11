package service

import (
	"context"
	"encoding/json"
	"strings"

	v1 "backend/api/user/v1"
	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/crypto/bcrypt"
)

var (
	localUser IUser
)

// User returns the user service instance.
func User() IUser {
	return localUser
}

// RegisterUser sets the instance used by user related handlers.
func RegisterUser(i IUser) {
	localUser = i
}

var _ IUser = (*sUser)(nil)

func init() {
	RegisterUser(NewUser())
}

// NewUser creates a new user service instance.
func NewUser() *sUser {
	return &sUser{}
}

// IUser defines the user service interface.
type IUser interface {
	Info(ctx context.Context, token string) (res *v1.UserInfoRes, err error)
	List(ctx context.Context, in model.UserListIn) (out *model.UserListOut, err error)
	Get(ctx context.Context, id int64) (out *model.UserListItem, err error)
	Create(ctx context.Context, in model.UserCreateIn) (id int64, err error)
	Update(ctx context.Context, in model.UserUpdateIn) error
	Delete(ctx context.Context, id int64) error
}

type sUser struct{}

// Info returns the current authenticated user's profile by validating the JWT.
func (s *sUser) Info(ctx context.Context, token string) (res *v1.UserInfoRes, err error) {
	claims, err := parseToken(token)
	if err != nil {
		return nil, err
	}

	userID := gconv.Int64(claims["id"])
	if userID == 0 {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "user id not found in token")
	}

	var user entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns().Id, userID).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, gerror.NewCode(consts.ErrorCodeUserNotFound, "user not found")
	}

	roles := parseRoles(user.Roles)
	if len(roles) == 0 {
		roles = []string{consts.DefaultRole()}
	}

	homePath := user.HomePath
	if homePath == "" {
		homePath = "/dashboard"
	}

	res = &v1.UserInfoRes{
		UserId:   user.Id,
		Username: user.Username,
		RealName: user.RealName,
		Avatar:   user.Avatar,
		Roles:    roles,
		Desc:     user.RealName,
		HomePath: homePath,
		Token:    token,
	}
	return
}

// List returns paginated users filtered by username.
func (s *sUser) List(ctx context.Context, in model.UserListIn) (out *model.UserListOut, err error) {
	out = &model.UserListOut{}
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}
	tenantID := resolveTenantID(ctx)
	m := dao.SysUser.Ctx(ctx).
		Where(dao.SysUser.Columns().TenantId, tenantID).
		WhereNull(dao.SysUser.Columns().DeletedAt)

	if in.Username != "" {
		m = m.WhereLike(dao.SysUser.Columns().Username, "%"+in.Username+"%")
	}

	out.Total, err = m.Count()
	if err != nil {
		return nil, err
	}

	err = m.Page(in.Page, in.PageSize).
		OrderDesc(dao.SysUser.Columns().CreatedAt).
		Fields(
			dao.SysUser.Columns().Id,
			dao.SysUser.Columns().Username,
			dao.SysUser.Columns().RealName,
			dao.SysUser.Columns().Status,
			dao.SysUser.Columns().Roles,
			dao.SysUser.Columns().HomePath,
			dao.SysUser.Columns().Avatar,
			dao.SysUser.Columns().CreatedAt,
			dao.SysUser.Columns().ExtInfo,
			dao.SysUser.Columns().DeptId,
		).
		Scan(&out.Items)
	if err != nil {
		return nil, err
	}
	for _, item := range out.Items {
		item.ExtInfo = parseExtInfo(item.ExtInfo)
	}
	return out, nil
}

// Get returns a single user by id.
func (s *sUser) Get(ctx context.Context, id int64) (out *model.UserListItem, err error) {
	out = &model.UserListItem{}
	tenantID := resolveTenantID(ctx)
	err = dao.SysUser.Ctx(ctx).
		Where(dao.SysUser.Columns().TenantId, tenantID).
		Where(dao.SysUser.Columns().Id, id).
		WhereNull(dao.SysUser.Columns().DeletedAt).
		Fields(
			dao.SysUser.Columns().Id,
			dao.SysUser.Columns().Username,
			dao.SysUser.Columns().RealName,
			dao.SysUser.Columns().Status,
			dao.SysUser.Columns().Roles,
			dao.SysUser.Columns().HomePath,
			dao.SysUser.Columns().Avatar,
			dao.SysUser.Columns().CreatedAt,
			dao.SysUser.Columns().ExtInfo,
			dao.SysUser.Columns().DeptId,
		).
		Scan(out)
	if err != nil {
		return nil, err
	}
	if out.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	out.ExtInfo = parseExtInfo(out.ExtInfo)
	return out, nil
}

// Create creates a new user.
func (s *sUser) Create(ctx context.Context, in model.UserCreateIn) (id int64, err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return 0, err
	}
	tenantID := resolveTenantID(ctx)
	// Check if username already exists
	count, err := dao.SysUser.Ctx(ctx).
		Where(dao.SysUser.Columns().TenantId, tenantID).
		Where(dao.SysUser.Columns().Username, in.Username).
		Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, gerror.NewCodef(gcode.CodeValidationFailed, "Username '%s' already exists", in.Username)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	roles := strings.TrimSpace(in.Roles)
	if roles == "" {
		roles = "[]"
	}

	extInfoStr := marshalExtInfo(in.ExtInfo)
	deptValue := any(in.DeptId)
	if in.DeptId == 0 {
		deptValue = nil
	}

	result, err := dao.SysUser.Ctx(ctx).Data(g.Map{
		"tenant_id": tenantID,
		"username":  in.Username,
		"password":  string(hashedPassword),
		"real_name": in.RealName,
		"status":    in.Status,
		"roles":     roles,
		"home_path": in.HomePath,
		"avatar":    in.Avatar,
		"ext_info":  extInfoStr,
		"dept_id":   deptValue,
	}).Insert()
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// Update updates an existing user; if password is provided, re-hash it.
func (s *sUser) Update(ctx context.Context, in model.UserUpdateIn) error {
	if err := g.Validator().Data(in).Run(ctx); err != nil {
		return err
	}
	tenantID := resolveTenantID(ctx)
	// Ensure username unique except self
	count, err := dao.SysUser.Ctx(ctx).
		Where(dao.SysUser.Columns().TenantId, tenantID).
		Where(dao.SysUser.Columns().Username, in.Username).
		WhereNot(dao.SysUser.Columns().Id, in.ID).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Username '%s' already exists", in.Username)
	}

	roles := strings.TrimSpace(in.Roles)
	if roles == "" {
		roles = "[]"
	}

	extInfoStr := marshalExtInfo(in.ExtInfo)
	deptValue := any(in.DeptId)
	if in.DeptId == 0 {
		deptValue = nil
	}

	updateData := g.Map{
		dao.SysUser.Columns().Username: in.Username,
		dao.SysUser.Columns().RealName: in.RealName,
		dao.SysUser.Columns().Status:   in.Status,
		dao.SysUser.Columns().Roles:    roles,
		dao.SysUser.Columns().HomePath: in.HomePath,
		dao.SysUser.Columns().Avatar:   in.Avatar,
		dao.SysUser.Columns().ExtInfo:  extInfoStr,
		dao.SysUser.Columns().DeptId:   deptValue,
	}
	if in.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updateData[dao.SysUser.Columns().Password] = string(hashedPassword)
	}

	_, err = dao.SysUser.Ctx(ctx).Data(updateData).
		Where(dao.SysUser.Columns().TenantId, tenantID).
		Where(dao.SysUser.Columns().Id, in.ID).
		Update()
	return err
}

// Delete removes a user by id.
func (s *sUser) Delete(ctx context.Context, id int64) error {
	tenantID := resolveTenantID(ctx)
	_, err := dao.SysUser.Ctx(ctx).
		Where(dao.SysUser.Columns().TenantId, tenantID).
		Where(dao.SysUser.Columns().Id, id).
		Delete()
	return err
}

func marshalExtInfo(m map[string]any) string {
	if len(m) == 0 {
		return "{}"
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func parseExtInfo(raw map[string]any) map[string]any {
	// When scanned into map[string]any, it is already decoded.
	if raw != nil && len(raw) > 0 {
		return raw
	}
	return map[string]any{}
}
