// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares_test

import (
	"encoding/json"
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"github.com/ursiform/forest-wares"
	"io"
	"net/http"
)

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

func (app *router) authenticate(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
	ctx.Set(forest.SessionID, sessionID).Set(forest.SessionUserID, sessionUserID).Next(res, req)
}
func (app *router) initPostParse(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
	ctx.Set(forest.Body, new(postBody)).Next(res, req)
}
func (app *router) respondSuccess(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
	data := &responseFormat{Foo: "foo"}
	app.Response(res, http.StatusOK, forest.Success, forest.NoMessage).Write(data)
}

func (app *router) Route(path string) {
	app.Router.On("GET", path, app.respondSuccess)
	app.Router.On("GET", path+"/authenticate/failure", app.Ware("Authenticate"), app.respondSuccess)
	app.Router.On("GET", path+"/authenticate/success", app.authenticate, app.Ware("Authenticate"), app.respondSuccess)
	app.Router.On("GET", path+"/bad-request", app.Ware("BadRequest"))
	app.Router.On("GET", path+"/conflict", app.Ware("Conflict"))
	app.Router.On("GET", path+"/csrf", app.authenticate, app.Ware("CSRF"), app.respondSuccess)
	app.Router.On("GET", path+"/not-found", app.Ware("NotFound"))
	app.Router.On("GET", path+"/server-error", app.Ware("ServerError"))
	app.Router.On("GET", path+"/unauthorized", app.Ware("Unauthorized"))
	app.Router.On("POST", path+"/body-parser/failure/no-init", app.Ware("BodyParser"), app.respondSuccess)
	app.Router.On("POST", path+"/body-parser/success", app.initPostParse, app.Ware("BodyParser"), app.respondSuccess)
	app.Router.On("*", path, app.Ware("MethodNotAllowed"))
}

func newRouter(parent *forest.App) *router {
	wares.InstallBodyParser(parent)
	wares.InstallErrorWares(parent)
	wares.InstallSecurityWares(parent)
	return &router{parent}
}
