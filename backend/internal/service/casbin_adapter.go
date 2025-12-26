package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const casbinRuleTable = "casbin_rule"

// CasbinAdapter persists policy rules using GoFrame's database layer.
type CasbinAdapter struct {
	ctx   context.Context
	db    gdb.DB
	table string
}

// NewCasbinAdapter creates a new adapter with the provided context.
func NewCasbinAdapter(ctx context.Context) *CasbinAdapter {
	if ctx == nil {
		ctx = context.Background()
	}
	return &CasbinAdapter{
		ctx:   ctx,
		db:    g.DB(),
		table: casbinRuleTable,
	}
}

// LoadPolicy loads all policy rules from storage.
func (a *CasbinAdapter) LoadPolicy(model model.Model) error {
	records, err := a.db.Ctx(a.ctx).Model(a.table).All()
	if err != nil {
		return err
	}
	for _, record := range records {
		line := recordToPolicyLine(record)
		if line == "" {
			continue
		}
		persist.LoadPolicyLine(line, model)
	}
	return nil
}

// SavePolicy saves all policy rules to storage.
func (a *CasbinAdapter) SavePolicy(model model.Model) error {
	records := policyRecordsFromModel(model)
	return a.db.Transaction(a.ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(a.table).Ctx(ctx).Delete(); err != nil {
			return err
		}
		if len(records) == 0 {
			return nil
		}
		_, err := tx.Model(a.table).Ctx(ctx).Data(records).Insert()
		return err
	})
}

// AddPolicy adds a policy rule to storage.
func (a *CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	if strings.TrimSpace(ptype) == "" {
		return gerror.New("policy type is empty")
	}
	record := policyToRecord(ptype, rule)
	_, err := a.db.Ctx(a.ctx).Model(a.table).Data(record).Insert()
	return err
}

// RemovePolicy removes a policy rule from storage.
func (a *CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	builder := a.db.Ctx(a.ctx).Model(a.table).Where("ptype", ptype)
	for i, value := range rule {
		if i > 5 {
			break
		}
		builder = builder.Where(fmt.Sprintf("v%d", i), value)
	}
	_, err := builder.Delete()
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter.
func (a *CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	if fieldIndex < 0 || fieldIndex > 5 {
		return nil
	}
	builder := a.db.Ctx(a.ctx).Model(a.table).Where("ptype", ptype)
	for i, value := range fieldValues {
		if value == "" {
			continue
		}
		index := fieldIndex + i
		if index > 5 {
			break
		}
		builder = builder.Where(fmt.Sprintf("v%d", index), value)
	}
	_, err := builder.Delete()
	return err
}

func recordToPolicyLine(record gdb.Record) string {
	ptype := strings.TrimSpace(record["ptype"].String())
	if ptype == "" {
		return ""
	}
	values := []string{ptype}
	for i := 0; i <= 5; i++ {
		value := strings.TrimSpace(record[fmt.Sprintf("v%d", i)].String())
		if value != "" {
			values = append(values, value)
		}
	}
	return strings.Join(values, ", ")
}

func policyRecordsFromModel(model model.Model) []map[string]interface{} {
	records := make([]map[string]interface{}, 0)
	for ptype, assertion := range model["p"] {
		for _, rule := range assertion.Policy {
			records = append(records, policyToRecord(ptype, rule))
		}
	}
	for ptype, assertion := range model["g"] {
		for _, rule := range assertion.Policy {
			records = append(records, policyToRecord(ptype, rule))
		}
	}
	return records
}

func policyToRecord(ptype string, rule []string) map[string]interface{} {
	record := map[string]interface{}{
		"ptype": ptype,
	}
	for i, value := range rule {
		if i > 5 {
			break
		}
		record[fmt.Sprintf("v%d", i)] = value
	}
	return record
}
