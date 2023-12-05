package main

import (
	"gitag.ir/cookthepot/services/vault/cmd/chef/cmd"
	"gitag.ir/cookthepot/services/vault/config"
)

func main() {
	config.Load()

	cmd.Execute()
}
