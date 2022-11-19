package micropub

import (
	"emperror.dev/errors"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
)

type Client struct {
	micropubURL string
	bearerAuth  string
	client      *http.Client

	mutex            *sync.Mutex
	discoveredConfig MicroPubConfig
}

func New(micropubURL string, bearerAuth string) *Client {
	return &Client{
		micropubURL: micropubURL,
		bearerAuth:  bearerAuth,
		client:      http.DefaultClient,
		mutex:       new(sync.Mutex),
	}
}

func (c *Client) newReq(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, errors.Wrap(err, "bad request")
	}
	req.Header.Set("Authorization", "Bearer "+c.bearerAuth)

	return req, nil
}

func (c *Client) doRequestReturningJson(respBody any, req *http.Request) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.Errorf("non-200 status code: %v", resp.StatusCode)
	}

	// TEMP
	if b, err := httputil.DumpResponse(resp, true); err == nil {
		log.Printf("%v", string(b))
	}
	// END

	if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
		return errors.Wrapf(err, "cannot decode response body")
	}

	return nil
}
