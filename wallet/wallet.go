package wallet

import (
	"case/setting"
	"case/utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"

	"golang.org/x/crypto/ripemd160"
)

// 实现一下钱包功能
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey
	Address    []byte
}

func newKeyPair() (ecdsa.PrivateKey, ecdsa.PublicKey) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	return *privateKey, privateKey.PublicKey
}
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := &Wallet{PrivateKey: private, PublicKey: public}
	wallet.getAddr()
	return wallet
}

func (w *Wallet) getAddr() {
	// 计算公钥的哈希
	pubKeyBytes := elliptic.Marshal(w.PublicKey.Curve, w.PublicKey.X, w.PublicKey.Y)
	hash := sha256.Sum256(pubKeyBytes)
	ripemd160Hasher := ripemd160.New()
	_, err := ripemd160Hasher.Write(hash[:])
	if err != nil {
		fmt.Println("Error hashing public key:", err)
		return
	}
	publicKeyHash := ripemd160Hasher.Sum(nil)

	// 添加版本前缀
	versionedPayload := append([]byte(setting.Version), publicKeyHash...)

	// 计算两次 SHA-256 哈希
	hash = sha256.Sum256(versionedPayload)
	hash = sha256.Sum256(hash[:])

	// 计算校验码
	checksum := hash[:4]

	// 生成比特币地址
	fullPayload := append(versionedPayload, checksum...)
	address := utils.Base58Encode(fullPayload)

	w.Address = address
}
func (w *Wallet) GetAddress() []byte {
	return w.Address
}
func Case() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))
	sign := utils.Sign(privateKey, hash[:])
	valid := utils.Verify(&privateKey.PublicKey, hash[:], sign)
	fmt.Println(valid)
	wallet := NewWallet()
	fmt.Println(wallet.Address)
}
