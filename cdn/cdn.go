package cdn

import (
	"fmt"
)

// GetHippoURL return the URL of the hippo on our CDN
// TODO: this is a fake URL for now
func GetHippoURL(id string) string {
	// TODO: This should go to a setting
	return fmt.Sprintf("http://cdn.apihippo.com:8000/%s", id)
}
