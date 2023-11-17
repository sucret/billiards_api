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

func newTable(db *gorm.DB) table {
	_table := table{}

	_table.tableDo.UseDB(db)
	_table.tableDo.UseModel(&model.Table{})

	tableName := _table.tableDo.TableName()
	_table.ALL = field.NewAsterisk(tableName)
	_table.TableID = field.NewInt32(tableName, "table_id")
	_table.Name = field.NewString(tableName, "name")
	_table.ShopID = field.NewInt32(tableName, "shop_id")
	_table.Status = field.NewInt32(tableName, "status")
	_table.Qrcode = field.NewString(tableName, "qrcode")
	_table.CreatedAt = field.NewField(tableName, "created_at")

	_table.fillFieldMap()

	return _table
}

type table struct {
	tableDo tableDo

	ALL       field.Asterisk
	TableID   field.Int32
	Name      field.String // 球桌名称
	ShopID    field.Int32  // 店铺id
	Status    field.Int32  // 状态，1｜开启，2｜关闭
	Qrcode    field.String // 开台二维码
	CreatedAt field.Field

	fieldMap map[string]field.Expr
}

func (t table) Table(newTableName string) *table {
	t.tableDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t table) As(alias string) *table {
	t.tableDo.DO = *(t.tableDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *table) updateTableName(table string) *table {
	t.ALL = field.NewAsterisk(table)
	t.TableID = field.NewInt32(table, "table_id")
	t.Name = field.NewString(table, "name")
	t.ShopID = field.NewInt32(table, "shop_id")
	t.Status = field.NewInt32(table, "status")
	t.Qrcode = field.NewString(table, "qrcode")
	t.CreatedAt = field.NewField(table, "created_at")

	t.fillFieldMap()

	return t
}

func (t *table) WithContext(ctx context.Context) ITableDo { return t.tableDo.WithContext(ctx) }

func (t table) TableName() string { return t.tableDo.TableName() }

func (t table) Alias() string { return t.tableDo.Alias() }

func (t *table) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *table) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 6)
	t.fieldMap["table_id"] = t.TableID
	t.fieldMap["name"] = t.Name
	t.fieldMap["shop_id"] = t.ShopID
	t.fieldMap["status"] = t.Status
	t.fieldMap["qrcode"] = t.Qrcode
	t.fieldMap["created_at"] = t.CreatedAt
}

func (t table) clone(db *gorm.DB) table {
	t.tableDo.ReplaceDB(db)
	return t
}

type tableDo struct{ gen.DO }

type ITableDo interface {
	gen.SubQuery
	Debug() ITableDo
	WithContext(ctx context.Context) ITableDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	As(alias string) gen.Dao
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ITableDo
	Not(conds ...gen.Condition) ITableDo
	Or(conds ...gen.Condition) ITableDo
	Select(conds ...field.Expr) ITableDo
	Where(conds ...gen.Condition) ITableDo
	Order(conds ...field.Expr) ITableDo
	Distinct(cols ...field.Expr) ITableDo
	Omit(cols ...field.Expr) ITableDo
	Join(table schema.Tabler, on ...field.Expr) ITableDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ITableDo
	RightJoin(table schema.Tabler, on ...field.Expr) ITableDo
	Group(cols ...field.Expr) ITableDo
	Having(conds ...gen.Condition) ITableDo
	Limit(limit int) ITableDo
	Offset(offset int) ITableDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ITableDo
	Unscoped() ITableDo
	Create(values ...*model.Table) error
	CreateInBatches(values []*model.Table, batchSize int) error
	Save(values ...*model.Table) error
	First() (*model.Table, error)
	Take() (*model.Table, error)
	Last() (*model.Table, error)
	Find() ([]*model.Table, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Table, err error)
	FindInBatches(result *[]*model.Table, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Table) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ITableDo
	Assign(attrs ...field.AssignExpr) ITableDo
	Joins(fields ...field.RelationField) ITableDo
	Preload(fields ...field.RelationField) ITableDo
	FirstOrInit() (*model.Table, error)
	FirstOrCreate() (*model.Table, error)
	FindByPage(offset int, limit int) (result []*model.Table, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ITableDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (t tableDo) Debug() ITableDo {
	return t.withDO(t.DO.Debug())
}

func (t tableDo) WithContext(ctx context.Context) ITableDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t tableDo) ReadDB() ITableDo {
	return t.Clauses(dbresolver.Read)
}

func (t tableDo) WriteDB() ITableDo {
	return t.Clauses(dbresolver.Write)
}

func (t tableDo) Clauses(conds ...clause.Expression) ITableDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t tableDo) Returning(value interface{}, columns ...string) ITableDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t tableDo) Not(conds ...gen.Condition) ITableDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t tableDo) Or(conds ...gen.Condition) ITableDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t tableDo) Select(conds ...field.Expr) ITableDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t tableDo) Where(conds ...gen.Condition) ITableDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t tableDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) ITableDo {
	return t.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (t tableDo) Order(conds ...field.Expr) ITableDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t tableDo) Distinct(cols ...field.Expr) ITableDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t tableDo) Omit(cols ...field.Expr) ITableDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t tableDo) Join(table schema.Tabler, on ...field.Expr) ITableDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t tableDo) LeftJoin(table schema.Tabler, on ...field.Expr) ITableDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t tableDo) RightJoin(table schema.Tabler, on ...field.Expr) ITableDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t tableDo) Group(cols ...field.Expr) ITableDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t tableDo) Having(conds ...gen.Condition) ITableDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t tableDo) Limit(limit int) ITableDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t tableDo) Offset(offset int) ITableDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t tableDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ITableDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t tableDo) Unscoped() ITableDo {
	return t.withDO(t.DO.Unscoped())
}

func (t tableDo) Create(values ...*model.Table) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t tableDo) CreateInBatches(values []*model.Table, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t tableDo) Save(values ...*model.Table) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t tableDo) First() (*model.Table, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Table), nil
	}
}

func (t tableDo) Take() (*model.Table, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Table), nil
	}
}

func (t tableDo) Last() (*model.Table, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Table), nil
	}
}

func (t tableDo) Find() ([]*model.Table, error) {
	result, err := t.DO.Find()
	return result.([]*model.Table), err
}

func (t tableDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Table, err error) {
	buf := make([]*model.Table, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t tableDo) FindInBatches(result *[]*model.Table, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t tableDo) Attrs(attrs ...field.AssignExpr) ITableDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t tableDo) Assign(attrs ...field.AssignExpr) ITableDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t tableDo) Joins(fields ...field.RelationField) ITableDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t tableDo) Preload(fields ...field.RelationField) ITableDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t tableDo) FirstOrInit() (*model.Table, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Table), nil
	}
}

func (t tableDo) FirstOrCreate() (*model.Table, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Table), nil
	}
}

func (t tableDo) FindByPage(offset int, limit int) (result []*model.Table, count int64, err error) {
	result, err = t.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = t.Offset(-1).Limit(-1).Count()
	return
}

func (t tableDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t tableDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t tableDo) Delete(models ...*model.Table) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *tableDo) withDO(do gen.Dao) *tableDo {
	t.DO = *do.(*gen.DO)
	return t
}
