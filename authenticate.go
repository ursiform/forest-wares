// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares

import (
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"net/http"
)

func Authenticate(app *forest.App) bear.HandlerFunc {
	authenticate := func(
		_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
		userID, ok := ctx.Get(forest.SessionUserID).(string)
		if !ok || len(userID) == 0 {
			app.Response(ctx, http.StatusUnauthorized, forest.Failure,
				app.Error("Unauthorized")).Write(nil)
			return
		}
		ctx.Next()
	}
	handler, _, _ := bear.Handlerize(authenticate)
	return handler
}
