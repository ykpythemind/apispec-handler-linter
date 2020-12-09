// Package Api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package Api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi"
)

// CreateUserRequest defines model for CreateUserRequest.
type CreateUserRequest struct {
	Email *string `json:"email,omitempty"`
	Name  *string `json:"name,omitempty"`
}

// CreateUserResponse defines model for CreateUserResponse.
type CreateUserResponse struct {
	User *User `json:"user,omitempty"`
}

// GetUsersRequest defines model for GetUsersRequest.
type GetUsersRequest map[string]interface{}

// GetUsersResponse defines model for GetUsersResponse.
type GetUsersResponse []User

// User defines model for User.
type User struct {
	Email openapi_types.Email `json:"email"`

	// Unique identifier for the given user.
	Id   int     `json:"id"`
	Name *string `json:"name,omitempty"`
}

// CreateUserJSONBody defines parameters for CreateUser.
type CreateUserJSONBody CreateUserRequest

// GetUsersJSONBody defines parameters for GetUsers.
type GetUsersJSONBody GetUsersRequest

// CreateUserRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody CreateUserJSONBody

// GetUsersRequestBody defines body for GetUsers for application/json ContentType.
type GetUsersJSONRequestBody GetUsersJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create New User
	// (POST /user)
	CreateUser(w http.ResponseWriter, r *http.Request)
	// List Users
	// (GET /users)
	GetUsers(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// CreateUser operation middleware
func (siw *ServerInterfaceWrapper) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	siw.Handler.CreateUser(w, r.WithContext(ctx))
}

// GetUsers operation middleware
func (siw *ServerInterfaceWrapper) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	siw.Handler.GetUsers(w, r.WithContext(ctx))
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerFromMux(si, chi.NewRouter())
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerFromMuxWithBaseURL(si, r, "")
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	r.Group(func(r chi.Router) {
		r.Post(baseURL+"/user", wrapper.CreateUser)
	})
	r.Group(func(r chi.Router) {
		r.Get(baseURL+"/users", wrapper.GetUsers)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RUTW/bOBD9K8TsHmVLdoLdjU4bt2gRtEiCFu4lyIGWRjYDiWTIkRPD8H8vhpS/ndY5",
	"NJfI5MybN2/ecAmFaazRqMlDvgRfzLCR4fODQ0k49ui+4XOLnvjQOmPRkcIQgo1UNX/QwiLk4MkpPYVV",
	"Alo2eOJilQApqvnoGD5ZR5vJExbEMLtB3hrt8ZhE69Hx/78dVpDDX+m2o7RrJ2WEN6t3wCfKf0biEL8j",
	"wBrg8OqX2VvqirDx57FNTtQ6pCqdkwsOHXcivDGfyrhGEuTdSXI8MFVyXIm+cMqSMhpyGGv13KJQJWpS",
	"lUInKuMEzVBM1Ry1YOX7WzClCaeR+dvjd/jcKocl5A9cM+kYPW67Da0cypnAaw9fZWPr2Nl1rQoU3xtF",
	"s8BbEt5VI+X4Jwyurv7tDbLexWCDn4PkjL7njP+nfNYvTLO+/4GOGywhJ9diApVynm5DE7EURIUGl8ME",
	"arm5iwQS8Gqqx/ajJD4cZoOrXvZfb3gJK/5LQOnK7JpHWgUJzNH5KPSgn7FqxqLmqxwu+lk/gwSspFno",
	"N1273Jrow/1JRTsLKTS+bMbCTpAccFPuOR7iFNDTyJQLBiuMJtQBV1pbqyKkpU+ewZewK3xXiXHEyEzE",
	"J6xr88I3c1m3eHIY/3RqbIcxMZN+FTL3ZrEr+8hMYE/rrlIQNG7K7/bo+IkJyfvS3RtPwdMaC/ReuoWo",
	"FNal35j9+v5GkBHFscbBz3EngzjDLDtDzveT79b+BPu7L2GrfNs00i22TrjFF7F59IJ5AsEpBlb7zlg/",
	"L+/2xXmNHL6UcSX+mG5Hj+U5qn1VnkQUIfoL3TxI9rCE1rFjZ0Q2T9PaFLKeGU/5RZZlsHpc/QwAAP//",
	"pzUt00IHAAA=",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
