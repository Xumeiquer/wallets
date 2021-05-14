package main

import (
	"embed"

	"github.com/Xumeiquer/wallets/server/cmd"
)

//go:embed assets/*
var StaticContent embed.FS

func main() {
	cmd.Execute(StaticContent)
}
