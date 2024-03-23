package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Client struct {
	c *http.Client
}

func New() *Client {
	return &Client{
		c: http.DefaultClient,
	}
}

func (c *Client) Call(ctx context.Context, url string, request Request) (Response, error) {
	payload, err := request.JSON()
	if err != nil {
		return Response{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return Response{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return Response{}, err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Response{}, err
	}

	if response.Error != nil {
		return Response{}, response.Error
	}

	return response, nil
}
