package validateiap

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/imkira/gcp-iap-auth/jwt"
)

var cfg *jwt.Config

func init() {
	cfg = &jwt.Config{}

	aud, err := getAudience()
	if err != nil {
		log.Fatal(err)
	}

	err = initAudiences(aud)
	if err != nil {
		log.Fatal(err)
	}

	cfg.PublicKeys, err = jwt.FetchPublicKeys()
	if err != nil {
		log.Fatal(err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
}

func initAudiences(audiences string) error {
	str, err := extractAudiencesRegexp(audiences)
	if err != nil {
		return err
	}
	re, err := regexp.Compile(str)
	if err != nil {
		return fmt.Errorf("Invalid audiences regular expression %q (%v)", str, err)
	}
	cfg.MatchAudiences = re
	return nil
}

func extractAudiencesRegexp(audiences string) (string, error) {
	var strs []string
	for _, audience := range strings.Split(audiences, ",") {
		str, err := extractAudienceRegexp(audience)
		if err != nil {
			return "", err
		}
		strs = append(strs, str)
	}
	return strings.Join(strs, "|"), nil
}

func extractAudienceRegexp(audience string) (string, error) {
	if strings.HasPrefix(audience, "/") && strings.HasSuffix(audience, "/") {
		if len(audience) < 3 {
			return "", fmt.Errorf("Invalid audiences regular expression %q", audience)
		}
		return audience[1 : len(audience)-1], nil
	}
	return parseRawAudience(audience)
}

func parseRawAudience(audience string) (string, error) {
	_, err := jwt.ParseAudience(audience)
	if err != nil {
		return "", fmt.Errorf("Invalid audience %q (%v)", audience, err)
	}
	return fmt.Sprintf("^%s$", regexp.QuoteMeta(audience)), nil
}
