package telegram

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
