// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysMessageReadDao is the data access object for the table sys_message_read.
type SysMessageReadDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  SysMessageReadColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// SysMessageReadColumns defines and stores column names for the table sys_message_read.
type SysMessageReadColumns struct {
	Id        string //
	MessageId string //
	UserId    string //
	ReadAt    string //
}

// sysMessageReadColumns holds the columns for the table sys_message_read.
var sysMessageReadColumns = SysMessageReadColumns{
	Id:        "id",
	MessageId: "message_id",
	UserId:    "user_id",
	ReadAt:    "read_at",
}

// NewSysMessageReadDao creates and returns a new DAO object for table data access.
func NewSysMessageReadDao(handlers ...gdb.ModelHandler) *SysMessageReadDao {
	return &SysMessageReadDao{
		group:    "default",
		table:    "sys_message_read",
		columns:  sysMessageReadColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysMessageReadDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysMessageReadDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysMessageReadDao) Columns() SysMessageReadColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysMessageReadDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysMessageReadDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysMessageReadDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
