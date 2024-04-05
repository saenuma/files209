package files209

import (
	"encoding/base64"
	"encoding/json"
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

func (cl *Client) ListFiles(groupName, fileName string) (map[string]int64, error) {
	urlValues := url.Values{}
	urlValues.Add("key-str", cl.KeyStr)

	resp, err := httpCl.PostForm(fmt.Sprintf("%slist-files/%s", cl.Addr, groupName), urlValues)
	if err != nil {
		return nil, ConnError{err.Error()}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ConnError{err.Error()}
	}

	if resp.StatusCode == 200 {
		ret := make(map[string]int64)
		err = json.Unmarshal(body, &ret)
		if err != nil {
			return nil, ConnError{"json error\n" + err.Error()}
		}

		return ret, nil
	} else if resp.StatusCode == 400 {
		return nil, ValidationError{string(body)}
	} else {
		return nil, ServerError{string(body)}
	}
}
