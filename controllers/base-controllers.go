package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hilaoyu/go-utils/utilHttp"
)

type BaseController struct {
}

func (c *BaseController) GetPager(gc *gin.Context, pageSize int) (pager *utilHttp.Paginator) {

	pager = &utilHttp.Paginator{
		Request: gc.Request,
		PerPage: pageSize,
	}
	
	return
}