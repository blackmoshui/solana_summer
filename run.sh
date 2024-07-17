#!/bin/bash

wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

go mod init solana_summer

go get github.com/gagliardetto/binary
go get github.com/gagliardetto/solana-go
go get github.com/gagliardetto/solana-go/rpc
go get github.com/go-resty/resty/v2
go get github.com/tidwall/gjson
go get golang.org/x/exp/slog

go run main.go