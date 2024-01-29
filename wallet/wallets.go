package wallet

import "errors"

type Wallets struct {
	WalletsInfo map[string]*Wallet
}

func NewWallets() *Wallets {
	return &Wallets{
		WalletsInfo: make(map[string]*Wallet),
	}
}

func (ws *Wallets) InsertWallet(name string, wallet *Wallet) {
	ws.WalletsInfo[name] = wallet
}

func (ws *Wallets) SerachWallet(name string) (*Wallet, error) {
	if ws.isExist(name) {
		return ws.WalletsInfo[name], nil
	} else {
		return nil, errors.New("没这个人的钱包")
	}
}
func (ws *Wallets) isExist(name string) bool {
	_, ok := ws.WalletsInfo[name]
	return ok
}
func (ws *Wallets) NewWallet(name string) *Wallet {
	wallet := NewWallet()
	if ws.isExist(name) {
		return ws.WalletsInfo[name]
	}
	return wallet
}
