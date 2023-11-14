package request

type ExecuteSql struct {
	Sql string `form:"sql" json:"sql"`
}

func (ExecuteSql) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"sql.required": "sql语句不能为空",
	}
}
