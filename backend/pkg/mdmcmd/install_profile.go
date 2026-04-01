package mdmcmd

// InstallProfile generates an InstallProfile MDM command.
// The profileData should be the raw mobileconfig XML/plist data.
// The plist library encodes []byte as a <data> element (Apple MDM requirement).
func (b *CommandBuilder) InstallProfile(profileData []byte) ([]byte, string, error) {
	payload := map[string]any{
		"Payload": profileData, // []byte → plist <data> type, not <string>
	}

	return buildCommand("InstallProfile", payload)
}
