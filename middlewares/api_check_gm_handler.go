package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hilaoyu/gin-extend/engine"
	"github.com/hilaoyu/go-utils/utilEnc"
)

type GetGmKeyFunc func(appId string, gc *gin.Context) (publicKey []byte, privateKey []byte, err error)

func ApiCheckGmHandler(getRsaKey GetGmKeyFunc, debug ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := engine.GetResponse(c)

		apiData, appId, err := apiCheckGetDataFromGc(c)
		if err != nil {
			response.Failed(err.Error()).RenderApiJson(c)
			c.Abort()
			return
		}
		//fmt.Println(appId, apiData)

		publicKey, privateKey, err := getRsaKey(appId, c)

		if err != nil {
			response.Failed(fmt.Sprintf("读取解密密钥错误: %v", err)).RenderApiJson(c)
			c.Abort()
			return
		}

		encryptor := utilEnc.NewGmEncryptor()
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

		err = apiCheck(apiData, encryptor, utilEnc.ApiDataEncryptorTypeGm, c, debug...)
		if err != nil {
			response.Failed(err.Error()).RenderApiJson(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
