package telegram

import "fmt"

type Client struct {
	token string
}

func New() (*Client, error) {
	c := &Client{}

	t, err := mustToken()
	if err != nil {
		return nil, err
	}

	c.token = t

	return c, nil
}

func (c *Client) PrintToken() {
	fmt.Println(c.token)
}
