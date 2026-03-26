package serviceimpl

import (
	"bytes"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	"howett.net/plist"
)

type profileGeneratorImpl struct {
	orgName string
	prefix  string
}

func NewProfileGenerator(orgName, prefix string) service.ProfileGenerator {
	if prefix == "" {
		prefix = "com.thd.mdm"
	}
	return &profileGeneratorImpl{
		orgName: orgName,
		prefix:  prefix,
	}
}

func (g *profileGeneratorImpl) GenerateXML(ctx context.Context, p *ent.Profile) ([]byte, error) {
	// Root Profile Dictionary
	profileMap := map[string]any{
		"PayloadDisplayName":       p.Name,
		"PayloadIdentifier":        fmt.Sprintf("%s.profile.%d", g.prefix, p.ID),
		"PayloadOrganization":      g.orgName,
		"PayloadRemovalDisallowed": false,
		"PayloadType":              "Configuration",
		"PayloadUUID":              uuid.New().String(),
		"PayloadVersion":           1,
		"PayloadContent":           []any{},
	}

	payloadContent := []any{}

	// Handle Passcode/Security Settings
	if len(p.SecuritySettings) > 0 {
		payloadContent = append(payloadContent, g.createPayload("com.apple.mobiledevice.passwordpolicy", p.SecuritySettings))
	}

	// Handle Restrictions
	if len(p.Restrictions) > 0 {
		payloadContent = append(payloadContent, g.createPayload("com.apple.applicationaccess", p.Restrictions))
	}

	// Handle Network/Wifi
	if len(p.NetworkConfig) > 0 {
		// Basic mapping, assuming network_config might contain multiple configs or one
		payloadContent = append(payloadContent, g.createPayload("com.apple.wifi.managed", p.NetworkConfig))
	}

	profileMap["PayloadContent"] = payloadContent

	// Encode to PList XML
	var buf bytes.Buffer
	encoder := plist.NewEncoder(&buf)
	encoder.Indent("\t")
	if err := encoder.Encode(profileMap); err != nil {
		return nil, fmt.Errorf("failed to encode profile to plist: %w", err)
	}

	return buf.Bytes(), nil
}

func (g *profileGeneratorImpl) createPayload(payloadType string, settings map[string]any) map[string]any {
	payload := map[string]any{
		"PayloadDisplayName":  payloadType,
		"PayloadIdentifier":   fmt.Sprintf("%s.payload.%s.%s", g.prefix, payloadType, uuid.New().String()),
		"PayloadOrganization": g.orgName,
		"PayloadType":         payloadType,
		"PayloadUUID":         uuid.New().String(),
		"PayloadVersion":      1,
	}

	// Merge settings into payload
	for k, v := range settings {
		payload[k] = v
	}

	return payload
}
