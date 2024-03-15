package main

import (
	dscWallet "bitbucket.org/decimalteam/dsc-go-sdk/wallet"
	"cosmossdk.io/math"
	"fmt"
	"math/rand"
	// Required imports
	dscApi "bitbucket.org/decimalteam/dsc-go-sdk/api"
	dscTx "bitbucket.org/decimalteam/dsc-go-sdk/tx"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// PLEASE, DON'T USE THIS MNEMONIC OR ANY PUBLIC EXPOSED MNEMONIC IN MAINNET
	testMnemonicWords      = "repair furnace west loud peasant false six hockey poem tube now alien service phone hazard winter favorite away sand fuel describe version tragic vendor"
	testMnemonicPassphrase = ""
)

func random(min int, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	// Option 1. Generate private key (account) by mnemonic words (bip39)
	account, err := dscWallet.NewAccountFromMnemonicWords(testMnemonicWords, testMnemonicPassphrase)
	if err != nil {
		// Error handling
	}
	// Output: d01...
	fmt.Println(account.Address())

	api := dscApi.NewAPI("https://devnet-gate.decimalchain.com/api/")

	err = api.GetParameters()
	// ...error handling
	// now api has valid results for api.ChainID(), api.BaseCoin()

	// 3. Bind account
	accNumber, accSequence, _ := api.GetAccountNumberAndSequence(account.Address())
	// ...error handling
	account = account.WithChainID(api.ChainID()).WithSequence(accSequence).WithAccountNumber(accNumber)

	// 4. Create transaction message

	// or you can use message type directly
	msg := dscTx.MsgSendCoin{
		Sender:    account.Address(),
		Recipient: "d01spflnvz93jrgsseuam2c6lv2xys2e7cpn0wsky",
		Coin:      sdk.NewCoin(api.BaseCoin(), math.NewInt(1000000000000000000)),
	}

	tx, _ := dscTx.BuildTransaction(
		account,
		[]sdk.Msg{&msg},
		"",
		sdk.NewCoin(api.BaseCoin(), sdk.NewInt(0)),
	)
	// ...error handling

	// 5. Sign and send
	_ = tx.SignTransaction(account)
	// ...error handling
	bz, _ := tx.BytesToSend()
	// ...error handling

	// use one of methods:
	// 1) BroadcastTxSync: send transaction in SYNC mode and get transaction hash and
	// possible error of transaction check
	// You can check later transaction delivery by hash
	result, _ := api.BroadcastTxSync(bz)

	// only gate API
	// wait for block when using BroadcastTxSync
	// 6. Verify transaction delivery
	// NOTE: if transaction not in block already, you can get HTTP 404 error
	// If you want to be sure after every transaction, use BroadcastTxCommit
	fmt.Println(result.Hash)
}
