package delivery

// type PostEventParams struct {
// 	PubKey string
// 	PriKey string
// 	Msg    string
// }

// func postEventParam(ctx *gin.Context) *PostEventParams {
// 	req := &PostEventParams{}
// 	err := dlv.GetBodyFromRawData(ctx, req)
// 	if err != nil {
// 		//ctx.JSON(http.StatusBadRequest, nil)
// 		return nil
// 	}
// 	return req
// }

// func PostEventParam(ctx *gin.Context) *PostEventParams {
// 	p := postEventParam(ctx)
// 	if p == nil { // user default
// 		p = &PostEventParams{}
// 	}

// 	//if p.PubKey == "" || rpeq.PriKey == "" || p.Msg == "" {
// 	if p.PubKey == "" || p.PriKey == "" {
// 		cfg := config.GetConfig()
// 		p.PubKey = cfg.Nostr.PublicKey
// 		p.PriKey = cfg.Nostr.PrivateKey
// 	}

// 	return p
// }

// func PostEvent(ctx *gin.Context) {
// 	req := PostEventParam(ctx)
// 	url := config.GetRelayUrl()

// 	m := relayUcase.GetRelayManager()
// 	user, err := m.AddRelay(url, "", "")
// 	if err != nil {
// 		fmt.Println("AddUser Failed", err.Error())
// 		ctx.JSON(http.StatusInternalServerError, nil)
// 		return
// 	}
// 	err = user.PostEvent(req.Msg)
// 	if err != nil {
// 		fmt.Println("PostEvent Failed", err.Error())
// 		ctx.JSON(http.StatusInternalServerError, nil)
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, nil)
// }
