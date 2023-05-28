package delivery

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	userUcase "nostr-ex/pkg/app/user/usecase"
	"nostr-ex/pkg/config"
	dlv "nostr-ex/pkg/delivery"
)

type PostEventParams struct {
	PubKey string
	PriKey string
	Msg    string
}

func postEventParam(ctx *gin.Context) *PostEventParams {
	req := &PostEventParams{}
	err := dlv.GetBodyFromRawData(ctx, req)
	if err != nil {
		//ctx.JSON(http.StatusBadRequest, nil)
		return nil
	}
	return req
}

func PostEventParam(ctx *gin.Context) *PostEventParams {
	p := postEventParam(ctx)
	if p == nil { // user default
		p = &PostEventParams{}
	}

	//if p.PubKey == "" || rpeq.PriKey == "" || p.Msg == "" {
	if p.PubKey == "" || p.PriKey == "" {
		cfg := config.GetConfig()
		p.PubKey = cfg.Nostr.PublicKey
		p.PriKey = cfg.Nostr.PrivateKey
	}

	return p
}

func PostEvent(ctx *gin.Context) {
	req := PostEventParam(ctx)
	url := config.GetRelayUrl()

	m := userUcase.GetUserManager()
	user, err := m.AddUser(url, req.PubKey, req.PriKey)
	if err != nil {
		fmt.Println("AddUser Failed", err.Error())
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	err = user.PostEvent(req.Msg)
	if err != nil {
		fmt.Println("PostEvent Failed", err.Error())
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

// .. //
type ReqEventParams struct {
	PubKey string
}

func reqEventParam(ctx *gin.Context) *ReqEventParams {
	req := &ReqEventParams{}
	err := dlv.GetBodyFromRawData(ctx, req)
	if err != nil {
		//ctx.JSON(http.StatusBadRequest, nil)
		return nil
	}

	if req.PubKey == "" {
		//ctx.JSON(http.StatusBadRequest, nil)
		return nil
	}
	return req
}

func ReqEventParam(ctx *gin.Context) *ReqEventParams {
	p := reqEventParam(ctx)
	if p == nil { // user default
		cfg := config.GetConfig()
		p = &ReqEventParams{}
		p.PubKey = cfg.Nostr.PublicKey
	}

	return p
}

func ReqEvent(ctx *gin.Context) {
	req := CloseReqParam(ctx)
	url := config.GetRelayUrl()

	m := userUcase.GetUserManager()
	err := m.ReqEvent(url, req.PubKey)
	if err != nil {
		fmt.Println("ReqEvent Failed", err.Error())
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

// .. //
type CloseReqParams struct {
	PubKey string
}

func closeReqParam(ctx *gin.Context) *CloseReqParams {
	req := &CloseReqParams{}
	err := dlv.GetBodyFromRawData(ctx, req)
	if err != nil {
		//ctx.JSON(http.StatusBadRequest, nil)
		return nil
	}

	if req.PubKey == "" {
		//ctx.JSON(http.StatusBadRequest, nil)
		return nil
	}
	return req
}

func CloseReqParam(ctx *gin.Context) *CloseReqParams {
	p := closeReqParam(ctx)
	if p == nil { // user default
		cfg := config.GetConfig()
		p = &CloseReqParams{}
		p.PubKey = cfg.Nostr.PublicKey
	}

	return p
}

func CloseReq(ctx *gin.Context) {
	req := CloseReqParam(ctx)
	url := config.GetRelayUrl()
	m := userUcase.GetUserManager()
	err := m.CloseReq(url, req.PubKey)
	if err != nil {
		fmt.Println("CloseReq Failed", err.Error())
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
