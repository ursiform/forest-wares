// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares

import (
	"net/http"

	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
)

func Authenticate(app *forest.App) func(ctx *bear.Context) {
	return func(ctx *bear.Context) {
		userID, ok := ctx.Get(forest.SessionUserID).(string)
		if !ok || len(userID) == 0 {
			app.Response(ctx, http.StatusUnauthorized, forest.Failure,
				app.Error("Unauthorized")).Write(nil)
			return
		}
		ctx.Next()
	}
}
