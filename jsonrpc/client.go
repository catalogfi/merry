package jsonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/catalogfi/merry"
)

type Client interface {
	Fund(req RequestFund) error
	Update() error
	Version() (ResponseVersion, error)
	RelayBitcoin(req RequestRelay) (interface{}, error)
	RelayEthereum(req RequestRelay) (interface{}, error)
	RelayArbitrum(req RequestRelay) (interface{}, error)
}

type client struct {
	client *http.Client
	url    string
}

func NewClient(c *http.Client, url string) Client {
	return &client{client: c, url: url}
}

func (client *client) Fund(req RequestFund) error {
	return client.submitRequest(req, nil, MethodMerryFund)
}

func (client *client) Update() error {
	return client.submitRequest(nil, nil, MethodMerryUpdate)
}

func (client *client) Version() (ResponseVersion, error) {
	var res ResponseVersion
	if err := client.submitRequest(nil, &res, MethodMerryVersion); err != nil {
		return res, err
	}
	return res, nil
}

func (client *client) RelayBitcoin(request RequestRelay) (interface{}, error) {
	var res interface{}
	if err := client.submitRequest(request, &res, MethodMerryBitcoin); err != nil {
		return res, err
	}
	return res, nil
}

func (client *client) RelayEthereum(request RequestRelay) (interface{}, error) {
	var res interface{}
	if err := client.submitRequest(request, &res, MethodMerryEthereum); err != nil {
		return res, err
	}
	return res, nil
}

func (client *client) RelayArbitrum(request RequestRelay) (interface{}, error) {
	var res interface{}
	if err := client.submitRequest(request, &res, MethodMerryArbitrum); err != nil {
		return res, err
	}
	return res, nil
}

func (client *client) submitRequest(req, res interface{}, method string) error {
	data := []byte{}
	if req != nil {
		var err error
		data, err = json.Marshal(req)
		if err != nil {
			return err
		}
	}
	request := merry.Request{
		Version: "2.0",
		ID:      rand.Int(),
		Method:  method,
		Params:  data,
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(request); err != nil {
		return err
	}

	resp, err := http.Post(client.url, "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response merry.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		if response.Error == nil {
			return fmt.Errorf("request failed : status code = %v", resp.StatusCode)
		}
		return fmt.Errorf("request failed : code = %v, err = %v", response.Error.Code, response.Error.Data)
	}

	if res == nil {
		return nil
	}
	return json.Unmarshal(response.Result, &res)
}
