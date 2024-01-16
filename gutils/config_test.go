package gutils

import (
	"log"
	"testing"

	"github.com/Dcarbon/go-shared/libs/utils"
)

func TestGetDB(t *testing.T) {
	var config = Config{
		Port: utils.IntEnv("PORT", 9035),
		Name: "adfdf",
		DbUrl: utils.StringEnv(
			"DB",
			"postgres://admin:244466666@localhost:5433/custody",
		),
	}
	log.Println(config.GetDBUrl())
}
