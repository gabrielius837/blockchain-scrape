package blockchainscrape

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	URL = "https://eth-mainnet.g.alchemy.com/v2"
)

func InitRequest(apiKey string, payload io.Reader) ([]byte, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", URL, apiKey), payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetBlock(apiKey string, index uint64, number uint64) (*BlockResponse, error) {
	if number >= 4370000 {
		fmt.Fprintln(os.Stderr, "byzantium reached update program")
		os.Exit(1)
	}
	body := fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\"0x%x\",true],\"id\":%d}", number, index)
	payload := strings.NewReader(body)
	bytes, err := InitRequest(apiKey, payload)
	if err != nil {
		return nil, err
	}

	resp := &BlockResponse{}
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetCode(apiKey string, index uint64, address []byte) (*CodeResponse, error) {
	length := len(address)
	if length != 20 {
		return nil, fmt.Errorf("unexpected address, expected 20, got %d", length)
	}
	body := fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"method\":\"eth_getCode\",\"params\":[\"0x%x\",\"latest\"],\"id\":%d}", address, index)
	payload := strings.NewReader(body)
	bytes, err := InitRequest(apiKey, payload)
	if err != nil {
		return nil, err
	}

	resp := &CodeResponse{}
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
