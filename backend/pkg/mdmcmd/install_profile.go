package mdmcmd

import "encoding/base64"

// InstallProfile generates an InstallProfile MDM command.
// The profileData should be the raw mobileconfig XML/plist data.
func (b *CommandBuilder) InstallProfile(profileData []byte) ([]byte, string, error) {
	// Profile data must be base64 encoded in the Payload field
	encodedProfile := base64.StdEncoding.EncodeToString(profileData)

	payload := map[string]any{
		"Payload": encodedProfile,
	}

	return buildCommand("InstallProfile", payload)
}
