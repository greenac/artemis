package rest

import (
	"bytes"
	"github.com/greenac/artemis/pkg/errs"
	"io"
	"net/http"
)

var _ IClient = (*Client)(nil)

type Client struct {
	BaseHeaders *Headers
	HttpClient  IHttpClient
	BodyReader  func(r io.Reader) ([]byte, error)
	GetRequest  func(method, url string, body io.Reader) (*http.Request, error)
}

func (c *Client) Get(url string, headers *Headers, params UrlParams, cookies *[]http.Cookie) (Response, errs.IGenError) {
	req, err := c.GetRequest("GET", url, nil)
	if err != nil {
		ge := errs.GenError{}
		return Response{}, ge.AddMsg("Client:Get:failed to get new request with url: " + url + "with error: " + err.Error())
	}

	return c.makeRequest(req, c.makeHeaders(headers), cookies)
}

func (c *Client) PostBody(url string, headers *Headers, body []byte, cookies *[]http.Cookie) (Response, errs.IGenError) {
	req, err := c.GetRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		ge := errs.GenError{}
		return Response{}, ge.AddMsg("Client:PostBody:failed to get new request with url: " + url + "with error: " + err.Error())
	}

	return c.makeRequest(req, c.makeHeaders(headers), cookies)
}

func (c *Client) PostUrl(url string, headers *Headers, params UrlParams, cookies *[]http.Cookie) (Response, errs.IGenError) {
	req, err := c.GetRequest("POST", url, nil)
	if err != nil {
		ge := errs.GenError{}
		return Response{}, ge.AddMsg("Client:PostUrl:failed to get new request with url: " + url + "with error: " + err.Error())
	}

	return c.makeRequest(req, c.makeHeaders(headers), cookies)
}

func (c *Client) makeHeaders(headers *Headers) map[string]string {
	reqHeaders := map[string]string{}

	if c.BaseHeaders != nil {
		for key, hds := range *c.BaseHeaders {
			reqHeaders[key] = hds.String()
		}
	}

	if headers != nil {
		for key, hds := range *headers {
			reqHeaders[key] = hds.String()
		}
	}

	return reqHeaders
}

func (c *Client) makeRequest(req *http.Request, headers map[string]string, cookies *[]http.Cookie) (Response, errs.IGenError) {
	for k, h := range headers {
		req.Header.Add(k, h)
	}

	if cookies != nil {
		for _, ck := range *cookies {
			req.AddCookie(&ck)
		}
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		ge := errs.GenError{}
		return Response{}, ge.AddMsg("Client:makeRequest:failed to make request with error: " + err.Error())
	}

	var body []byte
	if res.Body != nil {
		b, err := c.BodyReader(res.Body)
		if err != nil {
			ge := errs.GenError{}
			return Response{}, ge.AddMsg("Client:makeRequest:failed to read response body with error: " + err.Error())
		}

		body = b
	}

	return Response{StatusCode: res.StatusCode, Status: res.Status, Body: body}, nil
}
