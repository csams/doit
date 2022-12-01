package cli

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pkg/browser"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func CreateClient(insecure bool) *http.Client {
	if insecure {
		// like http.DefaultTransport but with InsecureSkipVerify: true
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		return &http.Client{
			Transport: transport,
		}
	}
	return http.DefaultClient
}

type OIDCFlow struct {
	CompletedConfig

	Server *http.Server

	ClientContext context.Context
	Provider      *oidc.Provider
	Verifier      *oidc.IDTokenVerifier
	OAuth2Config  *oauth2.Config
}

/*
"issuer": "https://localhost/realms/todoapp",
"authorization_endpoint": "https://localhost/realms/todoapp/protocol/openid-connect/auth",
"token_endpoint": "https://localhost/realms/todoapp/protocol/openid-connect/token",
"introspection_endpoint": "https://localhost/realms/todoapp/protocol/openid-connect/token/introspect",
"userinfo_endpoint": "https://localhost/realms/todoapp/protocol/openid-connect/userinfo",
"end_session_endpoint": "https://localhost/realms/todoapp/protocol/openid-connect/logout",
"frontchannel_logout_session_supported": true,
"frontchannel_logout_supported": true,
"jwks_uri": "https://localhost/realms/todoapp/protocol/openid-connect/certs",
*/
func NewOIDCFlow(c CompletedConfig) (*OIDCFlow, error) {
	ctx := context.Background()
	client := CreateClient(c.InsecureClient)
	ctx = oidc.ClientContext(ctx, client)

	oidcConfig := &oidc.Config{ClientID: c.ClientId}

	provider, err := oidc.NewProvider(ctx, c.AuthorizationServerURL)
	if err != nil {
		return nil, err
	}

	oauthConfig := &oauth2.Config{
		ClientID:    c.ClientId,
		Endpoint:    provider.Endpoint(),
		RedirectURL: c.RedirectURL,
		Scopes:      []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oauthConfig.Client(ctx, nil)

	return &OIDCFlow{
		CompletedConfig: c,
		ClientContext:   ctx,
		Provider:        provider,
		Verifier:        provider.Verifier(oidcConfig),
		Server:          &http.Server{Addr: c.LocalAddr},
		OAuth2Config:    oauthConfig,
	}, nil

}

type TokenWrapper struct {
	*oauth2.Token
	IdToken string `json:"id_token,omitempty"`
}

func (l *OIDCFlow) SaveToken(tok *oauth2.Token) error {
	rawIDToken := tok.Extra("id_token").(string)
	data, err := json.Marshal(&TokenWrapper{Token: tok, IdToken: rawIDToken})
	if err != nil {
		return err
	}
	return os.WriteFile(l.TokenFile, data, 0600)
}

func (l *OIDCFlow) GetSavedToken() (string, error) {
	rawToken, err := os.ReadFile(l.TokenFile)
	if err == nil {
		var tok TokenWrapper
		err = json.Unmarshal(rawToken, &tok)
		if err != nil {
			return "", err
		}

		// if the saved token is still valid, just use it
		if _, err := l.Verifier.Verify(l.ClientContext, tok.IdToken); err == nil {
			return tok.IdToken, err
		}

		// if it's not still valid, this should automatically refresh
		if newTok, err := l.OAuth2Config.TokenSource(l.ClientContext, tok.Token).Token(); err == nil {
			if newTok.AccessToken != tok.AccessToken {
				if err = l.SaveToken(newTok); err != nil {
					return "", err
				}
			}

			// check again just to make sure
			rawIDToken := newTok.Extra("id_token").(string)
			if _, err := l.Verifier.Verify(l.ClientContext, rawIDToken); err != nil {
				return "", err
			}

			// we have a valid id token
			return rawIDToken, nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

func (l *OIDCFlow) GetIdToken() (string, error) {
	tok, err := l.GetSavedToken()
	if err == nil {
		return tok, nil
	}

	// if we couldn't get a saved token, we have to run the oauth flow
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	stop := make(chan struct{})
	var idData string
	var idErr error

	mux := http.NewServeMux()
	l.Server.Handler = mux

	// we'll open the browser here to set the initial state and nonce
	// before redirecting over to the SSO server for login
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		state, err := randString(16)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		nonce, err := randString(16)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		setCallbackCookie(w, r, "state", state)
		setCallbackCookie(w, r, "nonce", nonce)

		http.Redirect(w, r, l.OAuth2Config.AuthCodeURL(state, oidc.Nonce(nonce), oauth2.AccessTypeOffline), http.StatusFound)
	})

	// The SSO server will redirect the browser back here. On successful login, the redirect
	// will contain the state and nonce we sent to it in the original redirect above along with
	// an auth code that can be exchanged on the back channel for an access token that comes with
	// an OIDC identity token.
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		defer close(stop)
		state, err := r.Cookie("state")
		if err != nil {
			idErr = err
			http.Error(w, "state not found", http.StatusBadRequest)
			return
		}

		if r.URL.Query().Get("state") != state.Value {
			msg := "state did not match"
			idErr = errors.New(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		oauth2Token, err := l.OAuth2Config.Exchange(l.ClientContext, r.URL.Query().Get("code"))
		if err != nil {
			idErr = err
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			msg := "no id_token field in oauth2 token"
			idErr = errors.New(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		idToken, err := l.Verifier.Verify(l.ClientContext, rawIDToken)
		if err != nil {
			idErr = err
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		nonce, err := r.Cookie("nonce")
		if err != nil {
			idErr = err
			http.Error(w, "nonce not found", http.StatusBadRequest)
			return
		}

		if idToken.Nonce != nonce.Value {
			msg := "nonce did not match"
			idErr = errors.New(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		if err = l.SaveToken(oauth2Token); err != nil {
			idErr = err
			http.Error(w, "Failed to save oauth token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		idData = rawIDToken
		w.Write([]byte("Authorization Successful"))
	})

	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)
		errChan <- l.Server.ListenAndServe()
	}()

	if err := browser.OpenURL("http://" + l.LocalAddr); err != nil {
		l.Server.Close()
		return "", err
	}

	select {
	case err := <-errChan:
		l.Server.Close()
		return "", err
	case <-sigChannel:
		return "", l.Server.Close()
	case <-stop:
		l.Server.Shutdown(l.ClientContext)
		return idData, idErr
	}
}
