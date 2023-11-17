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

func newShop(db *gorm.DB) shop {
	_shop := shop{}

	_shop.shopDo.UseDB(db)
	_shop.shopDo.UseModel(&model.Shop{})

	tableName := _shop.shopDo.TableName()
	_shop.ALL = field.NewAsterisk(tableName)
	_shop.ShopID = field.NewInt32(tableName, "shop_id")
	_shop.Name = field.NewString(tableName, "name")
	_shop.Status = field.NewInt32(tableName, "status")
	_shop.Address = field.NewString(tableName, "address")
	_shop.CreatedAt = field.NewField(tableName, "created_at")

	_shop.fillFieldMap()

	return _shop
}

type shop struct {
	shopDo shopDo

	ALL       field.Asterisk
	ShopID    field.Int32
	Name      field.String // 门店名称
	Status    field.Int32  // 门店状态，1｜开启，2｜关闭
	Address   field.String // 店铺地址
	CreatedAt field.Field

	fieldMap map[string]field.Expr
}

func (s shop) Table(newTableName string) *shop {
	s.shopDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s shop) As(alias string) *shop {
	s.shopDo.DO = *(s.shopDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *shop) updateTableName(table string) *shop {
	s.ALL = field.NewAsterisk(table)
	s.ShopID = field.NewInt32(table, "shop_id")
	s.Name = field.NewString(table, "name")
	s.Status = field.NewInt32(table, "status")
	s.Address = field.NewString(table, "address")
	s.CreatedAt = field.NewField(table, "created_at")

	s.fillFieldMap()

	return s
}

func (s *shop) WithContext(ctx context.Context) IShopDo { return s.shopDo.WithContext(ctx) }

func (s shop) TableName() string { return s.shopDo.TableName() }

func (s shop) Alias() string { return s.shopDo.Alias() }

func (s *shop) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *shop) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 5)
	s.fieldMap["shop_id"] = s.ShopID
	s.fieldMap["name"] = s.Name
	s.fieldMap["status"] = s.Status
	s.fieldMap["address"] = s.Address
	s.fieldMap["created_at"] = s.CreatedAt
}

func (s shop) clone(db *gorm.DB) shop {
	s.shopDo.ReplaceDB(db)
	return s
}

type shopDo struct{ gen.DO }

type IShopDo interface {
	gen.SubQuery
	Debug() IShopDo
	WithContext(ctx context.Context) IShopDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	As(alias string) gen.Dao
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IShopDo
	Not(conds ...gen.Condition) IShopDo
	Or(conds ...gen.Condition) IShopDo
	Select(conds ...field.Expr) IShopDo
	Where(conds ...gen.Condition) IShopDo
	Order(conds ...field.Expr) IShopDo
	Distinct(cols ...field.Expr) IShopDo
	Omit(cols ...field.Expr) IShopDo
	Join(table schema.Tabler, on ...field.Expr) IShopDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IShopDo
	RightJoin(table schema.Tabler, on ...field.Expr) IShopDo
	Group(cols ...field.Expr) IShopDo
	Having(conds ...gen.Condition) IShopDo
	Limit(limit int) IShopDo
	Offset(offset int) IShopDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IShopDo
	Unscoped() IShopDo
	Create(values ...*model.Shop) error
	CreateInBatches(values []*model.Shop, batchSize int) error
	Save(values ...*model.Shop) error
	First() (*model.Shop, error)
	Take() (*model.Shop, error)
	Last() (*model.Shop, error)
	Find() ([]*model.Shop, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Shop, err error)
	FindInBatches(result *[]*model.Shop, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Shop) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IShopDo
	Assign(attrs ...field.AssignExpr) IShopDo
	Joins(fields ...field.RelationField) IShopDo
	Preload(fields ...field.RelationField) IShopDo
	FirstOrInit() (*model.Shop, error)
	FirstOrCreate() (*model.Shop, error)
	FindByPage(offset int, limit int) (result []*model.Shop, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IShopDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (s shopDo) Debug() IShopDo {
	return s.withDO(s.DO.Debug())
}

func (s shopDo) WithContext(ctx context.Context) IShopDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s shopDo) ReadDB() IShopDo {
	return s.Clauses(dbresolver.Read)
}

func (s shopDo) WriteDB() IShopDo {
	return s.Clauses(dbresolver.Write)
}

func (s shopDo) Clauses(conds ...clause.Expression) IShopDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s shopDo) Returning(value interface{}, columns ...string) IShopDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s shopDo) Not(conds ...gen.Condition) IShopDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s shopDo) Or(conds ...gen.Condition) IShopDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s shopDo) Select(conds ...field.Expr) IShopDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s shopDo) Where(conds ...gen.Condition) IShopDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s shopDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IShopDo {
	return s.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (s shopDo) Order(conds ...field.Expr) IShopDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s shopDo) Distinct(cols ...field.Expr) IShopDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s shopDo) Omit(cols ...field.Expr) IShopDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s shopDo) Join(table schema.Tabler, on ...field.Expr) IShopDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s shopDo) LeftJoin(table schema.Tabler, on ...field.Expr) IShopDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s shopDo) RightJoin(table schema.Tabler, on ...field.Expr) IShopDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s shopDo) Group(cols ...field.Expr) IShopDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s shopDo) Having(conds ...gen.Condition) IShopDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s shopDo) Limit(limit int) IShopDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s shopDo) Offset(offset int) IShopDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s shopDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IShopDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s shopDo) Unscoped() IShopDo {
	return s.withDO(s.DO.Unscoped())
}

func (s shopDo) Create(values ...*model.Shop) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s shopDo) CreateInBatches(values []*model.Shop, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s shopDo) Save(values ...*model.Shop) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s shopDo) First() (*model.Shop, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Shop), nil
	}
}

func (s shopDo) Take() (*model.Shop, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Shop), nil
	}
}

func (s shopDo) Last() (*model.Shop, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Shop), nil
	}
}

func (s shopDo) Find() ([]*model.Shop, error) {
	result, err := s.DO.Find()
	return result.([]*model.Shop), err
}

func (s shopDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Shop, err error) {
	buf := make([]*model.Shop, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s shopDo) FindInBatches(result *[]*model.Shop, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s shopDo) Attrs(attrs ...field.AssignExpr) IShopDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s shopDo) Assign(attrs ...field.AssignExpr) IShopDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s shopDo) Joins(fields ...field.RelationField) IShopDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s shopDo) Preload(fields ...field.RelationField) IShopDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s shopDo) FirstOrInit() (*model.Shop, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Shop), nil
	}
}

func (s shopDo) FirstOrCreate() (*model.Shop, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Shop), nil
	}
}

func (s shopDo) FindByPage(offset int, limit int) (result []*model.Shop, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s shopDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s shopDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s shopDo) Delete(models ...*model.Shop) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *shopDo) withDO(do gen.Dao) *shopDo {
	s.DO = *do.(*gen.DO)
	return s
}
