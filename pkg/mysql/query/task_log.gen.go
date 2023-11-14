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

	"gin-api/pkg/mysql/model"
)

func newTaskLog(db *gorm.DB) taskLog {
	_taskLog := taskLog{}

	_taskLog.taskLogDo.UseDB(db)
	_taskLog.taskLogDo.UseModel(&model.TaskLog{})

	tableName := _taskLog.taskLogDo.TableName()
	_taskLog.ALL = field.NewAsterisk(tableName)
	_taskLog.TaskLogID = field.NewInt64(tableName, "task_log_id")
	_taskLog.TaskID = field.NewInt32(tableName, "task_id")
	_taskLog.Status = field.NewInt32(tableName, "status")
	_taskLog.StartTime = field.NewField(tableName, "start_time")
	_taskLog.EndTime = field.NewField(tableName, "end_time")
	_taskLog.Log = field.NewString(tableName, "log")
	_taskLog.CreatedAt = field.NewField(tableName, "created_at")
	_taskLog.UpdatedAt = field.NewField(tableName, "updated_at")

	_taskLog.fillFieldMap()

	return _taskLog
}

type taskLog struct {
	taskLogDo taskLogDo

	ALL       field.Asterisk
	TaskLogID field.Int64
	TaskID    field.Int32
	Status    field.Int32 // 任务状态：1|执行中，2|执行成功，3|执行失败，4|手动取消
	StartTime field.Field
	EndTime   field.Field
	Log       field.String
	CreatedAt field.Field
	UpdatedAt field.Field

	fieldMap map[string]field.Expr
}

func (t taskLog) Table(newTableName string) *taskLog {
	t.taskLogDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t taskLog) As(alias string) *taskLog {
	t.taskLogDo.DO = *(t.taskLogDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *taskLog) updateTableName(table string) *taskLog {
	t.ALL = field.NewAsterisk(table)
	t.TaskLogID = field.NewInt64(table, "task_log_id")
	t.TaskID = field.NewInt32(table, "task_id")
	t.Status = field.NewInt32(table, "status")
	t.StartTime = field.NewField(table, "start_time")
	t.EndTime = field.NewField(table, "end_time")
	t.Log = field.NewString(table, "log")
	t.CreatedAt = field.NewField(table, "created_at")
	t.UpdatedAt = field.NewField(table, "updated_at")

	t.fillFieldMap()

	return t
}

func (t *taskLog) WithContext(ctx context.Context) ITaskLogDo { return t.taskLogDo.WithContext(ctx) }

func (t taskLog) TableName() string { return t.taskLogDo.TableName() }

func (t taskLog) Alias() string { return t.taskLogDo.Alias() }

func (t *taskLog) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *taskLog) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 8)
	t.fieldMap["task_log_id"] = t.TaskLogID
	t.fieldMap["task_id"] = t.TaskID
	t.fieldMap["status"] = t.Status
	t.fieldMap["start_time"] = t.StartTime
	t.fieldMap["end_time"] = t.EndTime
	t.fieldMap["log"] = t.Log
	t.fieldMap["created_at"] = t.CreatedAt
	t.fieldMap["updated_at"] = t.UpdatedAt
}

func (t taskLog) clone(db *gorm.DB) taskLog {
	t.taskLogDo.ReplaceDB(db)
	return t
}

type taskLogDo struct{ gen.DO }

