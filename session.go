// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"net/http"
)

func SessionDel(app *forest.App, manager SessionManager) bear.HandlerFunc {
	sessionDel := func(
		_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
		sessionID, ok := ctx.Get(forest.SessionID).(string)
		if !ok {
			err := fmt.Errorf("SessionDel %s: %v",
				forest.SessionID, ctx.Get(forest.SessionID))
			ctx.Set(forest.Error, err)
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(ctx, http.StatusInternalServerError,
				forest.Failure, message).Write(nil)
			return
		}
		userID, ok := ctx.Get(forest.SessionUserID).(string)
		if !ok {
			err := fmt.Errorf("SessionDel %s: %v",
				forest.SessionUserID, ctx.Get(forest.SessionUserID))
			ctx.Set(forest.Error, err)
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(ctx, http.StatusInternalServerError,
				forest.Failure, message).Write(nil)
			return
		}
		if err := manager.Delete(sessionID, userID); err != nil {
			ctx.Set(forest.Error, err)
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(ctx, http.StatusInternalServerError,
				forest.Failure, message).Write(nil)
			return
		}
		ctx.Next()
	}
	return bear.HandlerFunc(sessionDel)
}

func SessionGet(app *forest.App, manager SessionManager) bear.HandlerFunc {
	sessionGet := func(
		_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
		cookieName := forest.SessionID
		createEmptySession := func(sessionID string) {
			path := app.CookiePath
			if path == "" {
				path = "/"
			}
			cookieValue := sessionID
			duration := app.Duration("Cookie")
			// Reset the cookie.
			app.SetCookie(ctx, path, cookieName, cookieValue, duration)
			manager.CreateEmpty(sessionID, ctx)
			ctx.Next()
		}
		cookie, err := ctx.Request.Cookie(cookieName)
		if err != nil || cookie.Value == "" {
			createEmptySession(uuid.New())
			return
		}
		sessionID := cookie.Value
		userID, userJSON, err := manager.Read(sessionID)
		if err != nil || userID == "" || userJSON == "" {
			createEmptySession(uuid.New())
			return
		}
		if err := manager.Create(sessionID, userID, userJSON, ctx); err != nil {
			println(fmt.Sprintf("error creating session: %s", err))
			defer func(sessionID string, userID string) {
				if err := manager.Delete(sessionID, userID); err != nil {
					println(fmt.Sprintf("error deleting session: %s", err))
				}
			}(sessionID, userID)
			createEmptySession(uuid.New())
			return
		}
		// If SessionRefresh is set to false, the session will not refresh;
		// if it's not set or if it's set to true, the session is refreshed.
		refresh, ok := ctx.Get(forest.SessionRefresh).(bool)
		if !ok || refresh {
			path := app.CookiePath
			if path == "" {
				path = "/"
			}
			cookieName := forest.SessionID
			cookieValue := sessionID
			duration := app.Duration("Cookie")
			// Refresh the cookie.
			app.SetCookie(ctx, path, cookieName, cookieValue, duration)
			err := manager.Update(sessionID, userID,
				userJSON, app.Duration("Session"))
			if err != nil {
				println(fmt.Sprintf("error updating session: %s", err))
			}
		}
		ctx.Next()
	}
	return bear.HandlerFunc(sessionGet)
}

func SessionSet(app *forest.App, manager SessionManager) bear.HandlerFunc {
	sessionSet := func(
		_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
		userJSON, err := manager.Marshal(ctx)
		if err != nil {
			ctx.Set(forest.Error, err)
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(ctx, http.StatusInternalServerError,
				forest.Failure, message).Write(nil)
			return
		}
		sessionID, ok := ctx.Get(forest.SessionID).(string)
		if !ok {
			err := fmt.Errorf("%s: %v",
				forest.SessionID, ctx.Get(forest.SessionID))
			ctx.Set(forest.Error, err)
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(ctx, http.StatusInternalServerError,
				forest.Failure, message).Write(nil)
			return
		}
		userID, ok := ctx.Get(forest.SessionUserID).(string)
		if !ok {
			err := fmt.Errorf("%s: %v",
				forest.SessionUserID, ctx.Get(forest.SessionUserID))
			ctx.Set(forest.Error, err)
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(ctx, http.StatusInternalServerError,
				forest.Failure, message).Write(nil)
			return
		}
		if err := manager.Update(sessionID, userID,
			string(userJSON), app.Duration("Session")); err != nil {
			ctx.Set(forest.Error, err)
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(ctx, http.StatusInternalServerError,
				forest.Failure, message).Write(nil)
			return
		}
		ctx.Next()
	}
	return bear.HandlerFunc(sessionSet)
}
