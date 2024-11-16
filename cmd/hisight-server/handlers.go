package main

import (
	// "encoding/json"
	// "log"

	"net/http"
	// "github.com/jacobmiller22/gauth/internal/authorization"
	// "github.com/jacobmiller22/gauth/internal/database"
)

func (app *hsApp) getCommands(w http.ResponseWriter, r *http.Request) {
	// queryParams := r.URL.Query()
	// clientId := queryParams.Get("client_id")
	// responseType := queryParams.Get("response_type")
	// redirectUri := queryParams.Get("redirect_uri")
	// scope := queryParams.Get("scope")
	// state := queryParams.Get("state")
	// codeChallengeMethod := queryParams.Get("code_challenge_method")
	// codeChallenge := queryParams.Get("code_challenge")
	// log.Printf("Authorize Handler called w/ Query Params: %v", queryParams)
	//
	// authReq := authorization.NewAuthorizationReq(responseType, clientId, redirectUri, scope, state, codeChallengeMethod, codeChallenge)
	//
	// log.Printf("Authorization Request: %+v", authReq)
	//
	// err := authReq.VerifyAuthorizationRequest(app.db.GetApplication)
	// switch err {
	// case authorization.ErrInvalidClientId, authorization.ErrInvalidRedirectUri, authorization.ErrInvalidResponseType, authorization.ErrInvalidScope, authorization.ErrUnauthorizationClient, authorization.ErrUnsupportedResponseType, authorization.ErrUnsupportedResponseType:
	// 	log.Printf("Error verifying request: %+v", err)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	//
	// authRes, err := authReq.ToAuthorizationRes()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
	//
	// // Store the pkce and state values, associate them with the code or token
	// app.db.SaveResponse(database.AuthorizationResponse{
	// 	ResponseType:        authRes.ResponseType,
	// 	Code:                authRes.Code,
	// 	Token:               authRes.Token,
	// 	State:               authRes.State,
	// 	CodeChallengeMethod: codeChallengeMethod,
	// 	CodeChallenge:       codeChallenge,
	// })
	//
	// // TODO: Send user to authorization server UI login screen
	//
	// location, err := authRes.Location()
	//
	// if err != nil {
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	//
	// w.Header().Set("Location", location.String())
	// w.WriteHeader(http.StatusFound)
}
