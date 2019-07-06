package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	pbtypes "github.com/the729/go-libra/generated/types"
	"github.com/the729/go-libra/types"
	"github.com/urfave/cli"
)

func (c *Client) QueryAccountState(addr types.AccountAddress) (*types.AccountStateWithProof, *types.LedgerInfo, error) {
	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.ac.UpdateToLatestLedger(ctx1, &pbtypes.UpdateToLatestLedgerRequest{
		ClientKnownVersion: 0,
		RequestedItems: []*pbtypes.RequestItem{
			&pbtypes.RequestItem{
				RequestedItems: &pbtypes.RequestItem_GetAccountStateRequest{
					GetAccountStateRequest: &pbtypes.GetAccountStateRequest{
						Address: addr,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("rpc failed: %v", err)
	}
	// respj, _ := json.MarshalIndent(resp, "", "    ")
	// log.Println(string(respj))

	li := &types.LedgerInfoWithSignatures{}
	li.FromProto(resp.LedgerInfoWithSigs)
	err = li.Verify(c.verifier)
	if err != nil {
		return nil, nil, fmt.Errorf("verify failed: %v", err)
	}

	account := &types.AccountStateWithProof{}
	err = account.FromProtoResponse(resp.ResponseItems[0].GetGetAccountStateResponse())
	if err != nil {
		return nil, nil, fmt.Errorf("account state with proof from proto failed: %v", err)
	}

	err = account.Verify(addr, li.LedgerInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("account state with proof verify failed: %v", err)
	}

	if account.Blob != nil {
		err = account.Blob.ParseToMap()
		if err != nil {
			return nil, nil, fmt.Errorf("account blob cannot parse to map: %v", err)
		}
	}
	return account, li.LedgerInfo, nil
}

func (c *Client) GetLibraCoinResourceFromAccountBlob(blob *types.AccountBlob) (*types.AccountResource, error) {
	if blob.Map == nil {
		err := blob.ParseToMap()
		if err != nil {
			return nil, fmt.Errorf("account blob cannot parse to map: %v", err)
		}
	}
	res, err := blob.GetResource(&types.StructTag{
		Address: make([]byte, 32),
		Module:  "LibraAccount",
		Name:    "T",
	})
	if err != nil {
		return nil, fmt.Errorf("get resource failed: %v", err)
	}
	return res, nil
}

func (c *Client) GetAccountSequenceNumber(addr types.AccountAddress) (uint64, *types.LedgerInfo, error) {
	state, ledgerInfo, err := c.QueryAccountState(addr)
	if err != nil {
		return 0, nil, err
	}
	if state.Blob == nil {
		return 0, ledgerInfo, errors.New("sender account not present in ledger.")
	}
	resource, err := c.GetLibraCoinResourceFromAccountBlob(state.Blob)
	if err != nil {
		return 0, ledgerInfo, err
	}
	return resource.SequenceNumber, ledgerInfo, nil
}

func (c *Client) CmdQueryAccountState(ctx *cli.Context) error {
	c.Connect()
	defer c.Disconnect()
	c.LoadTrustedPeers()
	c.LoadAccounts()

	addr := ctx.Args().Get(0)
	account, err := c.GetAccount(addr)
	if err != nil {
		return err
	}
	accountState, ledgerInfo, err := c.QueryAccountState(account.Address)
	if err != nil {
		return err
	}

	log.Printf("Ledger info: version %d, time %d", ledgerInfo.Version, ledgerInfo.TimestampUsec)

	if accountState.Blob != nil {
		log.Printf("Account version: %d", accountState.Version)
		log.Printf("Libra coin resource:")
		libraCoin, err := c.GetLibraCoinResourceFromAccountBlob(accountState.Blob)
		if err != nil {
			fmt.Println("Account does not contain libra coin resource.")
			return nil
		}
		spew.Dump(libraCoin)
	} else {
		fmt.Println("Account is not present in the ledger.")
	}
	return nil
}
