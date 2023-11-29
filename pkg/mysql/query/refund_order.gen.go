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

func newRefundOrder(db *gorm.DB) refundOrder {
	_refundOrder := refundOrder{}

	_refundOrder.refundOrderDo.UseDB(db)
	_refundOrder.refundOrderDo.UseModel(&model.RefundOrder{})

	tableName := _refundOrder.refundOrderDo.TableName()
	_refundOrder.ALL = field.NewAsterisk(tableName)
	_refundOrder.RefundOrderID = field.NewInt32(tableName, "refund_order_id")
	_refundOrder.PaymentOrderID = field.NewInt32(tableName, "payment_order_id")
	_refundOrder.OrderID = field.NewInt32(tableName, "order_id")
	_refundOrder.OrderNum = field.NewString(tableName, "order_num")
	_refundOrder.Status = field.NewInt32(tableName, "status")
	_refundOrder.RefundNum = field.NewString(tableName, "refund_num")
	_refundOrder.Amount = field.NewInt32(tableName, "amount")
	_refundOrder.WxRefundID = field.NewString(tableName, "wx_refund_id")
	_refundOrder.CreatedAt = field.NewField(tableName, "created_at")
	_refundOrder.UpdatedAt = field.NewField(tableName, "updated_at")

	_refundOrder.fillFieldMap()

	return _refundOrder
}

type refundOrder struct {
	refundOrderDo refundOrderDo

	ALL            field.Asterisk
	RefundOrderID  field.Int32
	PaymentOrderID field.Int32 // 付款单id
	OrderID        field.Int32
	OrderNum       field.String
	Status         field.Int32  // 退款状态，1｜待退款，2｜退款成功，3｜退款关闭，4｜退款处理中，5｜退款异常，微信对应枚举值：SUCCESS|退款成功，CLOSED｜退款关闭，PROCESSING｜退款处理中，ABNORMAL｜退款异常
	RefundNum      field.String // 退款单号，不可重复，同意单号多次请求只会退一次
	Amount         field.Int32  // 退款金额
	WxRefundID     field.String // 微信支付退款单号
	CreatedAt      field.Field
	UpdatedAt      field.Field

	fieldMap map[string]field.Expr
}

func (r refundOrder) Table(newTableName string) *refundOrder {
	r.refundOrderDo.UseTable(newTableName)
	return r.updateTableName(newTableName)
}

func (r refundOrder) As(alias string) *refundOrder {
	r.refundOrderDo.DO = *(r.refundOrderDo.As(alias).(*gen.DO))
	return r.updateTableName(alias)
}

func (r *refundOrder) updateTableName(table string) *refundOrder {
	r.ALL = field.NewAsterisk(table)
	r.RefundOrderID = field.NewInt32(table, "refund_order_id")
	r.PaymentOrderID = field.NewInt32(table, "payment_order_id")
	r.OrderID = field.NewInt32(table, "order_id")
	r.OrderNum = field.NewString(table, "order_num")
	r.Status = field.NewInt32(table, "status")
	r.RefundNum = field.NewString(table, "refund_num")
	r.Amount = field.NewInt32(table, "amount")
	r.WxRefundID = field.NewString(table, "wx_refund_id")
	r.CreatedAt = field.NewField(table, "created_at")
	r.UpdatedAt = field.NewField(table, "updated_at")

	r.fillFieldMap()

	return r
}

func (r *refundOrder) WithContext(ctx context.Context) IRefundOrderDo {
	return r.refundOrderDo.WithContext(ctx)
}

func (r refundOrder) TableName() string { return r.refundOrderDo.TableName() }

func (r refundOrder) Alias() string { return r.refundOrderDo.Alias() }

func (r *refundOrder) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := r.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (r *refundOrder) fillFieldMap() {
	r.fieldMap = make(map[string]field.Expr, 10)
	r.fieldMap["refund_order_id"] = r.RefundOrderID
	r.fieldMap["payment_order_id"] = r.PaymentOrderID
	r.fieldMap["order_id"] = r.OrderID
	r.fieldMap["order_num"] = r.OrderNum
	r.fieldMap["status"] = r.Status
	r.fieldMap["refund_num"] = r.RefundNum
	r.fieldMap["amount"] = r.Amount
	r.fieldMap["wx_refund_id"] = r.WxRefundID
	r.fieldMap["created_at"] = r.CreatedAt
	r.fieldMap["updated_at"] = r.UpdatedAt
}

