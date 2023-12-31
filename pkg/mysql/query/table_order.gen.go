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

func newTableOrder(db *gorm.DB) tableOrder {
	_tableOrder := tableOrder{}

	_tableOrder.tableOrderDo.UseDB(db)
	_tableOrder.tableOrderDo.UseModel(&model.TableOrder{})

	tableName := _tableOrder.tableOrderDo.TableName()
	_tableOrder.ALL = field.NewAsterisk(tableName)
	_tableOrder.TableOrderID = field.NewInt32(tableName, "table_order_id")
	_tableOrder.OrderID = field.NewInt32(tableName, "order_id")
	_tableOrder.UserID = field.NewInt32(tableName, "user_id")
	_tableOrder.Status = field.NewInt(tableName, "status")
	_tableOrder.ShopID = field.NewInt32(tableName, "shop_id")
	_tableOrder.TableID = field.NewInt32(tableName, "table_id")
	_tableOrder.CouponID = field.NewInt32(tableName, "coupon_id")
	_tableOrder.UserCouponID = field.NewInt32(tableName, "user_coupon_id")
	_tableOrder.Amount = field.NewInt32(tableName, "amount")
	_tableOrder.PayAmount = field.NewInt32(tableName, "pay_amount")
	_tableOrder.CreatedAt = field.NewField(tableName, "created_at")
	_tableOrder.UpdatedAt = field.NewField(tableName, "updated_at")
	_tableOrder.StartedAt = field.NewField(tableName, "started_at")
	_tableOrder.TerminatedAt = field.NewField(tableName, "terminated_at")
	_tableOrder.Price = field.NewInt32(tableName, "price")

	_tableOrder.fillFieldMap()

	return _tableOrder
}

type tableOrder struct {
	tableOrderDo tableOrderDo

	ALL          field.Asterisk
	TableOrderID field.Int32
	OrderID      field.Int32
	UserID       field.Int32
	Status       field.Int // 订单状态，1｜待支付，2｜支付完成，3｜已退款
	ShopID       field.Int32
	TableID      field.Int32
	CouponID     field.Int32 // 优惠券ID
	UserCouponID field.Int32 // 用户优惠券ID
	Amount       field.Int32 // 订单退过押金之后的金额，在结束订单的时候回写（这里表示的是这个订单实际应该支付的金额，是去除掉优惠券优惠时间的金额）
	PayAmount    field.Int32 // 订单支付金额
	CreatedAt    field.Field // 创建时间
	UpdatedAt    field.Field // 更新时间
	StartedAt    field.Field // 支付时间
	TerminatedAt field.Field // 终止时间
	Price        field.Int32 // 价格

	fieldMap map[string]field.Expr
}

func (t tableOrder) Table(newTableName string) *tableOrder {
	t.tableOrderDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t tableOrder) As(alias string) *tableOrder {
	t.tableOrderDo.DO = *(t.tableOrderDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *tableOrder) updateTableName(table string) *tableOrder {
	t.ALL = field.NewAsterisk(table)
	t.TableOrderID = field.NewInt32(table, "table_order_id")
	t.OrderID = field.NewInt32(table, "order_id")
	t.UserID = field.NewInt32(table, "user_id")
	t.Status = field.NewInt(table, "status")
	t.ShopID = field.NewInt32(table, "shop_id")
	t.TableID = field.NewInt32(table, "table_id")
	t.CouponID = field.NewInt32(table, "coupon_id")
	t.UserCouponID = field.NewInt32(table, "user_coupon_id")
	t.Amount = field.NewInt32(table, "amount")
	t.PayAmount = field.NewInt32(table, "pay_amount")
	t.CreatedAt = field.NewField(table, "created_at")
	t.UpdatedAt = field.NewField(table, "updated_at")
	t.StartedAt = field.NewField(table, "started_at")
	t.TerminatedAt = field.NewField(table, "terminated_at")
	t.Price = field.NewInt32(table, "price")

	t.fillFieldMap()

	return t
}

func (t *tableOrder) WithContext(ctx context.Context) ITableOrderDo {
	return t.tableOrderDo.WithContext(ctx)
}

func (t tableOrder) TableName() string { return t.tableOrderDo.TableName() }

func (t tableOrder) Alias() string { return t.tableOrderDo.Alias() }

func (t *tableOrder) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *tableOrder) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 15)
	t.fieldMap["table_order_id"] = t.TableOrderID
	t.fieldMap["order_id"] = t.OrderID
	t.fieldMap["user_id"] = t.UserID
	t.fieldMap["status"] = t.Status
	t.fieldMap["shop_id"] = t.ShopID
	t.fieldMap["table_id"] = t.TableID
	t.fieldMap["coupon_id"] = t.CouponID
	t.fieldMap["user_coupon_id"] = t.UserCouponID
	t.fieldMap["amount"] = t.Amount
	t.fieldMap["pay_amount"] = t.PayAmount
	t.fieldMap["created_at"] = t.CreatedAt
	t.fieldMap["updated_at"] = t.UpdatedAt
	t.fieldMap["started_at"] = t.StartedAt
	t.fieldMap["terminated_at"] = t.TerminatedAt
	t.fieldMap["price"] = t.Price
}

func (t tableOrder) clone(db *gorm.DB) tableOrder {
	t.tableOrderDo.ReplaceDB(db)
	return t
}

type tableOrderDo struct{ gen.DO }

