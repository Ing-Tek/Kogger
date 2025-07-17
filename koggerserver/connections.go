package koggerserver

import (
	"encoding/json"
	"net/http"

	models "github.com/Ing-Tek/Kogger/models"
	"github.com/mediocregopher/radix/v3"
)

var (
	PkiVaultServiceHost string
	PkiVaultServicePort string
	JwtSecretKey        string
	RedisHost           string
	RedisPort           string
	RedisPassword       string
	RedisPool           *radix.Pool
)

func writeStandardResponse(r *http.Request, w http.ResponseWriter, grpcToken string, message string) {
	responseObj := &models.Response{
		Token:   grpcToken,
		Message: message,
	}

	response, _ := json.Marshal(responseObj)
	w.Write(response)
}
