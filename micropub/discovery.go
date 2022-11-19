package micropub

func (c *Client) Config() (MicroPubConfig, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.discoveredConfig.MediaEndpoint == "" {
		m, err := c.getDiscoverPostConfig()
		if err != nil {
			return MicroPubConfig{}, err
		}
		c.discoveredConfig = m
	}
	return c.discoveredConfig, nil
}

func (c *Client) getDiscoverPostConfig() (pc MicroPubConfig, err error) {
	req, err := c.newReq("GET", c.micropubURL+"?q=config", nil)
	if err != nil {
		return MicroPubConfig{}, err
	}

	err = c.doRequestReturningJson(&pc, req)
	return pc, err
}

type MicroPubConfig struct {
	MediaEndpoint string        `json:"media-endpoint"`
	Destination   []Destination `json:"destination"`
}

type Destination struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}
