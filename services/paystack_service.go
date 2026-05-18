package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ResolveAccountResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
		BankID        int    `json:"bank_id"`
	} `json:"data"`
}

func ResolveAccountNumber(accountNumber, bankCode string) (string, error) {
	secretKey := os.Getenv("PAYSTACK_SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("Paystack secret key not set")
	}

	url := fmt.Sprintf("https://api.paystack.co/bank/resolve?account_number=%s&bank_code=%s", accountNumber, bankCode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to resolve account: status %d", resp.StatusCode)
	}

	var result ResolveAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if !result.Status {
		return "", fmt.Errorf("resolution failed: %s", result.Message)
	}

	return result.Data.AccountName, nil
}