func (r refundOrder) clone(db *gorm.DB) refundOrder {
	r.refundOrderDo.ReplaceDB(db)
	return r
}

type refundOrderDo struct{ gen.DO }

type IRefundOrderDo interface {
	gen.SubQuery
	Debug() IRefundOrderDo
	WithContext(ctx context.Context) IRefundOrderDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	As(alias string) gen.Dao
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IRefundOrderDo
	Not(conds ...gen.Condition) IRefundOrderDo
	Or(conds ...gen.Condition) IRefundOrderDo
	Select(conds ...field.Expr) IRefundOrderDo
	Where(conds ...gen.Condition) IRefundOrderDo
	Order(conds ...field.Expr) IRefundOrderDo
	Distinct(cols ...field.Expr) IRefundOrderDo
	Omit(cols ...field.Expr) IRefundOrderDo
	Join(table schema.Tabler, on ...field.Expr) IRefundOrderDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IRefundOrderDo
	RightJoin(table schema.Tabler, on ...field.Expr) IRefundOrderDo
	Group(cols ...field.Expr) IRefundOrderDo
	Having(conds ...gen.Condition) IRefundOrderDo
	Limit(limit int) IRefundOrderDo
	Offset(offset int) IRefundOrderDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IRefundOrderDo
	Unscoped() IRefundOrderDo
	Create(values ...*model.RefundOrder) error
	CreateInBatches(values []*model.RefundOrder, batchSize int) error
	Save(values ...*model.RefundOrder) error
	First() (*model.RefundOrder, error)
	Take() (*model.RefundOrder, error)
	Last() (*model.RefundOrder, error)
	Find() ([]*model.RefundOrder, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.RefundOrder, err error)
	FindInBatches(result *[]*model.RefundOrder, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.RefundOrder) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IRefundOrderDo
	Assign(attrs ...field.AssignExpr) IRefundOrderDo
	Joins(fields ...field.RelationField) IRefundOrderDo
	Preload(fields ...field.RelationField) IRefundOrderDo
	FirstOrInit() (*model.RefundOrder, error)
	FirstOrCreate() (*model.RefundOrder, error)
	FindByPage(offset int, limit int) (result []*model.RefundOrder, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IRefundOrderDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (r refundOrderDo) Debug() IRefundOrderDo {
	return r.withDO(r.DO.Debug())
}

func (r refundOrderDo) WithContext(ctx context.Context) IRefundOrderDo {
	return r.withDO(r.DO.WithContext(ctx))
}

func (r refundOrderDo) ReadDB() IRefundOrderDo {
	return r.Clauses(dbresolver.Read)
}

func (r refundOrderDo) WriteDB() IRefundOrderDo {
	return r.Clauses(dbresolver.Write)
}

func (r refundOrderDo) Clauses(conds ...clause.Expression) IRefundOrderDo {
	return r.withDO(r.DO.Clauses(conds...))
}

func (r refundOrderDo) Returning(value interface{}, columns ...string) IRefundOrderDo {
	return r.withDO(r.DO.Returning(value, columns...))
}

func (r refundOrderDo) Not(conds ...gen.Condition) IRefundOrderDo {
	return r.withDO(r.DO.Not(conds...))
}

func (r refundOrderDo) Or(conds ...gen.Condition) IRefundOrderDo {
	return r.withDO(r.DO.Or(conds...))
}

func (r refundOrderDo) Select(conds ...field.Expr) IRefundOrderDo {
	return r.withDO(r.DO.Select(conds...))
}

func (r refundOrderDo) Where(conds ...gen.Condition) IRefundOrderDo {
	return r.withDO(r.DO.Where(conds...))
}

func (r refundOrderDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IRefundOrderDo {
	return r.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (r refundOrderDo) Order(conds ...field.Expr) IRefundOrderDo {
	return r.withDO(r.DO.Order(conds...))
}

func (r refundOrderDo) Distinct(cols ...field.Expr) IRefundOrderDo {
	return r.withDO(r.DO.Distinct(cols...))
}

func (r refundOrderDo) Omit(cols ...field.Expr) IRefundOrderDo {
	return r.withDO(r.DO.Omit(cols...))
}

func (r refundOrderDo) Join(table schema.Tabler, on ...field.Expr) IRefundOrderDo {
	return r.withDO(r.DO.Join(table, on...))
}

func (r refundOrderDo) LeftJoin(table schema.Tabler, on ...field.Expr) IRefundOrderDo {
	return r.withDO(r.DO.LeftJoin(table, on...))
}

func (r refundOrderDo) RightJoin(table schema.Tabler, on ...field.Expr) IRefundOrderDo {
	return r.withDO(r.DO.RightJoin(table, on...))
}

func (r refundOrderDo) Group(cols ...field.Expr) IRefundOrderDo {
	return r.withDO(r.DO.Group(cols...))
}

func (r refundOrderDo) Having(conds ...gen.Condition) IRefundOrderDo {
	return r.withDO(r.DO.Having(conds...))
}

func (r refundOrderDo) Limit(limit int) IRefundOrderDo {
	return r.withDO(r.DO.Limit(limit))
}

func (r refundOrderDo) Offset(offset int) IRefundOrderDo {
	return r.withDO(r.DO.Offset(offset))
}

func (r refundOrderDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IRefundOrderDo {
	return r.withDO(r.DO.Scopes(funcs...))
}

func (r refundOrderDo) Unscoped() IRefundOrderDo {
	return r.withDO(r.DO.Unscoped())
}

func (r refundOrderDo) Create(values ...*model.RefundOrder) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Create(values)
}

func (r refundOrderDo) CreateInBatches(values []*model.RefundOrder, batchSize int) error {
	return r.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (r refundOrderDo) Save(values ...*model.RefundOrder) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Save(values)
}

func (r refundOrderDo) First() (*model.RefundOrder, error) {
	if result, err := r.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.RefundOrder), nil
	}
}

func (r refundOrderDo) Take() (*model.RefundOrder, error) {
	if result, err := r.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.RefundOrder), nil
	}
}

