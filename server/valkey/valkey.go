package valkey

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
)

type ValkeyClient struct {
	Client valkey.Client
}

func NewValkeyClient() *ValkeyClient {
	return &ValkeyClient{}
}

func (c *ValkeyClient) Initialize(connectionString string) error {
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{connectionString}})
	if err != nil {
		return err
	}

	resp := client.Do(context.Background(), client.B().Ping().Build())
	if resp.Error() != nil {
		return resp.Error()
	}

	c.Client = client
	fmt.Println("Valkey client initialized successfully")
	return nil
}

func (c *ValkeyClient) SetValue(key string, value string) error {
	resp := c.Client.Do(context.Background(), c.Client.B().Set().Key(key).Value(value).Nx().Build())
	return resp.Error()
}

func (c *ValkeyClient) GetValue(key string) (string, error) {
	resp := c.Client.Do(context.Background(), c.Client.B().Get().Key(key).Build())
	if resp.Error() != nil {
		return "", resp.Error()
	}
	return resp.String(), nil
}
