package micropub

import (
	"bytes"
	"emperror.dev/errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

func (c *Client) Upload(dest string, r io.Reader, filename string) (ur UploadResponse, err error) {
	cfg, err := c.Config()
	if err != nil {
		return UploadResponse{}, err
	}

	bfr, contentType, err := c.prepUploadBody(dest, r, filename)
	if err != nil {
		return UploadResponse{}, err
	}

	req, err := http.NewRequest("POST", cfg.MediaEndpoint, bfr)
	if err != nil {
		return UploadResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.bearerAuth)
	req.Header.Set("Content-Type", contentType)

	err = c.doRequestReturningJson(&ur, req)
	return ur, err
}

func (c *Client) prepUploadBody(dest string, r io.Reader, filename string) (*bytes.Buffer, string, error) {
	var b bytes.Buffer

	fw := multipart.NewWriter(&b)
	defer fw.Close()

	if err := fw.WriteField("h", "entry"); err != nil {
		return nil, "", errors.Wrap(err, "cannot write form 'h'")
	}

	// https://indieweb.org/Micropub-extensions#Destination
	if err := fw.WriteField("mp-destination", dest); err != nil {
		return nil, "", errors.Wrap(err, "cannot write form 'mp-destination'")
	}

	hdr := textproto.MIMEHeader{}
	hdr.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%v"`, filename))
	hdr.Set("Content-Type", "image/jpeg") // !! TEMP
	hdr.Set("Content-Transfer-Encoding", "binary")

	pw, err := fw.CreatePart(hdr)
	if err != nil {
		return nil, "", err
	}
	if _, err := io.Copy(pw, r); err != nil {
		return nil, "", errors.Wrap(err, "cannot write photo payload")
	}
	return &b, fw.FormDataContentType(), nil
}

type UploadResponse struct {
	URL    string `json:"url"`
	Poster string `json:"poster"`
}
