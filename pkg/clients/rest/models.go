package rest

import (
	"github.com/greenac/artemis/pkg/errs"
	"net/http"
	"strings"
)

type HeaderValue []string

func (hv HeaderValue) String() string {
	return strings.Join(hv, ", ")
}

type Headers map[string]HeaderValue

type UrlParamType string

const (
	UrlParamTypeEqual              UrlParamType = "="
	UrlParamTypeLessThan           UrlParamType = "<"
	UrlParamTypeLessThanEqualTo    UrlParamType = "<="
	UrlParamTypeGreaterThan        UrlParamType = ">"
	UrlParamTypeGreaterThanEqualTo UrlParamType = ">="
)

type UrlParam struct {
	LVal    string
	RVal    string
	Compare UrlParamType
}

func (up UrlParam) String() string {
	return up.LVal + string(up.Compare) + up.RVal
}

type UrlParams []UrlParam

func (ups UrlParams) String() string {
	s := ""
	for i, p := range ups {
		s += p.String()
		if i < len(ups)-1 {
			s += "?"
		}
	}

	return s
}

type Response struct {
	StatusCode int
	Status     string
	Body       []byte
}

type IHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type IClient interface {
	Get(url string, headers *Headers, params UrlParams, cookies *[]http.Cookie) (Response, errs.IGenError)
	PostBody(url string, headers *Headers, body []byte, cookies *[]http.Cookie) (Response, errs.IGenError)
	PostUrl(url string, headers *Headers, params UrlParams, cookies *[]http.Cookie) (Response, errs.IGenError)
}
