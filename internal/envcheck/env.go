package envcheck

import (
	"os"

	"speedy-cli/internal/common"

	"github.com/joho/godotenv"
)

type CheckItem struct {
	Key     string `json:"key"`
	Present bool   `json:"present"`
}

func Check(file string, required []string) (common.Result, []CheckItem) {
	_ = godotenv.Overload(file)
	items := make([]CheckItem, 0, len(required))
	missing := 0
	for _, key := range required {
		val, ok := os.LookupEnv(key)
		present := ok && val != ""
		if !present {
			missing++
		}
		items = append(items, CheckItem{Key: key, Present: present})
	}

	res := common.Result{Status: common.StatusSuccess, Message: "env check complete"}
	if missing > 0 {
		res.Status = common.StatusWarning
		res.Suggestion = "Add missing vars to your .env and restart the app"
	}
	return res, items
}
