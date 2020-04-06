package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"

	"github.com/the729/go-libra/crypto"
)

func cmdCreateAccounts(ctx *cli.Context) error {
	if _, err := os.Stat(WalletFile); err == nil {
		log.Printf("wallet file (%s) already exists.", WalletFile)
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
		account := &AccountConfig{
			PrivateKey: crypto.PrivateKey(prikey),
		}
		hasher.Sum(account.Address[:0])
		wallet.Accounts = append(wallet.Accounts, account)
	}

	f, err := os.Create(WalletFile)
	if err != nil {
		log.Printf("cannot create wallet file in %s", WalletFile)
		return nil
	}
	defer f.Close()
	err = toml.NewEncoder(f).Encode(wallet)
	if err != nil {
		log.Printf("cannot encode toml file: %v", err)
	}

	return cmdListAccounts(ctx)
}

func cmdListAccounts(ctx *cli.Context) error {
	wallet, err := LoadAccounts(WalletFile)
	if err != nil {
		log.Fatal(err)
	}

	for addr, account := range wallet.Accounts {
		log.Printf("account: %s   authkey prefix: %s   prikey prefix: %s\n",
			addr,
			hex.EncodeToString(account.AuthKey[:16]),
			hex.EncodeToString(account.PrivateKey[:4]),
		)
	}

	return nil
}

func cmdMint(ctx *cli.Context) error {
	wallet, err := LoadAccounts(WalletFile)
	if err != nil {
		log.Fatal(err)
	}

	receiver, err := wallet.GetAccount(ctx.Args().Get(0))
	if err != nil {
		return err
	}

	amount, err := strconv.Atoi(ctx.Args().Get(1))
	if err != nil {
		return err
	}
	amountMicro := uint64(amount) * 1000000

	faucetURL := fmt.Sprintf("http://faucet.testnet.libra.org/?amount=%d&auth_key=%s", amountMicro, hex.EncodeToString(receiver.AuthKey))
	log.Printf("Going to POST to faucet service: %s", faucetURL)

	resp, err := http.PostForm(faucetURL, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	log.Printf("Respone (code=%d): %s", resp.StatusCode, string(body))

	return nil
}
