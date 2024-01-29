package utils

import (
	"bytes"
	"case/setting"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"os"
)

func GetDataHash(data []byte) []byte {
	var hash [32]byte
	hash = sha256.Sum256(data)
	return hash[:]
}

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DbExists() bool {
	if _, err := os.Stat(setting.DbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Base58Encode 对数据进行 Base58 编码
func Base58Encode(input []byte) []byte {

	var result []byte

	x := new(big.Int).SetBytes(input)
	base := big.NewInt(int64(len(alphabet)))

	zero := big.NewInt(0)
	mod := &big.Int{}
	for x.Cmp(zero) > 0 {
		x.DivMod(x, base, mod)
		result = append(result, alphabet[mod.Int64()])
	}

	// 反转结果
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	// 处理前导零
	for _, b := range input {
		if b != 0 {
			break
		}
		result = append([]byte{alphabet[0]}, result...)
	}

	return result
}

// Base58Decode 对数据进行 Base58 解码
func Base58Decode(input []byte) []byte {
	// 创建一个字母到索引的映射
	indexMap := make(map[byte]int)
	for i, char := range alphabet {
		indexMap[byte(char)] = i
	}

	// 创建一个大整数来存储解码结果
	decoded := big.NewInt(0)

	// 对每个输入字符进行解码
	for _, char := range input {
		// 查找字符在字母表中的索引
		idx, ok := indexMap[char]
		if !ok {
			fmt.Println("Invalid character in input:", char)
			return nil
		}
		// 将索引乘以58并加到decoded中
		decoded.Mul(decoded, big.NewInt(58))
		decoded.Add(decoded, big.NewInt(int64(idx)))
	}

	// 将大整数转换为字节数组
	decodedBytes := decoded.Bytes()

	// 处理前导零
	for _, char := range input {
		if char == alphabet[0] {
			decodedBytes = append([]byte{0x00}, decodedBytes...)
		} else {
			break
		}
	}

	return decodedBytes
}
func Sign(privateKey *ecdsa.PrivateKey, hash []byte) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		fmt.Println("Error signing data:", err)
		return nil
	}
	// 签名插入'-'进行分割 r 和 s
	signature := append(r.Bytes(), []byte("-")...)
	signature = append(signature, s.Bytes()...)
	return signature
}
func Verify(publicKey *ecdsa.PublicKey, hash []byte, signature []byte) bool {
	index := bytes.Index(signature, []byte{'-'})
	if index == -1 {
		return false
	}
	r, s := signature[:index], signature[index+1:]
	bigr := new(big.Int).SetBytes(r)
	bigs := new(big.Int).SetBytes(s)
	isValid := ecdsa.Verify(publicKey, hash[:], bigr, bigs)
	return isValid
}
