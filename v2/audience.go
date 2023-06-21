package validateiap

import (
	"fmt"
	"os"

	"cloud.google.com/go/compute/metadata"
)

func getAudience() (string, error) {
	// Return the audience set in the environment first,
	// then fall back the auto-generating an App Engine compatible value.
	if aud := os.Getenv("IAP_AUDIENCE"); aud != "" {
		return aud, nil
	} else {
		return getAppEngineAudience()
	}
}

func getAppEngineAudience() (string, error) {
	projectNumber, err := metadata.NumericProjectID()
	if err != nil {
		return "", fmt.Errorf("metadata.NumericProjectID: %v", err)
	}

	projectID, err := metadata.ProjectID()
	if err != nil {
		return "", fmt.Errorf("metadata.ProjectID: %v", err)
	}

	return "/projects/" + projectNumber + "/apps/" + projectID, nil
}
