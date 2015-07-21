// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares

import (
	"bytes"
	"encoding/json"
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"io/ioutil"
	"net/http"
)

func CSRF(app *forest.App) bear.HandlerFunc {
	type postBody struct {
		SessionID string `json:"sessionid"` // forest.SessionID == "sessionid"
	}
	csrf := func(
		_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
		if ctx.Request.Body == nil {
			app.Response(ctx, http.StatusBadRequest,
				forest.Failure, app.Error("CSRF")).Write(nil)
			return
		}
		pb := new(postBody)
		body, _ := ioutil.ReadAll(ctx.Request.Body)
		if body == nil || len(body) < 2 { // smallest JSON body is {}, 2 chars
			message := app.Error("Parse")
			app.Response(ctx,
				http.StatusBadRequest, forest.Failure, message).Write(nil)
			return
		}
		// set ctx.Request.Body back to an untouched io.ReadCloser
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		if err := json.Unmarshal(body, pb); err != nil {
			message := app.Error("Parse") + ": " + err.Error()
			app.Response(ctx,
				http.StatusBadRequest, forest.Failure, message).Write(nil)
			return
		}
		sessionID, ok := ctx.Get(forest.SessionID).(string)
		if !ok || sessionID != pb.SessionID {
			app.Response(ctx, http.StatusBadRequest,
				forest.Failure, app.Error("CSRF")).Write(nil)
		} else {
			ctx.Next()
		}
	}
	return bear.HandlerFunc(csrf)
}
