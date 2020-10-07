package micro

import (
	"os"

	"github.com/avayayu/micro/net/constants"
)

func JudgeEnv() constants.ServerMode {
	mode := os.Getenv("cloudbrainMode")

	if mode == constants.Prod.String() {
		return constants.Prod
	} else {
		return constants.Dev
	}
}
