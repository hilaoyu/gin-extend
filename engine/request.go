package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/hilaoyu/go-utils/utilHttp"
)

const (
	ContextVariablesKeyPager = "_gin_extend_context_variables_key_pager"
)

func ReqGetPager(gc *gin.Context, pageSize int) (pager *utilHttp.Paginator) {

	pager = utilHttp.NewPaginator(gc.Request, pageSize, nil)
	gc.Set(ContextVariablesKeyPager, pager)
	return
}
