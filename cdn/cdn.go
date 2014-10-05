package cdn

import (
	"fmt"

	"github.com/127biscuits/apihippo.com/settings"
)

// GetHippoURL return the URL of the hippo on our CDN
func GetHippoURL(id string) string {
	return fmt.Sprintf("%s/%s", settings.Config.CDN.Address, id)
}
