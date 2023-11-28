package api

import "github.com/gin-gonic/gin"

type orderApi struct{}

var OrderApi = new(orderApi)

func (*orderApi) List(c *gin.Context) {

}
