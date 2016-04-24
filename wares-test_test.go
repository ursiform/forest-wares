// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

/*
Package wares_test contains tests and examples for package wares. The goal is
100% code coverage.
*/
package wares_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ursiform/forest"
)

const (
	arbitraryJSON             = "{\"foo\": \"bar\"}"
	customSafeErrorMessage    = "custom safe error message"
	customUnsafeErrorMessage  = "custom unsafe error message"
	root                      = "/test"
	sessionIDExistent         = "SOME-EXISTENT-SESSION-ID"
	sessionIDNonExistent      = "SOME-NONEXISTENT-SESSION-ID"
	sessionIDWithDeleteError  = "SOME-EXISTENT-SESSION-ID-THAT-FAILS-TO-DELETE"
	sessionIDWithMarshalError = "SOME-EXISTENT-SESSION-ID-THAT-FAILS-TO-MARSHAL"
	sessionIDWithUserDestruct = "SOME-EXISTENT-SESSION-ID-THAT-USER-DESTRUCTS"
	sessionIDWithSelfDestruct = "SOME-EXISTENT-SESSION-ID-THAT-SELF-DESTRUCTS"
	sessionIDWithUpdateError  = "SOME-EXISTENT-SESSION-ID-THAT-FAILS-TO-UPDATE"
	sessionUserID             = "SOME-USER-ID"
	sessionUserJSON           = "{\"id\": \"" + sessionUserID + "\"}"
)

type requested struct {
	auth   string
	body   []byte
	method string
	path   string
}

type wanted struct {
	code    int
	success bool
	data    interface{}
}

func makeRequest(t *testing.T, app *forest.App,
	params *requested, want *wanted) (*http.Response, *forest.Response) {
	var request *http.Request
	method := params.method
	auth := params.auth
	path := params.path
	body := params.body
	if body != nil {
		request, _ = http.NewRequest(method, path, bytes.NewBuffer(body))
	} else {
		request, _ = http.NewRequest(method, path, nil)
	}
	if len(auth) > 0 {
		request.AddCookie(&http.Cookie{Name: forest.SessionID, Value: auth})
	}
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	responseData := new(forest.Response)
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
		return nil, responseData
	}
	if err := json.Unmarshal(responseBody, responseData); err != nil {
		t.Errorf("unmarshal error: %v when attempting to read: %s",
			err, string(responseBody))
		return nil, responseData
	}
	if response.Code != want.code {
		t.Errorf("%s %s want: %d (%s) got: %d %s, body: %s", method, path,
			want.code, http.StatusText(want.code), response.Code,
			http.StatusText(response.Code), string(responseBody))
		return nil, responseData
	}
	if responseData.Success != want.success {
		t.Errorf("%s %s should return success: %t", method, path, want.success)
		return nil, responseData
	}
	return &http.Response{Header: response.Header()}, responseData
}

func TestAuthenticateFailure(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/authenticate/failure"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusUnauthorized, success: false}
	makeRequest(t, app, params, want)
}

func TestAuthenticateSuccess(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/authenticate/success"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusOK, success: true}
	makeRequest(t, app, params, want)
}

func TestBadRequest(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/bad-request"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusBadRequest, success: false}
	makeRequest(t, app, params, want)
}

func TestBodyParserFailureBadInput(t *testing.T) {
	debug := false
	method := "POST"
	path := root + "/body-parser/success"
	body := []byte("{BAD JSON}")
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{body: body, method: method, path: path}
	want := &wanted{code: http.StatusBadRequest, success: false}
	makeRequest(t, app, params, want)
}

func TestBodyParserFailureBodyNil(t *testing.T) {
	debug := false
	method := "POST"
	path := root + "/body-parser/success"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusBadRequest, success: false}
	makeRequest(t, app, params, want)
}

func TestBodyParserFailureNoInit(t *testing.T) {
	debug := false
	method := "POST"
	path := root + "/body-parser/failure/no-init"
	body := []byte(arbitraryJSON)
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{body: body, method: method, path: path}
	want := &wanted{code: http.StatusInternalServerError, success: false}
	makeRequest(t, app, params, want)
}

func TestBodyParserSuccess(t *testing.T) {
	debug := false
	method := "POST"
	path := root + "/body-parser/success"
	body := []byte(arbitraryJSON)
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{body: body, method: method, path: path}
	want := &wanted{code: http.StatusOK, success: true}
	makeRequest(t, app, params, want)
}

func TestConflict(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/conflict"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusConflict, success: false}
	makeRequest(t, app, params, want)
}

func TestCSRFFailureBodyNil(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/csrf"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusBadRequest, success: false}
	makeRequest(t, app, params, want)
}

func TestCSRFFailureBodyParse(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/csrf"
	body := []byte("{BAD JSON}")
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{body: body, method: method, path: path}
	want := &wanted{code: http.StatusBadRequest, success: false}
	makeRequest(t, app, params, want)
}

func TestCSRFFailureBodyTooShort(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/csrf"
	body := []byte("{")
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{body: body, method: method, path: path}
	want := &wanted{code: http.StatusBadRequest, success: false}
	makeRequest(t, app, params, want)
}

func TestCSRFFailureWrongSessionID(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/csrf"
	body := []byte(fmt.Sprintf("{\"sessionid\": \"WRONG-SESSION-ID\"}"))
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{body: body, method: method, path: path}
	want := &wanted{code: http.StatusBadRequest, success: false}
	makeRequest(t, app, params, want)
}

