package request

type RedisQuery struct {
	Method string `form:"method" json:"method"`
	Query  string `form:"query" json:"query"`
}

func (RedisQuery) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"query.required": "查询语句不能为空",
	}
}
