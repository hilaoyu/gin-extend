package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/hilaoyu/go-utils/utilHttp"
)

func ReqGetPager(gc *gin.Context, pageSize int) (pager *utilHttp.Paginator) {

	pager = &utilHttp.Paginator{
		Request: gc.Request,
		PerPage: pageSize,
	}

	return
}
