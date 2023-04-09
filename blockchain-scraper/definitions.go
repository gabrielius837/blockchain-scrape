package blockchainscrape

import "fmt"

type HexUint64 uint64
type Address [ADDRESS_LENGTH]byte
type NullAddress []byte
type Hash [HASH_LENGTH]byte

const (
	ADDRESS_LENGTH = 20
	HASH_LENGTH    = 32
	UINT64_LENGTH  = 8
	ZERO_LENGTH    = 0
)

func convertAndCheckLength(input []byte, expectedLength int) ([]byte, error) {
	bytes, err := HexBytesToBytes(input)
	if err != nil {
		return nil, err
	}

	length := len(bytes)
	if expectedLength != length {
		return nil, fmt.Errorf("unexpected length, expected %d, got %d", expectedLength, length)
	}

	return bytes, nil
}

func (address *Address) UnmarshalJSON(input []byte) error {
	bytes, err := convertAndCheckLength(input, ADDRESS_LENGTH)
	if err != nil {
		return err
	}

	for i := 0; i < ADDRESS_LENGTH; i++ {
		address[i] = bytes[i]
	}
	return nil
}

func (hash *Hash) UnmarshalJSON(input []byte) error {
	bytes, err := convertAndCheckLength(input, HASH_LENGTH)
	if err != nil {
		return err
	}

	for i := 0; i < HASH_LENGTH; i++ {
		hash[i] = bytes[i]
	}
	return nil
}

func (nullAddress *NullAddress) UnmarshalJSON(input []byte) error {
	length := len(input)
	// check for "0x" or 0x
	if length == 4 && input[0] == '"' && input[1] == '0' && input[2] == 'x' && input[3] == '"' ||
		length == 2 && input[0] == '0' && input[1] == 'x' {
		return nil
	}

	bytes, err := convertAndCheckLength(input, ADDRESS_LENGTH)
	if err != nil {
		return err
	}

	result := NullAddress(make([]byte, ADDRESS_LENGTH))
	for i := 0; i < ADDRESS_LENGTH; i++ {
		result[i] = bytes[i]
	}
	*nullAddress = result
	return nil
}

func (hexUint *HexUint64) UnmarshalJSON(input []byte) error {
	number, err := HexBytesToUint64(input)
	if err != nil {
		return err
	}

	*hexUint = HexUint64(number)
	return nil
}

func (nullAddress NullAddress) IsNull() bool {
	length := len(nullAddress)
	return length == ZERO_LENGTH
}

type Block struct {
	Number       HexUint64     `json:"number"`
	Miner        Address       `json:"miner"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	From  Address     `json:"from"`
	To    NullAddress `json:"to"`
	Hash  Hash        `json:"hash"`
	Nonce HexUint64   `json:"nonce"`
}

type BlockResponse struct {
	Id     uint64 `json:"id"`
	Result Block  `json:"result"`
}

type CodeResponse struct {
	Id     uint64 `json:"id"`
	Result string `json:"result"`
}
