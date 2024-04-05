package files209

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
)

func (cl *Client) WriteFile(groupName, fileName string, toWrite []byte) error {
	urlValues := url.Values{}
	urlValues.Add("key-str", cl.KeyStr)
	urlValues.Add("name", fileName)
	urlValues.Add("dataB64", base64.StdEncoding.EncodeToString(toWrite))

	resp, err := httpCl.PostForm(fmt.Sprintf("%swrite-file/%s", cl.Addr, groupName), urlValues)
	if err != nil {
		return ConnError{err.Error()}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ConnError{err.Error()}
	}

	if resp.StatusCode == 200 {
		return nil
	} else if resp.StatusCode == 400 {
		return ValidationError{string(body)}
	} else {
		return ServerError{string(body)}
	}

}

func (cl *Client) DeleteFile(groupName, fileName string) error {
	urlValues := url.Values{}
	urlValues.Add("key-str", cl.KeyStr)

	resp, err := httpCl.PostForm(fmt.Sprintf("%sdelete-file/%s/%s", cl.Addr, groupName, fileName), urlValues)
	if err != nil {
		return ConnError{err.Error()}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ConnError{err.Error()}
	}

	if resp.StatusCode == 200 {
		return nil
	} else if resp.StatusCode == 400 {
		return ValidationError{string(body)}
	} else {
		return ServerError{string(body)}
	}

}
