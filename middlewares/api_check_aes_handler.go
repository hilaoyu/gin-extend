package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hilaoyu/gin-extend/engine"
	"github.com/hilaoyu/go-utils/utilEnc"
)

type GetAesSecretAndEnDataFunc func(gc *gin.Context) (secret string, enData string, err error)

func ApiCheckAesHandler(getSecret GetAesSecretAndEnDataFunc, debug ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		response := engine.GetResponse(c)

		secret, apiData, err := getSecret(c)
		if err != nil {
			response.Failed(fmt.Sprintf("获取密钥错误: %v", err)).RenderApiJson(c)
			c.Abort()
			return
		}

		encryptor := utilEnc.NewAesEncryptor(secret)

		err = apiCheck(apiData, encryptor, c, debug...)
		if err != nil {
			response.Failed(err.Error()).RenderApiJson(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
