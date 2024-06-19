package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hilaoyu/gin-extend/engine"
	"github.com/hilaoyu/go-utils/utilCache"
	"github.com/hilaoyu/go-utils/utilEnc"
	"time"
)

const ApiCheckAesApiDataKey = "_api_check_aes_api_data"
const ApiCheckAesEncryptorKey = "_api_check_aes_encryptor"

type GetSecretFunc func(appId string) (secret string, err error)

var dataIdCache = utilCache.NewCache("_api_check_aes_api_data_id_", time.Duration(5)*time.Minute).RegisterStoreMemory(10000)

func ApiCheckAesHandler(getSecret GetSecretFunc, debug ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var appId string
		var apiData string
		response := engine.GetResponse(c)

		if "get" == c.Request.Method {
			appId = c.GetString("app_id")
			apiData = c.GetString("data")
		} else {
			var input = struct {
				AppId string `json:"app_id,omitempty" form:"app_id"`
				Data  string `json:"data,omitempty" form:"data"`
			}{}
			err := c.ShouldBind(&input)
			if err != nil {
				response.Failed(fmt.Sprintf("参数错误: %v", err)).RenderApiJson(c)
				c.Abort()
				return
			}

			appId = input.AppId
			apiData = input.Data
		}

		//fmt.Println(appId, apiData)

		secret, err := getSecret(appId)
		if err != nil {
			response.Failed(fmt.Sprintf("解密错误: %v", err)).RenderApiJson(c)
			c.Abort()
			return
		}

		encryptor := utilEnc.NewAesEncryptor(secret)

		data := struct {
			DataId       string `json:"_data_id"`
			Timestamp    int64  `json:"_timestamp"`
			TimestampOld int64  `json:"timestamp"`
		}{}
		err = encryptor.Decrypt(apiData, &data)
		if nil != err {
			response.Failed(fmt.Sprintf("解密错误: %v", err)).RenderApiJson(c)
			c.Abort()
			return
		}

		now := time.Now().UTC()

		if len(debug) <= 0 || !debug[0] {
			t := data.Timestamp
			if t <= 0 {
				t = data.TimestampOld
			}
			if t <= now.Add(time.Duration(-3)*time.Minute).Unix() || t > now.Add(time.Duration(3)*time.Minute).Unix() {
				response.Failed("解密错误: 数据已过期").RenderApiJson(c)
				c.Abort()
				return
			}
		}
		if len(debug) <= 1 || !debug[1] {
			if "" == data.DataId {
				response.Failed("数据错误").RenderApiJson(c)
				c.Abort()
				return
			}

			if dataIdCache.GetBool(data.DataId) {
				response.Failed("数据已处理过了").RenderApiJson(c)
				c.Abort()
				return
			}
			dataIdCache.Set(data.DataId, true)
		}

		c.Set(ApiCheckAesApiDataKey, apiData)
		c.Set(ApiCheckAesEncryptorKey, encryptor)

		c.Next()
	}
}

func ApiCheckAesGetEnData(gc *gin.Context) string {

	return gc.GetString(ApiCheckAesApiDataKey)
}
func ApiCheckAesGetEncryptor(gc *gin.Context) (encryptor *utilEnc.AesEncryptor, err error) {
	enData := gc.GetString(ApiCheckAesApiDataKey)
	if "" == enData {
		err = fmt.Errorf("密文数据为空")
		return
	}
	encryptorTemp, exist := gc.Get(ApiCheckAesEncryptorKey)
	if !exist {
		err = fmt.Errorf("加密器不存在")
		return
	}
	encryptor, ok := encryptorTemp.(*utilEnc.AesEncryptor)
	if !ok {
		err = fmt.Errorf("加密器类型错误")
	}

	return
}
func ApiCheckAesDecryptData(gc *gin.Context, v interface{}) (err error) {
	enData := ApiCheckAesGetEnData(gc)
	if "" == enData {
		err = fmt.Errorf("密文数据为空")
		return
	}

	encryptor, err := ApiCheckAesGetEncryptor(gc)
	if nil != err {
		return
	}

	err = encryptor.Decrypt(enData, v)
	return
}