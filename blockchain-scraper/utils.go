package blockchainscrape

import (
	"encoding/binary"
	"fmt"

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

func stripZeroes(input []byte) []byte {
	length := len(input)
	if length == 0 {
		return input
	}

	count := 0

	for i := 0; i < length; i++ {
		if input[i] > 0 {
			break
		}
		count++
	}

	switch count {
	case 0:
		return input
	case length:
		return nil
	default:
		return input[count:]
	}
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

	slice := transaction.From[:]
	bytes := RlpEncodeArray([][]byte{slice, nonceBytes})
	fmt.Printf("0x%x\n", bytes)
	hasher := sha3.NewLegacyKeccak256()
	_, err := hasher.Write(bytes)
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

// https://ethereum.org/en/developers/docs/data-structures-and-encoding/rlp/
func RlpEncodeArray(input [][]byte) []byte {
	length := len(input)

	for i := 0; i < length; i++ {
		fmt.Printf("0x%x\n", input[i])
	}
	switch length {
	case 0:
		return []byte{0xc0}
	case 1:
		return RlpEncodeBytes(input[0])
	default:
		sum := 0
		for i := 0; i < length; i++ {
			result := RlpEncodeBytes(input[i])
			sum += len(result)
			input[i] = result
		}

		if sum < 56 {
			result := make([]byte, sum+1)
			result[0] = byte(0xc0 + sum)
			offset := 1
			for i := 0; i < length; i++ {
				copy(result[offset:], input[i])
				offset += len(input[i])
			}
			return result
		} else {
			lengthRepresentation := make([]byte, 4)
			binary.BigEndian.PutUint32(lengthRepresentation, uint32(length))
			lengthRepresentation = stripZeroes(lengthRepresentation)
			lengthOfRepresentation := len(lengthRepresentation)
			result := make([]byte, 1+lengthOfRepresentation+length)
			result[0] = byte(0xf7 + lengthOfRepresentation)
			copy(result[1:], lengthRepresentation)
			offset := 1 + lengthOfRepresentation
			for i := 1; i < length; i++ {
				copy(result[offset:], input[i])
				offset += len(input[i])
			}
			return result
		}
		//return bytes.Join(input, []byte{})
	}
}

func RlpEncodeBytes(input []byte) []byte {
	length := len(input)

	switch length {
	case 0:
		return []byte{0x80}
	case 1:
		value := input[0]
		if value < 128 {
			return input
		} else {
			return []byte{0x81, value}
		}
	default:
		if length < 56 {
			result := make([]byte, length+1)
			result[0] = 0x80 + byte(length)
			copy(result[1:], input)
			return result
		} else {
			lengthRepresentation := make([]byte, 4)
			binary.BigEndian.PutUint32(lengthRepresentation, uint32(length))
			lengthRepresentation = stripZeroes(lengthRepresentation)
			lengthOfRepresentation := len(lengthRepresentation)
			result := make([]byte, 1+lengthOfRepresentation+length)
			result[0] = byte(0xb7 + lengthOfRepresentation)
			copy(result[1:], lengthRepresentation)
			copy(result[(1+lengthOfRepresentation):], input)
			return result
		}
	}
}
