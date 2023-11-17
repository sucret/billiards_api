// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"billiards/pkg/mysql/model"
)

func newRoleNode(db *gorm.DB) roleNode {
	_roleNode := roleNode{}

	_roleNode.roleNodeDo.UseDB(db)
	_roleNode.roleNodeDo.UseModel(&model.RoleNode{})

	tableName := _roleNode.roleNodeDo.TableName()
	_roleNode.ALL = field.NewAsterisk(tableName)
	_roleNode.RoleNodeID = field.NewInt32(tableName, "role_node_id")
	_roleNode.RoleID = field.NewInt32(tableName, "role_id")
	_roleNode.NodeID = field.NewInt32(tableName, "node_id")

	_roleNode.fillFieldMap()

	return _roleNode
}

type roleNode struct {
	roleNodeDo roleNodeDo

	ALL        field.Asterisk
	RoleNodeID field.Int32
	RoleID     field.Int32 // 角色ID
	NodeID     field.Int32 // 节点ID

	fieldMap map[string]field.Expr
}

func (r roleNode) Table(newTableName string) *roleNode {
	r.roleNodeDo.UseTable(newTableName)
	return r.updateTableName(newTableName)
}

func (r roleNode) As(alias string) *roleNode {
	r.roleNodeDo.DO = *(r.roleNodeDo.As(alias).(*gen.DO))
	return r.updateTableName(alias)
}

func (r *roleNode) updateTableName(table string) *roleNode {
	r.ALL = field.NewAsterisk(table)
	r.RoleNodeID = field.NewInt32(table, "role_node_id")
	r.RoleID = field.NewInt32(table, "role_id")
	r.NodeID = field.NewInt32(table, "node_id")

	r.fillFieldMap()

	return r
}

func (r *roleNode) WithContext(ctx context.Context) IRoleNodeDo { return r.roleNodeDo.WithContext(ctx) }

func (r roleNode) TableName() string { return r.roleNodeDo.TableName() }

func (r roleNode) Alias() string { return r.roleNodeDo.Alias() }

func (r *roleNode) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := r.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (r *roleNode) fillFieldMap() {
	r.fieldMap = make(map[string]field.Expr, 3)
	r.fieldMap["role_node_id"] = r.RoleNodeID
	r.fieldMap["role_id"] = r.RoleID
	r.fieldMap["node_id"] = r.NodeID
}

func (r roleNode) clone(db *gorm.DB) roleNode {
	r.roleNodeDo.ReplaceDB(db)
	return r
}

type roleNodeDo struct{ gen.DO }

