// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares

import (
	"fmt"
	"net/http"

	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
)

func BodyParser(app *forest.App) func(ctx *bear.Context) {
	return func(ctx *bear.Context) {
		destination, ok := ctx.Get(forest.Body).(Populater)
		if !ok {
			ctx.Set(forest.Error,
				fmt.Errorf("(*forest.App).BodyParser unitialized"))
			message := safeErrorMessage(app, ctx, app.Error("Parse"))
			app.Response(ctx, http.StatusInternalServerError,
				forest.Failure, message).Write(nil)
			return
		}
		if ctx.Request.Body == nil {
			ctx.Set(forest.SafeError,
				fmt.Errorf("%s: body is empty", app.Error("Parse")))
			message := safeErrorMessage(app, ctx, app.Error("Parse"))
			app.Response(ctx, http.StatusBadRequest,
				forest.Failure, message).Write(nil)
			return
		}
		if err := destination.Populate(ctx.Request.Body); err != nil {
			ctx.Set(forest.SafeError,
				fmt.Errorf("%s: %s", app.Error("Parse"), err))
			message := safeErrorMessage(app, ctx, app.Error("Parse"))
			app.Response(ctx, http.StatusBadRequest,
				forest.Failure, message).Write(nil)
			return
		}
		ctx.Next()
	}
}
