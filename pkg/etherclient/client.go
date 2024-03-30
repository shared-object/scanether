package etherclient

import (
	"bytes"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"strconv"

	dto "github.com/xllwhoami/ethergrab/pkg/dto"
)

type Client struct {
	Endpoint string
}

func NewClient(endpoint_url string) *Client {
	return &Client{Endpoint: endpoint_url}
}

func (c *Client) GetLatestBlockNumber() (string, error) {
	request, err := json.Marshal(dto.JSONRPCRequestDTO{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  make([]any, 0),
		ID:      1,
	})

	if err != nil {
		return "", err
	}

	response, err := http.Post(c.Endpoint, "application/json", bytes.NewBuffer(request))

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	result := new(dto.ETHBlockNumberResponseDTO)

	json.Unmarshal(body, result)

	return result.Result, nil

}

func (c *Client) GetBlockByNumber(blockno string) (dto.ETHGetBlockByBumberResponseDTO, error) {
	request, err := json.Marshal(dto.JSONRPCRequestDTO{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []any{blockno, true},
		ID:      1,
	})

	if err != nil {
		return dto.ETHGetBlockByBumberResponseDTO{}, err
	}

	response, err := http.Post(c.Endpoint, "application/json", bytes.NewBuffer(request))

	body, err := io.ReadAll(response.Body)

	print(string(body))

	if err != nil {
		return dto.ETHGetBlockByBumberResponseDTO{}, err

	}

	result := new(dto.ETHGetBlockByBumberResponseDTO)

	json.Unmarshal(body, result)

	return *result, nil

}

func (c *Client) NumberFromHex(hexString string) int64 {
	number := new(big.Int)
	number.SetString(hexString, 0)

	return number.Int64()
}

func (c *Client) NumberToHex(number int64) string {
	return "0x" + strconv.FormatInt(int64(number), 16)
}