type IRoleNodeDo interface {
	gen.SubQuery
	Debug() IRoleNodeDo
	WithContext(ctx context.Context) IRoleNodeDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	As(alias string) gen.Dao
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IRoleNodeDo
	Not(conds ...gen.Condition) IRoleNodeDo
	Or(conds ...gen.Condition) IRoleNodeDo
	Select(conds ...field.Expr) IRoleNodeDo
	Where(conds ...gen.Condition) IRoleNodeDo
	Order(conds ...field.Expr) IRoleNodeDo
	Distinct(cols ...field.Expr) IRoleNodeDo
	Omit(cols ...field.Expr) IRoleNodeDo
	Join(table schema.Tabler, on ...field.Expr) IRoleNodeDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IRoleNodeDo
	RightJoin(table schema.Tabler, on ...field.Expr) IRoleNodeDo
	Group(cols ...field.Expr) IRoleNodeDo
	Having(conds ...gen.Condition) IRoleNodeDo
	Limit(limit int) IRoleNodeDo
	Offset(offset int) IRoleNodeDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IRoleNodeDo
	Unscoped() IRoleNodeDo
	Create(values ...*model.RoleNode) error
	CreateInBatches(values []*model.RoleNode, batchSize int) error
	Save(values ...*model.RoleNode) error
	First() (*model.RoleNode, error)
	Take() (*model.RoleNode, error)
	Last() (*model.RoleNode, error)
	Find() ([]*model.RoleNode, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.RoleNode, err error)
	FindInBatches(result *[]*model.RoleNode, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.RoleNode) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IRoleNodeDo
	Assign(attrs ...field.AssignExpr) IRoleNodeDo
	Joins(fields ...field.RelationField) IRoleNodeDo
	Preload(fields ...field.RelationField) IRoleNodeDo
	FirstOrInit() (*model.RoleNode, error)
	FirstOrCreate() (*model.RoleNode, error)
	FindByPage(offset int, limit int) (result []*model.RoleNode, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IRoleNodeDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (r roleNodeDo) Debug() IRoleNodeDo {
	return r.withDO(r.DO.Debug())
}

func (r roleNodeDo) WithContext(ctx context.Context) IRoleNodeDo {
	return r.withDO(r.DO.WithContext(ctx))
}

func (r roleNodeDo) ReadDB() IRoleNodeDo {
	return r.Clauses(dbresolver.Read)
}

func (r roleNodeDo) WriteDB() IRoleNodeDo {
	return r.Clauses(dbresolver.Write)
}

func (r roleNodeDo) Clauses(conds ...clause.Expression) IRoleNodeDo {
	return r.withDO(r.DO.Clauses(conds...))
}

func (r roleNodeDo) Returning(value interface{}, columns ...string) IRoleNodeDo {
	return r.withDO(r.DO.Returning(value, columns...))
}

func (r roleNodeDo) Not(conds ...gen.Condition) IRoleNodeDo {
	return r.withDO(r.DO.Not(conds...))
}

func (r roleNodeDo) Or(conds ...gen.Condition) IRoleNodeDo {
	return r.withDO(r.DO.Or(conds...))
}

func (r roleNodeDo) Select(conds ...field.Expr) IRoleNodeDo {
	return r.withDO(r.DO.Select(conds...))
}

func (r roleNodeDo) Where(conds ...gen.Condition) IRoleNodeDo {
	return r.withDO(r.DO.Where(conds...))
}

func (r roleNodeDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IRoleNodeDo {
	return r.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (r roleNodeDo) Order(conds ...field.Expr) IRoleNodeDo {
	return r.withDO(r.DO.Order(conds...))
}

func (r roleNodeDo) Distinct(cols ...field.Expr) IRoleNodeDo {
	return r.withDO(r.DO.Distinct(cols...))
}

func (r roleNodeDo) Omit(cols ...field.Expr) IRoleNodeDo {
	return r.withDO(r.DO.Omit(cols...))
}

func (r roleNodeDo) Join(table schema.Tabler, on ...field.Expr) IRoleNodeDo {
	return r.withDO(r.DO.Join(table, on...))
}

func (r roleNodeDo) LeftJoin(table schema.Tabler, on ...field.Expr) IRoleNodeDo {
	return r.withDO(r.DO.LeftJoin(table, on...))
}

func (r roleNodeDo) RightJoin(table schema.Tabler, on ...field.Expr) IRoleNodeDo {
	return r.withDO(r.DO.RightJoin(table, on...))
}

func (r roleNodeDo) Group(cols ...field.Expr) IRoleNodeDo {
	return r.withDO(r.DO.Group(cols...))
}

func (r roleNodeDo) Having(conds ...gen.Condition) IRoleNodeDo {
	return r.withDO(r.DO.Having(conds...))
}

func (r roleNodeDo) Limit(limit int) IRoleNodeDo {
	return r.withDO(r.DO.Limit(limit))
}

func (r roleNodeDo) Offset(offset int) IRoleNodeDo {
	return r.withDO(r.DO.Offset(offset))
}

func (r roleNodeDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IRoleNodeDo {
	return r.withDO(r.DO.Scopes(funcs...))
}

func (r roleNodeDo) Unscoped() IRoleNodeDo {
	return r.withDO(r.DO.Unscoped())
}

func (r roleNodeDo) Create(values ...*model.RoleNode) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Create(values)
}

func (r roleNodeDo) CreateInBatches(values []*model.RoleNode, batchSize int) error {
	return r.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (r roleNodeDo) Save(values ...*model.RoleNode) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Save(values)
}

func (r roleNodeDo) First() (*model.RoleNode, error) {
	if result, err := r.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoleNode), nil
	}
}

func (r roleNodeDo) Take() (*model.RoleNode, error) {
	if result, err := r.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoleNode), nil
	}
}

func (r roleNodeDo) Last() (*model.RoleNode, error) {
	if result, err := r.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoleNode), nil
	}
}

func (r roleNodeDo) Find() ([]*model.RoleNode, error) {
	result, err := r.DO.Find()
	return result.([]*model.RoleNode), err
}

func (r roleNodeDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.RoleNode, err error) {
	buf := make([]*model.RoleNode, 0, batchSize)
	err = r.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (r roleNodeDo) FindInBatches(result *[]*model.RoleNode, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return r.DO.FindInBatches(result, batchSize, fc)
}

func (r roleNodeDo) Attrs(attrs ...field.AssignExpr) IRoleNodeDo {
	return r.withDO(r.DO.Attrs(attrs...))
}

func (r roleNodeDo) Assign(attrs ...field.AssignExpr) IRoleNodeDo {
	return r.withDO(r.DO.Assign(attrs...))
}

func (r roleNodeDo) Joins(fields ...field.RelationField) IRoleNodeDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Joins(_f))
	}
	return &r
}

func (r roleNodeDo) Preload(fields ...field.RelationField) IRoleNodeDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Preload(_f))
	}
	return &r
}

func (r roleNodeDo) FirstOrInit() (*model.RoleNode, error) {
	if result, err := r.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoleNode), nil
	}
}

func (r roleNodeDo) FirstOrCreate() (*model.RoleNode, error) {
	if result, err := r.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoleNode), nil
	}
}

func (r roleNodeDo) FindByPage(offset int, limit int) (result []*model.RoleNode, count int64, err error) {
	result, err = r.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = r.Offset(-1).Limit(-1).Count()
	return
}

func (r roleNodeDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = r.Count()
	if err != nil {
		return
	}

	err = r.Offset(offset).Limit(limit).Scan(result)
	return
}

func (r roleNodeDo) Scan(result interface{}) (err error) {
	return r.DO.Scan(result)
}

func (r roleNodeDo) Delete(models ...*model.RoleNode) (result gen.ResultInfo, err error) {
	return r.DO.Delete(models)
}

func (r *roleNodeDo) withDO(do gen.Dao) *roleNodeDo {
	r.DO = *do.(*gen.DO)
	return r
}
