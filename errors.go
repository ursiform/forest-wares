// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares

import (
	"net/http"

	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
)

func ErrorsBadRequest(app *forest.App) func(ctx *bear.Context) {
	return func(ctx *bear.Context) {
		app.Response(
			ctx,
			http.StatusBadRequest,
			forest.Failure,
			safeErrorMessage(app, ctx, app.Error("Generic"))).Write(nil)
	}
}

func ErrorsConflict(app *forest.App) func(ctx *bear.Context) {
	return func(ctx *bear.Context) {
		app.Response(
			ctx,
			http.StatusConflict,
			forest.Failure,
			safeErrorMessage(app, ctx, app.Error("Generic"))).Write(nil)
	}
}

func ErrorsMethodNotAllowed(app *forest.App) func(ctx *bear.Context) {
	return func(ctx *bear.Context) {
		app.Response(
			ctx,
			http.StatusMethodNotAllowed,
			forest.Failure,
			safeErrorMessage(app, ctx, app.Error("MethodNotAllowed"))).Write(nil)
	}
}

func ErrorsNotFound(app *forest.App) func(ctx *bear.Context) {
	return func(ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("NotFound"))
		app.Response(ctx, http.StatusNotFound, forest.Failure, message).Write(nil)
	}
}

func ErrorsServerError(app *forest.App) func(ctx *bear.Context) {
	return func(ctx *bear.Context) {
		app.Response(
			ctx,
			http.StatusInternalServerError,
			forest.Failure,
			safeErrorMessage(app, ctx, app.Error("Generic"))).Write(nil)
	}
}

func ErrorsUnauthorized(app *forest.App) func(ctx *bear.Context) {
	return func(ctx *bear.Context) {
		app.Response(
			ctx,
			http.StatusUnauthorized,
			forest.Failure,
			safeErrorMessage(app, ctx, app.Error("Unauthorized"))).Write(nil)
	}
}
