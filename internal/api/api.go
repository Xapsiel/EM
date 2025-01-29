package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Xapsiel/EffectiveMobile/internal/config"
	"github.com/Xapsiel/EffectiveMobile/internal/model"
)

type Client struct {
	Domain string
	Client *http.Client
}

func NewClient(cfg config.APIConfig) *Client {
	return &Client{
		Domain: cfg.Domain,
		Client: http.DefaultClient,
	}
}

func (c *Client) GetInfo(group, song string) (model.Song, error) {
	params := url.Values{}
	params.Set("group", group)
	params.Set("song", song)
	req, err := http.NewRequest(http.MethodGet, c.Domain+"/info?"+params.Encode(), nil)
	if err != nil {
		return model.Song{}, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return model.Song{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return model.Song{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var songResponse model.Song
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Song{}, err
	}
	if err := json.Unmarshal(body, &songResponse); err != nil {
		return model.Song{}, err
	}
	return songResponse, nil
}
