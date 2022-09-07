package client

func getTestClient() *Client {
	clt := NewClient(&ClientConfig{
		Endpoint: "https://valhalla1.openstreetmap.de",
	})

	return clt
}
