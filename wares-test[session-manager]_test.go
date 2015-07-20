// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares_test

import (
	"errors"
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
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
		return errors.New("manager.Create error")
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
	if sessionID == sessionIDWithDeleteError {
		return errors.New("manager.Delete error")
	}
	return nil
}
func (manager *sessionManager) Marshal(ctx *bear.Context) ([]byte, error) {
	return nil, nil
}
func (manager *sessionManager) Read(sessionID string) (userID string, userJSON string, err error) {
	if sessionID == sessionIDNonExistent {
		return "", "", nil
	} else {
		return sessionUserID, sessionUserJSON, nil
	}
}
func (manager *sessionManager) Revoke(userID string) error {
	return nil
}
func (manager *sessionManager) Update(sessionID string, userID string, userJSON string, duration time.Duration) error {
	return nil
}
