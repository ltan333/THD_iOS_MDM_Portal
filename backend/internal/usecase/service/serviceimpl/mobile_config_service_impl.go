package serviceimpl

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
)

type mobileConfigServiceImpl struct {
	mobileConfigRepo repository.MobileConfigRepository
}

func NewMobileConfigService(mobileConfigRepo repository.MobileConfigRepository) service.MobileConfigService {
	return &mobileConfigServiceImpl{mobileConfigRepo: mobileConfigRepo}
}

func (m *mobileConfigServiceImpl) GenerateXML(ctx context.Context, cmd service.GenerateMobileConfigXMLCommand) ([]byte, error) {
	mc, err := m.mobileConfigRepo.GetFullForExport(ctx, cmd.ID)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("MobileConfig không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất MobileConfig").WithError(err)
	}

	xmlBytes, err := buildMobileConfigXML(mc)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo XML").WithError(err)
	}

	return xmlBytes, nil
}

func buildMobileConfigXML(mc *ent.MobileConfig) ([]byte, error) {
	root := map[string]interface{}{
		"PayloadDescription":       mc.PayloadDescription,
		"PayloadDisplayName":       mc.PayloadDisplayName,
		"PayloadIdentifier":        mc.PayloadIdentifier,
		"PayloadOrganization":      mc.PayloadOrganization,
		"PayloadRemovalDisallowed": mc.PayloadRemovalDisallowed,
		"PayloadType":              mc.PayloadType,
		"PayloadUUID":              mc.PayloadUUID,
		"PayloadVersion":           mc.PayloadVersion,
		"PayloadContent":           buildPayloadContent(mc.Edges.Payloads),
	}

	// Create XML with indent
	buf := &bytes.Buffer{}
	enc := xml.NewEncoder(buf)
	enc.Indent("", "  ")

	// Header + DOCTYPE
	fmt.Fprint(buf, xml.Header)
	fmt.Fprintln(buf, `<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">`)

	plistStart := xml.StartElement{Name: xml.Name{Local: "plist"}, Attr: []xml.Attr{{Name: xml.Name{Local: "version"}, Value: "1.0"}}}

	if err := enc.EncodeElement(root, plistStart); err != nil {
		return nil, err
	}

	if err := enc.Flush(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buildPayloadContent(payloads []*ent.Payload) []map[string]interface{} {
	var content []map[string]interface{}

	for _, p := range payloads {
		payloadDict := map[string]interface{}{
			"PayloadType":         p.PayloadType,
			"PayloadVersion":      p.PayloadVersion,
			"PayloadIdentifier":   p.PayloadIdentifier,
			"PayloadUUID":         p.PayloadUUID,
			"PayloadDisplayName":  p.PayloadDisplayName,
			"PayloadDescription":  p.PayloadDescription,
			"PayloadOrganization": p.PayloadOrganization,
		}

		// Convert properties to key-value pairs in the payload content
		for _, prop := range p.Edges.Properties {
			if prop.Edges.Definition == nil {
				continue
			}
			key := prop.Edges.Definition.Key
			val := prop.ValueJSON

			// Normalize value based on the defined value type
			payloadDict[key] = normalizeValue(val, prop.Edges.Definition.ValueType)
		}

		content = append(content, payloadDict)
	}

	return content
}

func normalizeValue(val interface{}, valueType string) interface{} {
	switch strings.ToLower(valueType) {
	case "bool":
		if b, ok := val.(bool); ok {
			return b
		}
		return false
	case "integer", "number":
		if f, ok := val.(float64); ok {
			return int(f)
		}
		return 0
	case "string":
		if s, ok := val.(string); ok {
			return s
		}
		return ""
	case "array":
		return val // []interface{}
	case "dict":
		return val // map[string]interface{}
	default:
		return val
	}
}
