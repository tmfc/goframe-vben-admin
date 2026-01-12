// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysMessageDao is the data access object for the table sys_message.
type SysMessageDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysMessageColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysMessageColumns defines and stores column names for the table sys_message.
type SysMessageColumns struct {
	Id         string //
	TenantId   string //
	Title      string //
	Content    string //
	MsgType    string //
	SenderId   string //
	TargetType string //
	Recipients string //
	Path       string //
	Params     string //
	CreatedAt  string //
	UpdatedAt  string //
	DeletedAt  string //
	JumpType   string // 0: None, 1: Route, 2: External
}

// sysMessageColumns holds the columns for the table sys_message.
var sysMessageColumns = SysMessageColumns{
	Id:         "id",
	TenantId:   "tenant_id",
	Title:      "title",
	Content:    "content",
	MsgType:    "msg_type",
	SenderId:   "sender_id",
	TargetType: "target_type",
	Recipients: "recipients",
	Path:       "path",
	Params:     "params",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
	JumpType:   "jump_type",
}

// NewSysMessageDao creates and returns a new DAO object for table data access.
func NewSysMessageDao(handlers ...gdb.ModelHandler) *SysMessageDao {
	return &SysMessageDao{
		group:    "default",
		table:    "sys_message",
		columns:  sysMessageColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysMessageDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysMessageDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysMessageDao) Columns() SysMessageColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysMessageDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysMessageDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *SysMessageDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
