package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hilaoyu/gin-extend/engine"
	"github.com/hilaoyu/go-utils/utilEnc"
)

type GetGmSm2KeyFunc func(gc *gin.Context) (publicKey []byte, privateKey []byte, enData string, err error)

func ApiCheckGmSm2Handler(getRsaKey GetGmSm2KeyFunc, debug ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := engine.GetResponse(c)

		publicKey, privateKey, apiData, err := getRsaKey(c)

		if err != nil {
			response.Failed(fmt.Sprintf("读取解密密钥错误: %v", err)).RenderApiJson(c)
			c.Abort()
			return
		}

		encryptor := utilEnc.NewGmSm2Encryptor()
		_, err = encryptor.SetSm2PublicKey(publicKey)
		if err != nil {
			response.Failed(fmt.Sprintf("设置公钥错误: %v", err)).RenderApiJson(c)
			c.Abort()
			return
		}

		_, err = encryptor.SetSm2PrivateKey(privateKey, nil)
		if err != nil {
			response.Failed(fmt.Sprintf("设置私钥错误: %v", err)).RenderApiJson(c)
			c.Abort()
			return
		}

		err = apiCheck(apiData, encryptor, c, debug...)
		if err != nil {
			response.Failed(err.Error()).RenderApiJson(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
