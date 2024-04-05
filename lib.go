// this package 'files209' is the golang library for communicating with the files209 server.
package files209

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ConnError struct {
	msg string // description of error
}

func (e ConnError) Error() string {
	return "connection Error: " + e.msg
}

type ValidationError struct {
	msg string // description of error
}

func (e ValidationError) Error() string {
	return "validation Error: " + e.msg
}

type ServerError struct {
	msg string // description of error
}

func (e ServerError) Error() string {
	return "server Error: " + e.msg
}

var httpCl *http.Client

func init() {
	config := &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{TLSClientConfig: config}

	httpCl = &http.Client{Transport: tr}
}

type Client struct {
	Addr   string
	KeyStr string
}

func NewClient(ip, keyStr string) Client {
	return Client{"https://" + ip + ":31822/", keyStr}
}

// Used whenever you changed the default port
func NewClientCustomPort(ip, keyStr string, port int) Client {
	return Client{"https://" + ip + fmt.Sprintf(":%d/", port), keyStr}
}

func (cl *Client) Ping() error {
	urlValues := url.Values{}
	urlValues.Set("key-str", cl.KeyStr)

	resp, err := httpCl.PostForm(cl.Addr+"is-files209", urlValues)
	if err != nil {
		return ConnError{err.Error()}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ConnError{err.Error()}
	}

	if resp.StatusCode == 200 {
		if string(body) == "yeah-files209" {
			return nil
		} else {
			return ConnError{"Unexpected Error in confirming that the server is a files209 store."}
		}
	} else {
		return ConnError{string(body)}
	}
}