func TestCSRFSuccess(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/csrf"
	body := []byte(fmt.Sprintf("{\"sessionid\": \"%s\"}", sessionIDExistent))
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{body: body, method: method, path: path}
	want := &wanted{code: http.StatusOK, success: true}
	makeRequest(t, app, params, want)
}

func TestMethodNotAllowed(t *testing.T) {
	debug := false
	method := "OPTIONS"
	path := root
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusMethodNotAllowed, success: false}
	makeRequest(t, app, params, want)
}

func TestNotFound(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/not-found"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusNotFound, success: false}
	makeRequest(t, app, params, want)
}

func TestSafeErrorFilter(t *testing.T) {
	debug := false
	method := "GET"
	app := forest.New(debug)
	app.SafeErrorFilter = func(err error) error {
		if err.Error() == customSafeErrorMessage {
			return err
		} else {
			return nil
		}
	}
	app.RegisterRoute(root, newRouter(app))
	// test safe errors via custom filter
	path := root + "/safe-error/success"
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusInternalServerError, success: false}
	_, forestResponse := makeRequest(t, app, params, want)
	if forestResponse.Message != customSafeErrorMessage {
		t.Errorf("%s %s should return message: %s",
			method, path, customSafeErrorMessage)
	}
	// test unsafe errors not passing through custom filter
	path = root + "/safe-error/failure"
	params = &requested{method: method, path: path}
	want = &wanted{code: http.StatusInternalServerError, success: false}
	_, forestResponse = makeRequest(t, app, params, want)
	if forestResponse.Message == customUnsafeErrorMessage {
		t.Errorf("%s %s should NOT return unsafe message: %s",
			method, path, customUnsafeErrorMessage)
	}
	// test unsafe errors passing through if app.Config.Debug is true
	app.Config.Debug = true
	path = root + "/safe-error/failure"
	params = &requested{method: method, path: path}
	want = &wanted{code: http.StatusInternalServerError, success: false}
	_, forestResponse = makeRequest(t, app, params, want)
	if forestResponse.Message != customUnsafeErrorMessage {
		t.Errorf("%s %s should return unsafe message if %s is true: %s",
			method, path, "app.Config.Debug", customUnsafeErrorMessage)
	}
}

func TestServerError(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/server-error"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusInternalServerError, success: false}
	makeRequest(t, app, params, want)
}

func TestSessionSetDelBadUserIDError(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-del"
	auth := sessionIDWithUserDestruct
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusInternalServerError, success: false}
	makeRequest(t, app, params, want)
}

func TestSessionSetDelBadSessionIDError(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-del"
	auth := sessionIDWithSelfDestruct
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusInternalServerError, success: false}
	makeRequest(t, app, params, want)
}

func TestSessionDelFailure(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-del"
	auth := sessionIDWithDeleteError
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusInternalServerError, success: false}
	makeRequest(t, app, params, want)
}

func TestSessionDelSuccess(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-del"
	auth := sessionIDExistent
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusOK, success: true}
	makeRequest(t, app, params, want)
}

func TestSessionGetSuccessCreateEmpty(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-get"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusOK, success: true}
	makeRequest(t, app, params, want)
}

func TestSessionGetSuccessCreateError(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-get/create-error"
	auth := sessionIDExistent
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusOK, success: true}
	makeRequest(t, app, params, want)
}

func TestSessionGetSuccessCreateErrorWithDeleteError(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-get/create-error"
	auth := sessionIDWithDeleteError
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusOK, success: true}
	makeRequest(t, app, params, want)
}

func TestSessionGetSuccessExistent(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-get"
	auth := sessionIDExistent
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusOK, success: true}
	makeRequest(t, app, params, want)
}

func TestSessionGetSuccessExistentUpdateError(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-get"
	auth := sessionIDWithUpdateError
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusOK, success: true}
	makeRequest(t, app, params, want)
}

func TestSessionGetSuccessNonexistent(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-get"
	auth := sessionIDNonExistent
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusOK, success: true}
	makeRequest(t, app, params, want)
}

func TestSessionSetBadSessionIDError(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-set"
	auth := sessionIDWithSelfDestruct
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusInternalServerError, success: false}
	makeRequest(t, app, params, want)
}

func TestSessionSetBadUserIDError(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-set"
	auth := sessionIDWithUserDestruct
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusInternalServerError, success: false}
	makeRequest(t, app, params, want)
}

func TestSessionSetMarshalError(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-set"
	auth := sessionIDWithMarshalError
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusInternalServerError, success: false}
	makeRequest(t, app, params, want)
}

func TestSessionSetSuccess(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-set"
	auth := sessionIDExistent
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusOK, success: true}
	makeRequest(t, app, params, want)
}

func TestSessionSetUpdateError(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/session-set"
	auth := sessionIDWithUpdateError
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{auth: auth, method: method, path: path}
	want := &wanted{code: http.StatusInternalServerError, success: false}
	makeRequest(t, app, params, want)
}

func TestUnauthorized(t *testing.T) {
	debug := false
	method := "GET"
	path := root + "/unauthorized"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusUnauthorized, success: false}
	makeRequest(t, app, params, want)
}
