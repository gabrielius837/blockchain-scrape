package blockchainscrape

import (
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

func charByteToByte(input byte) (byte, error) {
	switch input {
	case '0':
		return 0, nil
	case '1':
		return 1, nil
	case '2':
		return 2, nil
	case '3':
		return 3, nil
	case '4':
		return 4, nil
	case '5':
		return 5, nil
	case '6':
		return 6, nil
	case '7':
		return 7, nil
	case '8':
		return 8, nil
	case '9':
		return 9, nil
	case 'a', 'A':
		return 10, nil
	case 'b', 'B':
		return 11, nil
	case 'c', 'C':
		return 12, nil
	case 'd', 'D':
		return 13, nil
	case 'e', 'E':
		return 14, nil
	case 'f', 'F':
		return 15, nil
	}
	return 0, fmt.Errorf("unexpected byte: %d", input)
}

func boolToInt(input bool) int {
	if input {
		return 1
	}
	return 0
}

func HexBytesToBytes(bytes []byte) ([]byte, error) {
	bytes = stripQoutes(bytes)
	length := len(bytes)

	// check for '0x' pre suffix
	hasSuffix := length > 1 && bytes[0] == 48 && bytes[1] == 120

	if length == 0 || (length == 2 && hasSuffix) {
		return []byte{}, nil
	}

	odd := length%2 != 0

	strategy := 0 | boolToInt(hasSuffix) | boolToInt(odd)<<1
	switch strategy {
	case 0:
		break
	// remove suffix
	case 1:
		bytes = bytes[2:]
	// append '0' byte in front
	case 2:
		bytes = append(bytes, 0)
		copy(bytes[1:], bytes)
		bytes[0] = '0'
	// remove suffix and append '0' byte in front
	case 3:
		bytes = bytes[1:]
		bytes[0] = '0'
	default:
		return nil, fmt.Errorf("unexpected strategy: %d", strategy)
	}

	length = len(bytes) / 2
	output := make([]byte, length)
	for i := range output {
		k := i * 2

		first, err := charByteToByte(bytes[k])
		if err != nil {
			return nil, err
		}
		second, err := charByteToByte(bytes[k+1])
		if err != nil {
			return nil, err
		}
		output[i] = first*16 + second
	}
	return output, nil
}

func HexStringToBytes(hexString string) ([]byte, error) {
	bytes := []byte(hexString)
	result, err := HexBytesToBytes(bytes)

	return result, err
}

func bytesToUint64(bytes []byte) (uint64, error) {
	length := len(bytes)

	if length == 8 {
		return binary.BigEndian.Uint64(bytes), nil
	}

	if length > 8 {
		i := 0
		for ; i < length; i++ {
			if bytes[i] != 0 || length-i == 8 {
				break
			}
		}

		if length-i != 8 {
			return 0, fmt.Errorf("overflow prevented while converting hex string to Uint64")
		}

		bytes = bytes[i:]
	} else {
		temp := make([]byte, 8)
		offset := 8 - length
		for i := 0; i < length; i++ {
			temp[i+offset] = bytes[i]
		}
		bytes = temp
	}

	return binary.BigEndian.Uint64(bytes), nil
}

func HexBytesToUint64(bytes []byte) (uint64, error) {
	bytes, err := HexBytesToBytes(bytes)
	if err != nil {
		return 0, err
	}

	result, err := bytesToUint64(bytes)
	return result, err
}

func HexStringToUint64(hexString string) (uint64, error) {
	bytes, err := HexStringToBytes(hexString)
	if err != nil {
		return 0, err
	}

	result, err := bytesToUint64(bytes)
	return result, err
}

func stripQoutes(input []byte) []byte {
	length := len(input)
	if length > 1 && input[0] == '"' && input[length-1] == '"' {
		return input[1 : length-1]
	}

	return input
}

func stripZeroes(bytes []byte) []byte {
	length := len(bytes)
	if length == 0 {
		return bytes
	}

	index := 0
	for ; index < length; index++ {
		if bytes[index] != 0 {
			break
		}
	}

	return bytes[index:]
}

func (address Address) GetChecksumAddress() ([]byte, error) {
	inter := fmt.Sprintf("%x", address)
	bytes := []byte(inter)
	hasher := sha3.NewLegacyKeccak256()
	_, err := hasher.Write(bytes)
	if err != nil {
		return nil, err
	}

	hash := hasher.Sum(nil)

	for i := 0; i < 40; i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte >>= 4
		} else {
			hashByte &= 0xf
		}

		if bytes[i] > '9' && hashByte > 7 {
			bytes[i] -= 32
		}
	}

	return bytes, nil
}

func (transaction Transaction) ComputeContractAddres() (Address, error) {
	var contract Address
	if !transaction.To.IsNull() {
		return contract, fmt.Errorf("in order to compute new contract address need null 'to' address")
	}
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, uint64(transaction.Nonce))
	nonceBytes = stripZeroes(nonceBytes)

	bytes, err := rlp.EncodeToBytes([][]byte{transaction.From[:], nonceBytes})
	if err != nil {
		return contract, err
	}
	hasher := sha3.NewLegacyKeccak256()
	_, err = hasher.Write(bytes)
	if err != nil {
		return contract, err
	}

	result := hasher.Sum(nil)

	offset := 12
	for i := 0; i < ADDRESS_LENGTH; i++ {
		contract[i] = result[i+offset]
	}

	return contract, nil
}
