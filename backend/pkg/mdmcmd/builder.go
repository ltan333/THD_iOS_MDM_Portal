// Package mdmcmd provides utilities for building Apple MDM command plists.
package mdmcmd

import (
	"bytes"

	"github.com/google/uuid"
	"howett.net/plist"
)

// CommandBuilder generates MDM command plist data.
type CommandBuilder struct {
	orgPrefix string
}

// NewBuilder creates a new MDM command builder.
func NewBuilder(orgPrefix string) *CommandBuilder {
	if orgPrefix == "" {
		orgPrefix = "com.thd.mdm"
	}
	return &CommandBuilder{orgPrefix: orgPrefix}
}

// MDMCommand represents the structure of an Apple MDM command.
type MDMCommand struct {
	Command     map[string]any `plist:"Command"`
	CommandUUID string         `plist:"CommandUUID"`
}

// generateUUID creates a new UUID for command identification.
func generateUUID() string {
	return uuid.New().String()
}

// encodePlist encodes the given data structure to Apple plist XML format.
func encodePlist(data any) ([]byte, error) {
	var buf bytes.Buffer
	encoder := plist.NewEncoder(&buf)
	encoder.Indent("\t")
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// buildCommand creates an MDM command with the given request type and optional payload.
func buildCommand(requestType string, payload map[string]any) ([]byte, string, error) {
	cmdUUID := generateUUID()

	cmdDict := map[string]any{
		"RequestType": requestType,
	}

	// Merge payload into command dict
	for k, v := range payload {
		cmdDict[k] = v
	}

	mdmCmd := MDMCommand{
		Command:     cmdDict,
		CommandUUID: cmdUUID,
	}

	data, err := encodePlist(mdmCmd)
	if err != nil {
		return nil, "", err
	}

	return data, cmdUUID, nil
}
