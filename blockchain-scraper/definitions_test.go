package blockchainscrape

import (
	"testing"
)

func TestAddressMapping(t *testing.T) {
	//input := []byte("{\"jsonrpc\":\"2.0\",\"id\":0,\"result\":{\"number\":\"0xf4240\",\"hash\":\"0x8e38b4dbf6b11fcc3b9dee84fb7986e29ca0a02cecd8977c161ff7333329681e\",\"transactions\":[{\"blockHash\":\"0x8e38b4dbf6b11fcc3b9dee84fb7986e29ca0a02cecd8977c161ff7333329681e\",\"blockNumber\":\"0xf4240\",\"hash\":\"0xea1093d492a1dcb1bef708f771a99a96ff05dcab81ca76c31940300177fcf49f\",\"chainId\":\"0x0\",\"from\":\"0x39fa8c5f2793459d6622857e7d9fbb4bd91766d3\",\"gas\":\"0x1f8dc\",\"gasPrice\":\"0x12bfb19e60\",\"input\":\"0x\",\"nonce\":\"0x15\",\"r\":\"0xa254fe085f721c2abe00a2cd244110bfc0df5f4f25461c85d8ab75ebac11eb10\",\"s\":\"0x30b7835ba481955b20193a703ebc5fdffeab081d63117199040cdf5a91c68765\",\"to\":\"0xc083e9947cf02b8ffc7d3090ae9aea72df98fd47\",\"transactionIndex\":\"0x0\",\"type\":\"0x0\",\"v\":\"0x1c\",\"value\":\"0x56bc75e2d63100000\"}],\"difficulty\":\"0xb69de81a22b\",\"extraData\":\"0xd783010303844765746887676f312e352e31856c696e7578\",\"gasLimit\":\"0x2fefd8\",\"gasUsed\":\"0xc444\",\"logsBloom\":\"0x00000000000000000000000000000000000800000000000000000000000800000000000000000400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000\",\"miner\":\"0x2a65aca4d5fc5b5c859090a6c34d164135398226\",\"mixHash\":\"0x92c4129a0ae2361b452a9edeece55c12eceeab866316195e3d87fc1b005b6645\",\"nonce\":\"0xcd4c55b941cf9015\",\"parentHash\":\"0xb4fbadf8ea452b139718e2700dc1135cfc81145031c84b7ab27cd710394f7b38\",\"receiptsRoot\":\"0x20e3534540caf16378e6e86a2bf1236d9f876d3218fbc03958e6db1c634b2333\",\"sha3Uncles\":\"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347\",\"size\":\"0x300\",\"stateRoot\":\"0x0e066f3c2297a5cb300593052617d1bca5946f0caa0635fdb1b85ac7e5236f34\",\"timestamp\":\"0x56bfb415\",\"totalDifficulty\":\"0x630554d65f2cfe6a\",\"transactionsRoot\":\"0x65ba887fcb0826f616d01f736c1d2d677bcabde2f7fc25aa91cfbc0b3bad5cb3\",\"uncles\":[]}}")
	input := []byte("\"0x39fa8c5f2793459d6622857e7d9fbb4bd91766d3\"")
	var result Address
	expected := []byte{57, 250, 140, 95, 39, 147, 69, 157, 102, 34, 133, 126, 125, 159, 187, 75, 217, 23, 102, 211}
	err := result.UnmarshalJSON(input)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < ADDRESS_LENGTH; i++ {
		if result[i] != expected[i] {
			t.Fatalf("mismatch at index %d while comparing expected and result bytes\n%v\n%v", i, expected, result)
		}
	}
}

func TestHashMapping(t *testing.T) {
	input := []byte("\"0xea1093d492a1dcb1bef708f771a99a96ff05dcab81ca76c31940300177fcf49f\"")
	var result Hash
	expected := []byte{234, 16, 147, 212, 146, 161, 220, 177, 190, 247, 8, 247, 113, 169, 154, 150, 255, 5, 220, 171, 129, 202, 118, 195, 25, 64, 48, 1, 119, 252, 244, 159}
	err := result.UnmarshalJSON(input)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < HASH_LENGTH; i++ {
		if result[i] != expected[i] {
			t.Fatalf("mismatch at index %d while comparing expected and result bytes\n%v\n%v", i, expected, result)
		}
	}
}

func TestNullAddressMapping_WhenAddress(t *testing.T) {
	input := []byte("\"0x39fa8c5f2793459d6622857e7d9fbb4bd91766d3\"")
	var result NullAddress
	expected := []byte{57, 250, 140, 95, 39, 147, 69, 157, 102, 34, 133, 126, 125, 159, 187, 75, 217, 23, 102, 211}
	err := result.UnmarshalJSON(input)
	if err != nil {
		t.Fatal(err)
	}

	isNull := result.IsNull()

	if isNull {
		t.Fatalf("expected to have not null address, instead got:\n%v", result)
	}

	for i := 0; i < ADDRESS_LENGTH; i++ {
		if result[i] != expected[i] {
			t.Fatalf("mismatch at index %d while comparing expected and result bytes\n%v\n%v", i, expected, result)
		}
	}
}

func TestNullAddressMapping_WhenNull(t *testing.T) {
	input := []byte("\"0x\"")
	var result NullAddress
	err := result.UnmarshalJSON(input)
	if err != nil {
		t.Fatal(err)
	}

	isNull := result.IsNull()

	if !isNull {
		t.Fatalf("expected to have null address, instead got:\n%v", result)
	}

	length := len(result)
	if length != ZERO_LENGTH {
		t.Fatalf("expected to have result wit length of %d, instead got %d", ZERO_LENGTH, length)
	}
}

func TestHexUint64Mapping(t *testing.T) {
	input := []byte("\"0x1f8dc\"")
	expected := HexUint64(129244)
	var result HexUint64
	err := result.UnmarshalJSON(input)
	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatalf("expected and result values are different, %d != %d", expected, result)
	}
}
