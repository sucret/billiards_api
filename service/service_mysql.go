package service

import (
	"errors"
	"fmt"
	"gin-api/pkg/config"
	"gin-api/pkg/mysql"
	"github.com/go-redis/redis"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
)

type mysqlService struct {
	db    *gorm.DB
	redis *redis.Client
}

type tableList struct {
	TableName    string `json:"table_name"`
	TableComment string `json:"table_comment"`
}

var MysqlService = &mysqlService{
	db: mysql.GetDB(),
}

func (m *mysqlService) Tables() (list []tableList, err error) {
	sql := fmt.Sprintf("SELECT `table_name`,`table_comment` FROM `information_schema`.`tables` "+
		"WHERE `table_schema`= '%s'", config.GetConfig().Database.Database)

	m.db.Raw(sql).Scan(&list)
	return
}

func (m *mysqlService) Execute(sql string) (data []map[string]interface{}, err error) {
	prefix := strings.ToTitle(sql[:strings.Index(sql, " ")])
	if prefix != "SELECT" && prefix != "EXPLAIN" && prefix != "SHOW" {
		err = errors.New("仅支持查询操作")
		return
	}

	rows, err := m.db.Raw(sql).Rows()
	if err != nil {
		return
	}

	cols, _ := rows.Columns()

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			fmt.Printf("query table scan error, detail is [%v]\n", err.Error())
			continue
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = cast.ToString(*val)
		}

		data = append(data, m)
	}

	return
}
