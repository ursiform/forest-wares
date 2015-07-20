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
	badRequest := func(res http.ResponseWriter, req *http.Request,
		ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Generic"))
		app.Response(res,
			http.StatusBadRequest, forest.Failure, message).Write(nil)
	}
	return bear.HandlerFunc(badRequest)
}

func ErrorsConflict(app *forest.App) bear.HandlerFunc {
	conflict := func(res http.ResponseWriter, req *http.Request,
		ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Generic"))
		app.Response(res,
			http.StatusConflict, forest.Failure, message).Write(nil)
	}
	return bear.HandlerFunc(conflict)
}

func ErrorsMethodNotAllowed(app *forest.App) bear.HandlerFunc {
	methodNotAllowed := func(res http.ResponseWriter, req *http.Request,
		ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("MethodNotAllowed"))
		app.Response(res,
			http.StatusMethodNotAllowed, forest.Failure, message).Write(nil)
	}
	return bear.HandlerFunc(methodNotAllowed)
}

func ErrorsNotFound(app *forest.App) bear.HandlerFunc {
	notFound := func(res http.ResponseWriter, req *http.Request,
		ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("NotFound"))
		app.Response(res,
			http.StatusNotFound, forest.Failure, message).Write(nil)
	}
	return bear.HandlerFunc(notFound)
}

func ErrorsServerError(app *forest.App) bear.HandlerFunc {
	serverError := func(res http.ResponseWriter, req *http.Request,
		ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Generic"))
		app.Response(res,
			http.StatusInternalServerError, forest.Failure, message).Write(nil)
	}
	return bear.HandlerFunc(serverError)
}

func ErrorsUnauthorized(app *forest.App) bear.HandlerFunc {
	unauthorized := func(res http.ResponseWriter, req *http.Request,
		ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Unauthorized"))
		app.Response(res,
			http.StatusUnauthorized, forest.Failure, message).Write(nil)
	}
	return bear.HandlerFunc(unauthorized)
}
