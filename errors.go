// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares

import (
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"net/http"
)

func ErrorsBadRequest(app *forest.App) bear.HandlerFunc {
	badRequest := func(ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Generic"))
		app.Response(ctx,
			http.StatusBadRequest, forest.Failure, message).Write(nil)
	}
	handler, _, _ := bear.Handlerize(badRequest)
	return handler
}

func ErrorsConflict(app *forest.App) bear.HandlerFunc {
	conflict := func(
		_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Generic"))
		app.Response(ctx,
			http.StatusConflict, forest.Failure, message).Write(nil)
	}
	handler, _, _ := bear.Handlerize(conflict)
	return handler
}

func ErrorsMethodNotAllowed(app *forest.App) bear.HandlerFunc {
	methodNotAllowed := func(
		_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("MethodNotAllowed"))
		app.Response(ctx,
			http.StatusMethodNotAllowed, forest.Failure, message).Write(nil)
	}
	handler, _, _ := bear.Handlerize(methodNotAllowed)
	return handler
}

func ErrorsNotFound(app *forest.App) bear.HandlerFunc {
	notFound := func(
		_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("NotFound"))
		app.Response(ctx,
			http.StatusNotFound, forest.Failure, message).Write(nil)
	}
	handler, _, _ := bear.Handlerize(notFound)
	return handler
}

func ErrorsServerError(app *forest.App) bear.HandlerFunc {
	serverError := func(
		_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Generic"))
		app.Response(ctx,
			http.StatusInternalServerError, forest.Failure, message).Write(nil)
	}
	handler, _, _ := bear.Handlerize(serverError)
	return handler
}

func ErrorsUnauthorized(app *forest.App) bear.HandlerFunc {
	unauthorized := func(
		_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Unauthorized"))
		app.Response(ctx,
			http.StatusUnauthorized, forest.Failure, message).Write(nil)
	}
	handler, _, _ := bear.Handlerize(unauthorized)
	return handler
}
