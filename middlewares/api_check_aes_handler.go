package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hilaoyu/gin-extend/engine"
	"github.com/hilaoyu/go-utils/utilEnc"
)

type GetSecretFunc func(appId string, gc *gin.Context) (secret string, err error)

func ApiCheckAesHandler(getSecret GetSecretFunc, debug ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		response := engine.GetResponse(c)

		apiData, appId, err := apiCheckGetDataFromGc(c)
		if err != nil {
			response.Failed(err.Error()).RenderApiJson(c)
			c.Abort()
			return
		}

		//fmt.Println(appId, apiData)

		secret, err := getSecret(appId, c)
		if err != nil {
			response.Failed(fmt.Sprintf("获取密钥错误: %v", err)).RenderApiJson(c)
			c.Abort()
			return
		}

		encryptor := utilEnc.NewAesEncryptor(secret)

		err = apiCheck(apiData, encryptor, utilEnc.ApiDataEncryptorTypeAes, c, debug...)
		if err != nil {
			response.Failed(err.Error()).RenderApiJson(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
