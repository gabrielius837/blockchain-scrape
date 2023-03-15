package blockchainscrape

import (
	"fmt"
	"testing"
)

func TestHexToBytes_WhenDefault(t *testing.T) {
	hex := "0x00112233445566778899aaAAbbBBccCCddDDeeEEffFF"
	expected := []byte{0, 17, 34, 51, 68, 85, 102, 119, 136, 153, 170, 170, 187, 187, 204, 204, 221, 221, 238, 238, 255, 255}

	result, err := HexStringToBytes(hex)

	if err != nil {
		t.Fatal(err)
	}
	if len(expected) != len(result) {
		t.Fatalf("expected lengths to match, instead expected: %d result: %d", len(expected), len(result))
	}

	for i := 0; i < len(result); i++ {
		expectedByte := expected[i]
		resultByte := result[i]
		if expectedByte != resultByte {
			t.Fatalf("expected and result bytes are not matching at index %d\n%v\n%v", i, expected, result)
		}
	}
}

func TestHexToBytes_WhenOdd(t *testing.T) {
	hex := "0x0112233445566778899aaAAbbBBccCCddDDeeEEffFF"
	expected := []byte{0, 17, 34, 51, 68, 85, 102, 119, 136, 153, 170, 170, 187, 187, 204, 204, 221, 221, 238, 238, 255, 255}

	result, err := HexStringToBytes(hex)

	if err != nil {
		t.Fatal(err)
	}
	if len(expected) != len(result) {
		t.Fatalf("expected lengths to match, instead expected: %d result: %d", len(expected), len(result))
	}

	for i := 0; i < len(result); i++ {
		expectedByte := expected[i]
		resultByte := result[i]
		if expectedByte != resultByte {
			t.Fatalf("expected and result bytes are not matching at index %d\n%v\n%v", i, expected, result)
		}
	}
}

func TestHexToBytes_WhenQoutes(t *testing.T) {
	hex := "\"0x00112233445566778899aaAAbbBBccCCddDDeeEEffFF\""
	expected := []byte{0, 17, 34, 51, 68, 85, 102, 119, 136, 153, 170, 170, 187, 187, 204, 204, 221, 221, 238, 238, 255, 255}

	result, err := HexStringToBytes(hex)

	if err != nil {
		t.Fatal(err)
	}
	if len(expected) != len(result) {
		t.Fatalf("expected lengths to match, instead expected: %d result: %d", len(expected), len(result))
	}

	for i := 0; i < len(result); i++ {
		expectedByte := expected[i]
		resultByte := result[i]
		if expectedByte != resultByte {
			t.Fatalf("expected and result bytes are not matching at index %d\n%v\n%v", i, expected, result)
		}
	}
}

func TestHexToBytes_WhenQuotesAndOdd(t *testing.T) {
	hex := "\"0x0112233445566778899aaAAbbBBccCCddDDeeEEffFF\""
	expected := []byte{0, 17, 34, 51, 68, 85, 102, 119, 136, 153, 170, 170, 187, 187, 204, 204, 221, 221, 238, 238, 255, 255}

	result, err := HexStringToBytes(hex)

	if err != nil {
		t.Fatal(err)
	}
	if len(expected) != len(result) {
		t.Fatalf("expected lengths to match, instead expected: %d result: %d", len(expected), len(result))
	}

	for i := 0; i < len(result); i++ {
		expectedByte := expected[i]
		resultByte := result[i]
		if expectedByte != resultByte {
			t.Fatalf("expected and result bytes are not matching at index %d\n%v\n%v", i, expected, result)
		}
	}
}

func TestHexToUint64_WhenDefault(t *testing.T) {
	hex := "0x1122334455667788"
	expected := uint64(1234605616436508552)

	result, err := HexStringToUint64(hex)
	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatalf("expected and result are different, %d != %d", expected, result)
	}
}

func TestHexToUint64_WithTruncated(t *testing.T) {
	hex := "0x300"
	expected := uint64(768)

	result, err := HexStringToUint64(hex)
	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatalf("expected and result are different, %d != %d", expected, result)
	}
}

func TestHexToUint64_WithExcessZeroesAndTruncated(t *testing.T) {
	hex := "0x000000000000000300"
	expected := uint64(768)

	result, err := HexStringToUint64(hex)
	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatalf("expected and result are different, %d != %d", expected, result)
	}
}

func TestHexToUint64_WithExcessZeroes(t *testing.T) {
	hex := "0x000000001122334455667788"
	expected := uint64(1234605616436508552)

	result, err := HexStringToUint64(hex)
	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatalf("expected and result are different, %d != %d", expected, result)
	}
}

func TestHexToUint64_WithExcessZeroesAndBytes(t *testing.T) {
	hex := "0x000000aa1122334455667788"

	_, err := HexStringToUint64(hex)
	if err == nil {
		t.Fatal("expected to fail")
	}
}

func TestHexToUint64_WithOverflow(t *testing.T) {
	hex := "0xff1122334455667788"

	result, err := HexStringToUint64(hex)
	if err == nil {
		t.Fatalf("expected to fail but got %d", result)
	}
}

func TestGetChecksumAddress(t *testing.T) {
	address := Address([ADDRESS_LENGTH]byte{76, 231, 243, 100, 108, 201, 130, 246, 70, 157, 91, 36, 62, 221, 54, 233, 175, 137, 192, 240})
	expected := "0x4cE7f3646cc982F6469D5B243edd36E9aF89c0F0"

	checksum, err := address.GetChecksumAddress()
	if err != nil {
		t.Fatal(err)
	}

	result := fmt.Sprintf("0x%s", checksum)

	if expected != result {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestComputeContractAddress(t *testing.T) {
	address := Address([ADDRESS_LENGTH]byte{76, 231, 243, 100, 108, 201, 130, 246, 70, 157, 91, 36, 62, 221, 54, 233, 175, 137, 192, 240})
	nonce := HexUint64(65278)
	expected := Address([ADDRESS_LENGTH]byte{120, 26, 88, 71, 4, 207, 91, 28, 117, 84, 141, 172, 24, 88, 112, 40, 206, 160, 110, 191})
	tx := Transaction{From: address, Nonce: nonce}

	result, err := tx.ComputeContractAddres()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < ADDRESS_LENGTH; i++ {
		expectedByte := expected[i]
		resultByte := result[i]
		if expectedByte != resultByte {
			t.Fatalf("expected and result bytes are not matching at index %d\n%v\n%v", i, expected, result)
		}
	}
}
