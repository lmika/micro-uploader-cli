package micropub

import (
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) Post(dest string, body string) (pr PostResponse, err error) {
	f := make(url.Values)
	f.Set("h", "entry")
	f.Set("content", body)
	f.Set("mp-destination", dest)

	req, err := http.NewRequest("POST", c.micropubURL, strings.NewReader(f.Encode()))
	if err != nil {
		return PostResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.bearerAuth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	err = c.doRequestReturningJson(&pr, req)
	return pr, err
}

type PostResponse struct {
	URL     string `json:"url"`
	Preview string `json:"preview"`
	Edit    string `json:"edit"`
}
