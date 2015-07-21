// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares_test

import (
	"encoding/json"
	"errors"
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"github.com/ursiform/forest-wares"
	"io"
	"net/http"
)

// implements Populater
type postBody struct {
	Foo string `json:"foo"`
}

func (pb *postBody) Populate(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(pb)
}

type responseFormat struct {
	Foo string `json:"foo"`
}

type router struct{ *forest.App }

func (app *router) authenticate(
	_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
	ctx.Set(forest.SessionID, sessionIDExistent)
	ctx.Set(forest.SessionUserID, sessionUserID)
	ctx.Next()
}
func (app *router) customSafeErrorFilterFailure(
	_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
	ctx.Set(forest.Error, errors.New(customUnsafeErrorMessage))
	app.Ware("ServerError")(ctx.ResponseWriter, ctx.Request, ctx)
}
func (app *router) customSafeErrorFilterSuccess(
	_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
	ctx.Set(forest.Error, errors.New(customSafeErrorMessage))
	app.Ware("ServerError")(ctx.ResponseWriter, ctx.Request, ctx)
}
func (app *router) initPostParse(
	_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
	ctx.Set(forest.Body, new(postBody)).Next()
}
func (app *router) respondSuccess(
	_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
	data := &responseFormat{Foo: "foo"}
	app.Response(
		ctx, http.StatusOK, forest.Success, forest.NoMessage).Write(data)
}
func (app *router) sessionCreateError(
	_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
	ctx.Set("testerror", true).Next()
}
func (app *router) sessionDelIntercept(
	_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
	sessionID := ctx.Get(forest.SessionID).(string)
	if sessionID == sessionIDWithSelfDestruct {
		ctx.Set(forest.SessionID, nil)
	}
	if sessionID == sessionIDWithUserDestruct {
		ctx.Set(forest.SessionUserID, nil)
	}
	ctx.Next()
}
func (app *router) sessionVerify(
	_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
	_, ok := ctx.Get(forest.SessionID).(string)
	if !ok {
		ctx.Set(forest.Error, errors.New("sessionVerify failed"))
		app.Ware("ServerError")(ctx.ResponseWriter, ctx.Request, ctx)
		return
	} else {
		ctx.Next()
	}
}

func (app *router) Route(path string) {
	app.Router.On("GET", path, app.respondSuccess)
	app.Router.On("GET", path+"/authenticate/failure",
		app.Ware("Authenticate"), app.respondSuccess)
	app.Router.On("GET", path+"/authenticate/success",
		app.authenticate, app.Ware("Authenticate"), app.respondSuccess)
	app.Router.On("GET", path+"/bad-request",
		app.Ware("BadRequest"))
	app.Router.On("GET", path+"/conflict",
		app.Ware("Conflict"))
	app.Router.On("GET", path+"/csrf",
		app.authenticate, app.Ware("CSRF"), app.respondSuccess)
	app.Router.On("GET", path+"/not-found",
		app.Ware("NotFound"))
	app.Router.On("GET", path+"/safe-error/failure",
		app.customSafeErrorFilterFailure)
	app.Router.On("GET", path+"/safe-error/success",
		app.customSafeErrorFilterSuccess)
	app.Router.On("GET", path+"/server-error",
		app.Ware("ServerError"))
	app.Router.On("GET", path+"/session-del",
		app.Ware("SessionGet"),
		app.sessionDelIntercept,
		app.Ware("SessionDel"), app.respondSuccess)
	app.Router.On("GET", path+"/session-get",
		app.Ware("SessionGet"), app.sessionVerify, app.respondSuccess)
	app.Router.On("GET", path+"/session-get/create-error",
		app.sessionCreateError, app.Ware("SessionGet"),
		app.sessionVerify, app.respondSuccess)
	app.Router.On("GET", path+"/session-set",
		app.Ware("SessionGet"), app.Ware("SessionSet"), app.respondSuccess)
	app.Router.On("GET", path+"/unauthorized",
		app.Ware("Unauthorized"))
	app.Router.On("POST", path+"/body-parser/failure/no-init",
		app.Ware("BodyParser"), app.respondSuccess)
	app.Router.On("POST", path+"/body-parser/success",
		app.initPostParse, app.Ware("BodyParser"), app.respondSuccess)
	app.Router.On("*", path, app.Ware("MethodNotAllowed"))
}

func newRouter(parent *forest.App) *router {
	manager := new(sessionManager)
	wares.InstallBodyParser(parent)
	wares.InstallErrorWares(parent)
	wares.InstallSecurityWares(parent)
	wares.InstallSessionWares(parent, manager)
	return &router{parent}
}