func (r refundOrderDo) Last() (*model.RefundOrder, error) {
	if result, err := r.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.RefundOrder), nil
	}
}

func (r refundOrderDo) Find() ([]*model.RefundOrder, error) {
	result, err := r.DO.Find()
	return result.([]*model.RefundOrder), err
}

func (r refundOrderDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.RefundOrder, err error) {
	buf := make([]*model.RefundOrder, 0, batchSize)
	err = r.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (r refundOrderDo) FindInBatches(result *[]*model.RefundOrder, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return r.DO.FindInBatches(result, batchSize, fc)
}

func (r refundOrderDo) Attrs(attrs ...field.AssignExpr) IRefundOrderDo {
	return r.withDO(r.DO.Attrs(attrs...))
}

func (r refundOrderDo) Assign(attrs ...field.AssignExpr) IRefundOrderDo {
	return r.withDO(r.DO.Assign(attrs...))
}

func (r refundOrderDo) Joins(fields ...field.RelationField) IRefundOrderDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Joins(_f))
	}
	return &r
}

func (r refundOrderDo) Preload(fields ...field.RelationField) IRefundOrderDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Preload(_f))
	}
	return &r
}

func (r refundOrderDo) FirstOrInit() (*model.RefundOrder, error) {
	if result, err := r.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.RefundOrder), nil
	}
}

func (r refundOrderDo) FirstOrCreate() (*model.RefundOrder, error) {
	if result, err := r.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.RefundOrder), nil
	}
}

func (r refundOrderDo) FindByPage(offset int, limit int) (result []*model.RefundOrder, count int64, err error) {
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

func (r refundOrderDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = r.Count()
	if err != nil {
		return
	}

	err = r.Offset(offset).Limit(limit).Scan(result)
	return
}

func (r refundOrderDo) Scan(result interface{}) (err error) {
	return r.DO.Scan(result)
}

func (r refundOrderDo) Delete(models ...*model.RefundOrder) (result gen.ResultInfo, err error) {
	return r.DO.Delete(models)
}

func (r *refundOrderDo) withDO(do gen.Dao) *refundOrderDo {
	r.DO = *do.(*gen.DO)
	return r
}