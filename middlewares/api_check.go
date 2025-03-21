package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hilaoyu/go-utils/utilCache"
	"github.com/hilaoyu/go-utils/utilEnc"
	"time"
)

const ApiCheckApiDataKey = "_api_check_api_data"
const ApiCheckEncryptorKey = "_api_check_encryptor"
const ApiCheckEncryptorTypeKey = "_api_check_encryptor_type"

var (
	dataIdCacheApiCheck = utilCache.NewCache("_api_check_api_data_id_", time.Duration(5)*time.Minute).RegisterStoreMemory(10000)
)

func apiCheckGetDataFromGc(gc *gin.Context) (apiData string, appId string, err error) {
	if "get" == gc.Request.Method {
		appId = gc.DefaultQuery("app_id", "")
		apiData = gc.DefaultQuery("data", "")
	} else {
		var input = struct {
			AppId string `json:"app_id,omitempty" form:"app_id"`
			Data  string `json:"data,omitempty" form:"data"`
		}{}
		err = gc.ShouldBind(&input)
		if err != nil {
			err = fmt.Errorf("参数错误: %v", err)
			return
		}

		appId = input.AppId
		apiData = input.Data
	}
	return
}
func apiCheck(apiData string, encryptor utilEnc.ApiDataEncryptor, encryptorType string, gc *gin.Context, debug ...bool) (err error) {
	data := struct {
		DataId       string `json:"_data_id"`
		Timestamp    int64  `json:"_timestamp"`
		TimestampOld int64  `json:"timestamp"`
	}{}
	err = encryptor.ApiDataDecrypt(apiData, &data)
	if nil != err {
		err = fmt.Errorf("解密错误: %v", err)
		return
	}

	now := time.Now().UTC()

	if len(debug) <= 0 || !debug[0] {
		t := data.Timestamp
		if t <= 0 {
			t = data.TimestampOld
		}
		if t <= now.Add(time.Duration(-3)*time.Minute).Unix() || t > now.Add(time.Duration(3)*time.Minute).Unix() {
			err = fmt.Errorf("解密错误: 数据已过期")
			return
		}
	}
	if len(debug) <= 1 || !debug[1] {
		if "" == data.DataId {
			err = fmt.Errorf("数据Id错误")
			return
		}

		if dataIdCacheApiCheck.GetBool(data.DataId) {
			err = fmt.Errorf("数据已处理过了")
			return
		}
		_ = dataIdCacheApiCheck.Set(data.DataId, true)
	}

	apiCheckSetEnData(gc, apiData)
	apiCheckSetEncryptor(gc, encryptor, encryptorType)
	return
}

func apiCheckSetEnData(gc *gin.Context, apiData string) {
	gc.Set(ApiCheckApiDataKey, apiData)
}
func apiCheckSetEncryptor(gc *gin.Context, encryptor utilEnc.ApiDataEncryptor, encryptorType string) {
	gc.Set(ApiCheckEncryptorKey, encryptor)
	gc.Set(ApiCheckEncryptorTypeKey, encryptorType)

}

func ApiCheckGetEnData(gc *gin.Context) string {
	return gc.GetString(ApiCheckApiDataKey)
}
func ApiCheckGetEncryptor(gc *gin.Context) (encryptor utilEnc.ApiDataEncryptor, err error) {
	encryptorType := gc.GetString(ApiCheckEncryptorTypeKey)

	encryptorTemp, exist := gc.Get(ApiCheckEncryptorKey)
	if !exist {
		err = fmt.Errorf("加密器不存在")
		return
	}
	var ok bool
	switch encryptorType {
	case utilEnc.ApiDataEncryptorTypeAes:
		encryptor, ok = encryptorTemp.(*utilEnc.AesEncryptor)
		break
	case utilEnc.ApiDataEncryptorTypeRsa:
		encryptor, ok = encryptorTemp.(*utilEnc.RsaEncryptor)
		break
	}

	if !ok {
		err = fmt.Errorf("加密器类型错误")
	}

	return
}

func ApiCheckGetAesEncryptor(gc *gin.Context) (encryptor *utilEnc.AesEncryptor, err error) {

	encryptorTemp, exist := gc.Get(ApiCheckEncryptorKey)
	if !exist {
		err = fmt.Errorf("加密器不存在")
		return
	}
	encryptorType := gc.GetString(ApiCheckEncryptorTypeKey)
	if utilEnc.ApiDataEncryptorTypeAes != encryptorType {
		err = fmt.Errorf("加密器类型不是AES")
		return
	}

	encryptor, ok := encryptorTemp.(*utilEnc.AesEncryptor)
	if !ok {
		err = fmt.Errorf("加密器类型错误")
	}

	return
}
func ApiCheckGetRsaEncryptor(gc *gin.Context) (encryptor *utilEnc.RsaEncryptor, err error) {

	encryptorTemp, exist := gc.Get(ApiCheckEncryptorKey)
	if !exist {
		err = fmt.Errorf("加密器不存在")
		return
	}
	encryptorType := gc.GetString(ApiCheckEncryptorTypeKey)
	if utilEnc.ApiDataEncryptorTypeRsa != encryptorType {
		err = fmt.Errorf("加密器类型不是RSA")
		return
	}

	encryptor, ok := encryptorTemp.(*utilEnc.RsaEncryptor)
	if !ok {
		err = fmt.Errorf("加密器类型错误")
	}

	return
}
func ApiCheckDecryptData(gc *gin.Context, v interface{}) (err error) {
	enData := ApiCheckGetEnData(gc)
	if "" == enData {
		err = fmt.Errorf("密文数据为空")
		return
	}

	encryptor, err := ApiCheckGetEncryptor(gc)
	if nil != err {
		return
	}

	err = encryptor.ApiDataDecrypt(enData, v)
	return
}
func ApiCheckEncryptData(gc *gin.Context, data interface{}) (enData string, err error) {
	encryptor, err := ApiCheckGetEncryptor(gc)
	if nil != err {
		return
	}
	enData, err = encryptor.ApiDataEncrypt(data)
	return
}
