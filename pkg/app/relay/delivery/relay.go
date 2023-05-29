package delivery

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	relayUcase "nostr-ex/pkg/app/relay/usecase"
	dlv "nostr-ex/pkg/delivery"
)

type RelayListResp struct {
	Total int                `json:"total_num"`
	List  []relayUcase.Relay `json:"Relay"`
}

func GetRelays(ctx *gin.Context) {

	m := relayUcase.GetRelayManager()
	list := m.GetRelays()

	ctx.JSON(http.StatusOK, nil)

	data := RelayListResp{
		Total: len(list),
		List:  list,
	}
	response := dlv.ResBody{}
	response.ResCode = dlv.RES_OK
	response.Data = data

	ctx.JSON(http.StatusOK, response)
}

type AddRelayParams struct {
	URL string
}

func addRelayParamFromBody(ctx *gin.Context) *AddRelayParams {
	req := &AddRelayParams{}
	err := dlv.GetBodyFromRawData(ctx, req)
	if err != nil {
		//ctx.JSON(http.StatusBadRequest, nil)
		return nil
	}
	return req
}

func AddRelayParamFromBody(ctx *gin.Context) *AddRelayParams {
	p := addRelayParamFromBody(ctx)
	if p == nil { // user default
		p = &AddRelayParams{}
	}

	if p.URL == "" {
		return nil
	}

	return p
}

func AddRelay(ctx *gin.Context) {
	req := AddRelayParamFromBody(ctx)
	if req == nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	m := relayUcase.GetRelayManager()
	_, err := m.AddRelay(req.URL)
	if err != nil {
		fmt.Println("AddRelay Failed", err.Error())
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
