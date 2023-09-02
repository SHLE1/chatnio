package auth

import (
	"chat/utils"
	"encoding/json"
	"github.com/spf13/viper"
)

type BalanceResponse struct {
	Status  bool    `json:"status" required:"true"`
	Balance float32 `json:"balance"`
}

type PaymentResponse struct {
	Status bool `json:"status" required:"true"`
	Type   bool `json:"type"`
}

func GenerateOrder() string {
	return utils.Sha2Encrypt(utils.GenerateChar(32))
}

func GetBalance(username string) float32 {
	order := GenerateOrder()
	res, err := utils.Post("https://api.deeptrain.net/app/balance", map[string]string{
		"Content-Type": "application/json",
	}, map[string]interface{}{
		"password": viper.GetString("auth.access"),
		"user":     username,
		"hash":     utils.Sha2Encrypt(username + viper.GetString("auth.salt")),
		"order":    order,
		"sign":     utils.Sha2Encrypt(username + order + viper.GetString("auth.sign")),
	})

	if err != nil || res == nil || res.(map[string]interface{})["status"] == false {
		return 0.
	}

	converter, _ := json.Marshal(res)
	resp, _ := utils.Unmarshal[BalanceResponse](converter)
	return resp.Balance
}

func Pay(username string, amount float32) bool {
	order := GenerateOrder()
	res, err := utils.Post("https://api.deeptrain.net/app/payment", map[string]string{
		"Content-Type": "application/json",
	}, map[string]interface{}{
		"password": viper.GetString("auth.access"),
		"user":     username,
		"hash":     utils.Sha2Encrypt(username + viper.GetString("auth.salt")),
		"order":    order,
		"amount":   amount,
		"sign":     utils.Sha2Encrypt(username + order + viper.GetString("auth.sign")),
	})

	if err != nil || res == nil || res.(map[string]interface{})["status"] == false {
		return false
	}

	converter, _ := json.Marshal(res)
	resp, _ := utils.Unmarshal[PaymentResponse](converter)
	return resp.Type
}