type ITableOrderDo interface {
	gen.SubQuery
	Debug() ITableOrderDo
	WithContext(ctx context.Context) ITableOrderDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	As(alias string) gen.Dao
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ITableOrderDo
	Not(conds ...gen.Condition) ITableOrderDo
	Or(conds ...gen.Condition) ITableOrderDo
	Select(conds ...field.Expr) ITableOrderDo
	Where(conds ...gen.Condition) ITableOrderDo
	Order(conds ...field.Expr) ITableOrderDo
	Distinct(cols ...field.Expr) ITableOrderDo
	Omit(cols ...field.Expr) ITableOrderDo
	Join(table schema.Tabler, on ...field.Expr) ITableOrderDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ITableOrderDo
	RightJoin(table schema.Tabler, on ...field.Expr) ITableOrderDo
	Group(cols ...field.Expr) ITableOrderDo
	Having(conds ...gen.Condition) ITableOrderDo
	Limit(limit int) ITableOrderDo
	Offset(offset int) ITableOrderDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ITableOrderDo
	Unscoped() ITableOrderDo
	Create(values ...*model.TableOrder) error
	CreateInBatches(values []*model.TableOrder, batchSize int) error
	Save(values ...*model.TableOrder) error
	First() (*model.TableOrder, error)
	Take() (*model.TableOrder, error)
	Last() (*model.TableOrder, error)
	Find() ([]*model.TableOrder, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TableOrder, err error)
	FindInBatches(result *[]*model.TableOrder, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.TableOrder) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ITableOrderDo
	Assign(attrs ...field.AssignExpr) ITableOrderDo
	Joins(fields ...field.RelationField) ITableOrderDo
	Preload(fields ...field.RelationField) ITableOrderDo
	FirstOrInit() (*model.TableOrder, error)
	FirstOrCreate() (*model.TableOrder, error)
	FindByPage(offset int, limit int) (result []*model.TableOrder, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ITableOrderDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (t tableOrderDo) Debug() ITableOrderDo {
	return t.withDO(t.DO.Debug())
}

func (t tableOrderDo) WithContext(ctx context.Context) ITableOrderDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t tableOrderDo) ReadDB() ITableOrderDo {
	return t.Clauses(dbresolver.Read)
}

func (t tableOrderDo) WriteDB() ITableOrderDo {
	return t.Clauses(dbresolver.Write)
}

func (t tableOrderDo) Clauses(conds ...clause.Expression) ITableOrderDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t tableOrderDo) Returning(value interface{}, columns ...string) ITableOrderDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t tableOrderDo) Not(conds ...gen.Condition) ITableOrderDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t tableOrderDo) Or(conds ...gen.Condition) ITableOrderDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t tableOrderDo) Select(conds ...field.Expr) ITableOrderDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t tableOrderDo) Where(conds ...gen.Condition) ITableOrderDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t tableOrderDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) ITableOrderDo {
	return t.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (t tableOrderDo) Order(conds ...field.Expr) ITableOrderDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t tableOrderDo) Distinct(cols ...field.Expr) ITableOrderDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t tableOrderDo) Omit(cols ...field.Expr) ITableOrderDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t tableOrderDo) Join(table schema.Tabler, on ...field.Expr) ITableOrderDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t tableOrderDo) LeftJoin(table schema.Tabler, on ...field.Expr) ITableOrderDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t tableOrderDo) RightJoin(table schema.Tabler, on ...field.Expr) ITableOrderDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t tableOrderDo) Group(cols ...field.Expr) ITableOrderDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t tableOrderDo) Having(conds ...gen.Condition) ITableOrderDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t tableOrderDo) Limit(limit int) ITableOrderDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t tableOrderDo) Offset(offset int) ITableOrderDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t tableOrderDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ITableOrderDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t tableOrderDo) Unscoped() ITableOrderDo {
	return t.withDO(t.DO.Unscoped())
}

func (t tableOrderDo) Create(values ...*model.TableOrder) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t tableOrderDo) CreateInBatches(values []*model.TableOrder, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t tableOrderDo) Save(values ...*model.TableOrder) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t tableOrderDo) First() (*model.TableOrder, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.TableOrder), nil
	}
}

func (t tableOrderDo) Take() (*model.TableOrder, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.TableOrder), nil
	}
}

func (t tableOrderDo) Last() (*model.TableOrder, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.TableOrder), nil
	}
}

func (t tableOrderDo) Find() ([]*model.TableOrder, error) {
	result, err := t.DO.Find()
	return result.([]*model.TableOrder), err
}

func (t tableOrderDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TableOrder, err error) {
	buf := make([]*model.TableOrder, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t tableOrderDo) FindInBatches(result *[]*model.TableOrder, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t tableOrderDo) Attrs(attrs ...field.AssignExpr) ITableOrderDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t tableOrderDo) Assign(attrs ...field.AssignExpr) ITableOrderDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t tableOrderDo) Joins(fields ...field.RelationField) ITableOrderDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t tableOrderDo) Preload(fields ...field.RelationField) ITableOrderDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t tableOrderDo) FirstOrInit() (*model.TableOrder, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.TableOrder), nil
	}
}

func (t tableOrderDo) FirstOrCreate() (*model.TableOrder, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.TableOrder), nil
	}
}

func (t tableOrderDo) FindByPage(offset int, limit int) (result []*model.TableOrder, count int64, err error) {
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

func (t tableOrderDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t tableOrderDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t tableOrderDo) Delete(models ...*model.TableOrder) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *tableOrderDo) withDO(do gen.Dao) *tableOrderDo {
	t.DO = *do.(*gen.DO)
	return t
}
