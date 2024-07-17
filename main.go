package main

import (
	"context"
	"encoding/base64"
	"fmt"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log/slog"
	"time"
)

var (
	client = resty.New().SetTimeout(30 * time.Second)

	headers = map[string]string{
		"Origin": "https://jup.ag",
	}

	rpcClient = rpc.NewWithHeaders("https://jupiter-fe.helius-rpc.com/", headers)
)

func main() {
	wallet, _ := solana.WalletFromPrivateKeyBase58("这里填入你的私钥")
	fmt.Println(wallet.PublicKey().String())

	go fuck(wallet)
	go fuck(wallet)
	go fuck(wallet)
	go fuck(wallet)

	select {}
}

func fuck(wallet *solana.Wallet) {
	for {
		txRaw := GetTx(wallet.PublicKey().String())
		if txRaw == "" {
			continue
		}

		SendTx(wallet, txRaw)
	}
}

func GetTx(address string) string {
	payload := map[string]interface{}{
		"account": address,
	}
	resp, err := client.R().SetBody(payload).Post("https://proxy.dial.to/?url=https%3A%2F%2Fsolanasummer.click%2Fon%2Fmint")
	if err != nil {
		slog.Error("request blink error: ", slog.String("error", err.Error()))
		return ""
	}

	return gjson.Get(resp.String(), "transaction").String()
}

func SendTx(wallet *solana.Wallet, rawtx string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	txData, err := base64.StdEncoding.DecodeString(rawtx)
	if err != nil {
		slog.Error("decode tx error: ", slog.String("error", err.Error()))
		return
	}

	tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(txData))
	if err != nil {
		slog.Error("decode tx error: ", slog.String("error", err.Error()))
		return
	}

	_, err = tx.PartialSign(func(key solana.PublicKey) *solana.PrivateKey {
		if wallet.PublicKey().Equals(key) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		slog.Error("sign tx error: ", slog.String("error", err.Error()))
		return
	}

	_, err = rpcClient.SendTransaction(ctx, tx)
	if err != nil {
		slog.Error("send tx error: ", slog.String("error", err.Error()))
	}
}
