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
func apiCheck(apiData string, encryptor utilEnc.ApiDataEncryptor, gc *gin.Context, debug ...bool) (err error) {
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
	apiCheckSetEncryptor(gc, encryptor)
	return
}

func apiCheckSetEnData(gc *gin.Context, apiData string) {
	gc.Set(ApiCheckApiDataKey, apiData)
}
func apiCheckSetEncryptor(gc *gin.Context, encryptor utilEnc.ApiDataEncryptor) {
	gc.Set(ApiCheckEncryptorKey, encryptor)

}

func ApiCheckGetEnData(gc *gin.Context) string {
	return gc.GetString(ApiCheckApiDataKey)
}
func ApiCheckGetEncryptor(gc *gin.Context) (encryptor utilEnc.ApiDataEncryptor, err error) {
	encryptorTemp, exist := gc.Get(ApiCheckEncryptorKey)
	if !exist {
		err = fmt.Errorf("加密器不存在")
		return
	}
	var ok bool
	encryptor, ok = encryptorTemp.(utilEnc.ApiDataEncryptor)
	if !ok {
		err = fmt.Errorf("加密器错误")
	}

	return
}

func ApiCheckGetAesEncryptor(gc *gin.Context) (encryptor *utilEnc.AesEncryptor, err error) {

	encryptorTemp, err := ApiCheckGetEncryptor(gc)
	if nil != err {
		return
	}
	if utilEnc.ApiDataEncryptorTypeAes != encryptorTemp.EncryptorType() {
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

	encryptorTemp, err := ApiCheckGetEncryptor(gc)
	if nil != err {
		return
	}
	if utilEnc.ApiDataEncryptorTypeRsa != encryptorTemp.EncryptorType() {
		err = fmt.Errorf("加密器类型不是RSA")
		return
	}

	encryptor, ok := encryptorTemp.(*utilEnc.RsaEncryptor)
	if !ok {
		err = fmt.Errorf("加密器类型错误")
	}

	return
}
func ApiCheckGetGmSm2Encryptor(gc *gin.Context) (encryptor *utilEnc.GmSm2Encryptor, err error) {

	encryptorTemp, err := ApiCheckGetEncryptor(gc)
	if nil != err {
		return
	}
	if utilEnc.ApiDataEncryptorTypeGmSm2 != encryptorTemp.EncryptorType() {
		err = fmt.Errorf("加密器类型不是GM_SM2")
		return
	}

	encryptor, ok := encryptorTemp.(*utilEnc.GmSm2Encryptor)
	if !ok {
		err = fmt.Errorf("加密器类型错误")
	}

	return
}
func ApiCheckGetGmSm4Encryptor(gc *gin.Context) (encryptor *utilEnc.GmSm4Encryptor, err error) {

	encryptorTemp, err := ApiCheckGetEncryptor(gc)
	if nil != err {
		return
	}
	if utilEnc.ApiDataEncryptorTypeGmSm4 != encryptorTemp.EncryptorType() {
		err = fmt.Errorf("加密器类型不是GM_SM4")
		return
	}

	encryptor, ok := encryptorTemp.(*utilEnc.GmSm4Encryptor)
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
