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
	"time"
)

// implements SessionManager
type sessionManager struct{}

func (manager *sessionManager) Create(sessionID string, userID string, userJSON string, ctx *bear.Context) error {
	testError, ok := ctx.Get("testerror").(bool)
	if !ok {
		testError = false
	}
	if testError {
		return errors.New("man.Create error")
	} else {
		ctx.Set(forest.SessionID, sessionID)
		ctx.Set(forest.SessionUserID, userID)
		return nil
	}
}
func (manager *sessionManager) CreateEmpty(sessionID string, ctx *bear.Context) {
	ctx.Set(forest.SessionID, sessionID)
}
func (manager *sessionManager) Delete(sessionID string, userID string) error {
	return nil
}
func (manager *sessionManager) Marshal(ctx *bear.Context) ([]byte, error) {
	return nil, nil
}
func (manager *sessionManager) Read(sessionID string) (userID string, userJSON string, err error) {
	return sessionUserID, sessionUserJSON, nil
}
func (manager *sessionManager) Revoke(userID string) error {
	return nil
}
func (manager *sessionManager) Update(sessionID string, userID string, userJSON string, duration time.Duration) error {
	return nil
}

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
func (app *router) sessionVerify(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
	_, ok := ctx.Get(forest.SessionID).(string)
	if !ok {
		ctx.Set(forest.Error, errors.New("sessionVerify failed"))
		app.Ware("ServerError")(res, req, ctx)
		return
	} else {
		ctx.Next(res, req)
	}
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
	app.Router.On("GET", path+"/session-get", app.Ware("SessionGet"), app.sessionVerify, app.respondSuccess)
	app.Router.On("GET", path+"/unauthorized", app.Ware("Unauthorized"))
	app.Router.On("POST", path+"/body-parser/failure/no-init", app.Ware("BodyParser"), app.respondSuccess)
	app.Router.On("POST", path+"/body-parser/success", app.initPostParse, app.Ware("BodyParser"), app.respondSuccess)
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
