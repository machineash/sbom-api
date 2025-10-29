package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

func GetSecret() (map[string]interface{}, error) {
	vaultAddr := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")

	if vaultAddr == "" || vaultToken == "" {
		return nil, fmt.Errorf("VAULT_ADDR or VAULT_TOKEN not set")
	}

	client, err := api.NewClient(&api.Config{Address: vaultAddr})
	if err != nil {
		return nil, fmt.Errorf("vault client init error: %w", err)
	}

	client.SetToken(vaultToken)

	secret, err := client.Logical().Read("secret/data/app/config")
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}
	if secret == nil || secret.Data["data"] == nil {
		return nil, fmt.Errorf("no data found at path")
	}

	data := secret.Data["data"].(map[string]interface{})
	return data, nil
}
