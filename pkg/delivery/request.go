package delivery

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func GetRawData(ctx *gin.Context) ([]byte, error) {

	body, err := ctx.GetRawData()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetBodyFromRawData(ctx *gin.Context, out interface{}) error {

	body, err := GetRawData(ctx) //GetRawDataAndCacheInGin(ctx)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, out)
	return err

}

// ========================= //
