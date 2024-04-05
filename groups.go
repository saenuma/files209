package files209

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

func (cl *Client) ListGroups() ([]string, error) {
	urlValues := url.Values{}
	urlValues.Add("key-str", cl.KeyStr)

	resp, err := httpCl.PostForm(fmt.Sprintf("%slist-groups", cl.Addr), urlValues)
	if err != nil {
		return nil, ConnError{err.Error()}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ConnError{err.Error()}
	}

	if resp.StatusCode == 200 {
		ret := make([]string, 0)
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

func (cl *Client) DeleteGroup(groupName string) error {
	urlValues := url.Values{}
	urlValues.Add("key-str", cl.KeyStr)

	resp, err := httpCl.PostForm(fmt.Sprintf("%sdelete-group/%s", cl.Addr, groupName), urlValues)
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
