package client

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"

	lcrypto "github.com/the729/go-libra/crypto"
	"github.com/the729/go-libra/types"
)

type Account struct {
	PrivateKey     lcrypto.PrivateKey   `toml:"private_key"`
	Address        types.AccountAddress `toml:"address"`
	SequenceNumber uint64               `toml:"-"`
}

type WalletConfig struct {
	Accounts []*Account `toml:"accounts"`
}

func (c *Client) LoadAccounts() error {
	wallet := &WalletConfig{}
	_, err := toml.DecodeFile(c.WalletFile, wallet)
	if err != nil {
		return fmt.Errorf("toml decode file error: %v", err)
	}

	c.accounts = make(map[string]*Account)
	for _, account := range wallet.Accounts {
		if account.PrivateKey != nil {
			pubkey := ed25519.PrivateKey(account.PrivateKey).Public().(ed25519.PublicKey)
			hasher := sha3.New256()
			hasher.Write(pubkey)
			account.Address = hasher.Sum([]byte{})
		}
		c.accounts[hex.EncodeToString(account.Address)] = account
	}
	return nil
}

func (c *Client) GetAccount(prefix string) (*Account, error) {
	var seen *Account
	for addr, account := range c.accounts {
		if strings.HasPrefix(addr, prefix) {
			if seen != nil {
				return nil, fmt.Errorf("more than 1 accounts have prefix %s", prefix)
			}
			seen = account
		}
	}
	if seen != nil {
		return seen, nil
	}

	newAddr, err := hex.DecodeString(prefix)
	if err == nil && len(newAddr) == types.AccountAddressLength {
		return &Account{
			Address: newAddr,
		}, nil
	}

	return nil, fmt.Errorf("account not present in local config file")
}

func (c *Client) CmdCreateAccounts(ctx *cli.Context) error {
	if _, err := os.Stat(c.WalletFile); err == nil {
		log.Printf("wallet file (%s) already exists.", c.WalletFile)
		return nil
	}

	number, _ := strconv.Atoi(ctx.Args().Get(0))
	if number == 0 {
		number = 10
	}
	log.Printf("generating %d accounts...", number)
	wallet := &WalletConfig{}
	for i := 0; i < number; i++ {
		pubkey, prikey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}
		hasher := sha3.New256()
		hasher.Write(pubkey)
		account := &Account{
			PrivateKey: lcrypto.PrivateKey(prikey),
			Address:    hasher.Sum([]byte{}),
		}
		wallet.Accounts = append(wallet.Accounts, account)
	}

	f, err := os.Create(c.WalletFile)
	if err != nil {
		log.Printf("cannot create wallet file in %s", c.WalletFile)
		return nil
	}
	defer f.Close()
	err = toml.NewEncoder(f).Encode(wallet)
	if err != nil {
		log.Printf("cannot encode toml file: %v", err)
	}

	return c.CmdListAccounts(ctx)
}

func (c *Client) CmdListAccounts(ctx *cli.Context) error {
	if err := c.LoadAccounts(); err != nil {
		log.Printf("cannot load accounts: %v\n", err)
		return nil
	}

	for addr := range c.accounts {
		log.Printf("account: %s\n", addr)
	}

	return nil
}

func (c *Client) CmdMint(ctx *cli.Context) error {
	c.LoadAccounts()

	receiver, err := c.GetAccount(ctx.Args().Get(0))
	if err != nil {
		return err
	}

	amount, err := strconv.Atoi(ctx.Args().Get(1))
	if err != nil {
		return err
	}
	amountMicro := uint64(amount) * 1000000

	log.Printf("Please visit the following faucet service:")
	log.Printf("http://faucet.testnet.libra.org/?amount=%d&address=%s", amountMicro, hex.EncodeToString(receiver.Address))
	return nil
}