type ITaskLogDo interface {
	gen.SubQuery
	Debug() ITaskLogDo
	WithContext(ctx context.Context) ITaskLogDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	As(alias string) gen.Dao
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ITaskLogDo
	Not(conds ...gen.Condition) ITaskLogDo
	Or(conds ...gen.Condition) ITaskLogDo
	Select(conds ...field.Expr) ITaskLogDo
	Where(conds ...gen.Condition) ITaskLogDo
	Order(conds ...field.Expr) ITaskLogDo
	Distinct(cols ...field.Expr) ITaskLogDo
	Omit(cols ...field.Expr) ITaskLogDo
	Join(table schema.Tabler, on ...field.Expr) ITaskLogDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ITaskLogDo
	RightJoin(table schema.Tabler, on ...field.Expr) ITaskLogDo
	Group(cols ...field.Expr) ITaskLogDo
	Having(conds ...gen.Condition) ITaskLogDo
	Limit(limit int) ITaskLogDo
	Offset(offset int) ITaskLogDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ITaskLogDo
	Unscoped() ITaskLogDo
	Create(values ...*model.TaskLog) error
	CreateInBatches(values []*model.TaskLog, batchSize int) error
	Save(values ...*model.TaskLog) error
	First() (*model.TaskLog, error)
	Take() (*model.TaskLog, error)
	Last() (*model.TaskLog, error)
	Find() ([]*model.TaskLog, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TaskLog, err error)
	FindInBatches(result *[]*model.TaskLog, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.TaskLog) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ITaskLogDo
	Assign(attrs ...field.AssignExpr) ITaskLogDo
	Joins(fields ...field.RelationField) ITaskLogDo
	Preload(fields ...field.RelationField) ITaskLogDo
	FirstOrInit() (*model.TaskLog, error)
	FirstOrCreate() (*model.TaskLog, error)
	FindByPage(offset int, limit int) (result []*model.TaskLog, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ITaskLogDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (t taskLogDo) Debug() ITaskLogDo {
	return t.withDO(t.DO.Debug())
}

func (t taskLogDo) WithContext(ctx context.Context) ITaskLogDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t taskLogDo) ReadDB() ITaskLogDo {
	return t.Clauses(dbresolver.Read)
}

func (t taskLogDo) WriteDB() ITaskLogDo {
	return t.Clauses(dbresolver.Write)
}

func (t taskLogDo) Clauses(conds ...clause.Expression) ITaskLogDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t taskLogDo) Returning(value interface{}, columns ...string) ITaskLogDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t taskLogDo) Not(conds ...gen.Condition) ITaskLogDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t taskLogDo) Or(conds ...gen.Condition) ITaskLogDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t taskLogDo) Select(conds ...field.Expr) ITaskLogDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t taskLogDo) Where(conds ...gen.Condition) ITaskLogDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t taskLogDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) ITaskLogDo {
	return t.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (t taskLogDo) Order(conds ...field.Expr) ITaskLogDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t taskLogDo) Distinct(cols ...field.Expr) ITaskLogDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t taskLogDo) Omit(cols ...field.Expr) ITaskLogDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t taskLogDo) Join(table schema.Tabler, on ...field.Expr) ITaskLogDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t taskLogDo) LeftJoin(table schema.Tabler, on ...field.Expr) ITaskLogDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t taskLogDo) RightJoin(table schema.Tabler, on ...field.Expr) ITaskLogDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t taskLogDo) Group(cols ...field.Expr) ITaskLogDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t taskLogDo) Having(conds ...gen.Condition) ITaskLogDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t taskLogDo) Limit(limit int) ITaskLogDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t taskLogDo) Offset(offset int) ITaskLogDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t taskLogDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ITaskLogDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t taskLogDo) Unscoped() ITaskLogDo {
	return t.withDO(t.DO.Unscoped())
}

func (t taskLogDo) Create(values ...*model.TaskLog) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t taskLogDo) CreateInBatches(values []*model.TaskLog, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t taskLogDo) Save(values ...*model.TaskLog) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t taskLogDo) First() (*model.TaskLog, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.TaskLog), nil
	}
}

func (t taskLogDo) Take() (*model.TaskLog, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.TaskLog), nil
	}
}

func (t taskLogDo) Last() (*model.TaskLog, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.TaskLog), nil
	}
}

func (t taskLogDo) Find() ([]*model.TaskLog, error) {
	result, err := t.DO.Find()
	return result.([]*model.TaskLog), err
}

func (t taskLogDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TaskLog, err error) {
	buf := make([]*model.TaskLog, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t taskLogDo) FindInBatches(result *[]*model.TaskLog, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t taskLogDo) Attrs(attrs ...field.AssignExpr) ITaskLogDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t taskLogDo) Assign(attrs ...field.AssignExpr) ITaskLogDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t taskLogDo) Joins(fields ...field.RelationField) ITaskLogDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t taskLogDo) Preload(fields ...field.RelationField) ITaskLogDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t taskLogDo) FirstOrInit() (*model.TaskLog, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.TaskLog), nil
	}
}

func (t taskLogDo) FirstOrCreate() (*model.TaskLog, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.TaskLog), nil
	}
}

func (t taskLogDo) FindByPage(offset int, limit int) (result []*model.TaskLog, count int64, err error) {
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

func (t taskLogDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t taskLogDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t taskLogDo) Delete(models ...*model.TaskLog) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *taskLogDo) withDO(do gen.Dao) *taskLogDo {
	t.DO = *do.(*gen.DO)
	return t
}
