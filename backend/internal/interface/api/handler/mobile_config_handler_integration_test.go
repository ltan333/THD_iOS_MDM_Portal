package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/handler"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
)

type mobileConfigServiceStub struct {
	createFn func(ctx context.Context, cmd service.CreateMobileConfigCommand) (*ent.MobileConfig, error)
	updateFn func(ctx context.Context, cmd service.UpdateMobileConfigCommand) (*ent.MobileConfig, error)
	deleteFn func(ctx context.Context, id uint) error
}

func (s *mobileConfigServiceStub) Create(ctx context.Context, cmd service.CreateMobileConfigCommand) (*ent.MobileConfig, error) {
	if s.createFn != nil {
		return s.createFn(ctx, cmd)
	}
	return nil, nil
}

func (s *mobileConfigServiceStub) GenerateXML(ctx context.Context, cmd service.GenerateMobileConfigXMLCommand) ([]byte, error) {
	return nil, nil
}

func (s *mobileConfigServiceStub) Update(ctx context.Context, cmd service.UpdateMobileConfigCommand) (*ent.MobileConfig, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, cmd)
	}
	return nil, nil
}

func (s *mobileConfigServiceStub) Delete(ctx context.Context, id uint) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}

type createMobileConfigResponse struct {
	IsSuccess bool `json:"is_success"`
	Data      struct {
		ID                 uint   `json:"id"`
		Name               string `json:"name"`
		PayloadIdentifier  string `json:"payload_identifier"`
		PayloadDisplayName string `json:"payload_display_name"`
		PayloadType        string `json:"payload_type"`
	} `json:"data"`
	Message string `json:"message"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Fields  []struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		} `json:"fields"`
	} `json:"error"`
}

func TestCreateMobileConfig_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	called := false
	stub := &mobileConfigServiceStub{
		createFn: func(_ context.Context, cmd service.CreateMobileConfigCommand) (*ent.MobileConfig, error) {
			called = true
			require.Equal(t, "Corporate WiFi", cmd.Name)
			require.Len(t, cmd.Payloads, 1)
			require.Equal(t, "wifiSSID_STR", cmd.Payloads[0].Properties[0].Key)

			now := time.Now().UTC()
			return &ent.MobileConfig{
				ID:                 100,
				Name:               cmd.Name,
				PayloadIdentifier:  cmd.PayloadIdentifier,
				PayloadDisplayName: cmd.PayloadDisplayName,
				PayloadType:        cmd.PayloadType,
				PayloadUUID:        "generated-root-uuid",
				PayloadVersion:     1,
				CreatedAt:          now,
				UpdatedAt:          now,
				Edges: ent.MobileConfigEdges{
					Payloads: []*ent.Payload{
						{
							ID:                 200,
							PayloadType:        "com.apple.wifi.managed",
							PayloadIdentifier:  "com.thd.portal.wifi.payload",
							PayloadDisplayName: "WiFi Payload",
							PayloadUUID:        "payload-uuid-1",
							PayloadVersion:     1,
							Edges: ent.PayloadEdges{
								Properties: []*ent.PayloadProperty{
									{
										ID:        300,
										ValueJSON: map[string]interface{}{"value": "THD-Corp"},
										Edges: ent.PayloadPropertyEdges{
											Definition: &ent.PayloadPropertyDefinition{Key: "wifiSSID_STR"},
										},
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}

	r := gin.New()
	h := handler.NewMobileConfigHandler(stub)
	r.POST("/api/mobile-configs", h.Create)

	body := map[string]interface{}{
		"name":                 "Corporate WiFi",
		"payload_identifier":   "com.thd.portal.mobileconfig",
		"payload_type":         "Configuration",
		"payload_display_name": "THD Profile",
		"payload_version":      1,
		"payloads": []map[string]interface{}{
			{
				"payload_display_name": "WiFi Payload",
				"payload_identifier":   "com.thd.portal.wifi.payload",
				"payload_type":         "com.apple.wifi.managed",
				"payload_version":      1,
				"properties": []map[string]interface{}{
					{
						"key":        "wifiSSID_STR",
						"value_json": map[string]interface{}{"value": "THD-Corp"},
					},
				},
			},
		},
	}
	raw, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/mobile-configs", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	require.True(t, called)
	require.Equal(t, http.StatusCreated, res.Code)

	var got createMobileConfigResponse
	err = json.Unmarshal(res.Body.Bytes(), &got)
	require.NoError(t, err)
	require.True(t, got.IsSuccess)
	require.Equal(t, "Tạo mobile config thành công", got.Message)
	require.Equal(t, uint(100), got.Data.ID)
	require.Equal(t, "Corporate WiFi", got.Data.Name)
}

func TestCreateMobileConfig_Integration_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	called := false
	stub := &mobileConfigServiceStub{
		createFn: func(_ context.Context, _ service.CreateMobileConfigCommand) (*ent.MobileConfig, error) {
			called = true
			return nil, nil
		},
	}

	r := gin.New()
	h := handler.NewMobileConfigHandler(stub)
	r.POST("/api/mobile-configs", h.Create)

	invalidBody := map[string]interface{}{
		"payload_identifier":   "com.thd.bad",
		"payload_type":         "Configuration",
		"payload_display_name": "Bad",
		"payloads": []map[string]interface{}{
			{
				"payload_display_name": "WiFi Payload",
				"payload_identifier":   "com.thd.bad.payload",
				"payload_type":         "com.apple.wifi.managed",
				"properties": []map[string]interface{}{
					{
						"key":        "",
						"value_json": map[string]interface{}{"value": "X"},
					},
				},
			},
		},
	}
	raw, err := json.Marshal(invalidBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/mobile-configs", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	require.False(t, called)
	require.Equal(t, http.StatusBadRequest, res.Code)

	var got createMobileConfigResponse
	err = json.Unmarshal(res.Body.Bytes(), &got)
	require.NoError(t, err)
	require.False(t, got.IsSuccess)
	require.NotNil(t, got.Error)
	require.Equal(t, "VALIDATION_ERROR", got.Error.Code)
	require.NotEmpty(t, got.Error.Fields)
}

func TestCreateMobileConfig_Integration_WhitespaceValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	called := false
	stub := &mobileConfigServiceStub{
		createFn: func(_ context.Context, _ service.CreateMobileConfigCommand) (*ent.MobileConfig, error) {
			called = true
			return nil, nil
		},
	}

	r := gin.New()
	h := handler.NewMobileConfigHandler(stub)
	r.POST("/api/mobile-configs", h.Create)

	body := map[string]interface{}{
		"name":                 "   ",
		"payload_identifier":   "com.thd.bad",
		"payload_type":         "Configuration",
		"payload_display_name": "Bad",
		"payloads": []map[string]interface{}{
			{
				"payload_display_name": "WiFi Payload",
				"payload_identifier":   "com.thd.bad.payload",
				"payload_type":         "com.apple.wifi.managed",
			},
		},
	}
	raw, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/mobile-configs", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	require.False(t, called)
	require.Equal(t, http.StatusBadRequest, res.Code)

	var got createMobileConfigResponse
	err = json.Unmarshal(res.Body.Bytes(), &got)
	require.NoError(t, err)
	require.False(t, got.IsSuccess)
	require.NotNil(t, got.Error)
	require.Equal(t, "VALIDATION_ERROR", got.Error.Code)
	require.NotEmpty(t, got.Error.Fields)
	require.Equal(t, "name", got.Error.Fields[0].Field)
}

func TestCreateMobileConfig_Integration_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stub := &mobileConfigServiceStub{
		createFn: func(_ context.Context, _ service.CreateMobileConfigCommand) (*ent.MobileConfig, error) {
			return nil, apperror.ErrConflict.WithMessage("Payload identifier đã tồn tại")
		},
	}

	r := gin.New()
	h := handler.NewMobileConfigHandler(stub)
	r.POST("/api/mobile-configs", h.Create)

	body := map[string]interface{}{
		"name":                 "Conflict Config",
		"payload_identifier":   "com.thd.conflict",
		"payload_type":         "Configuration",
		"payload_display_name": "Conflict",
		"payloads": []map[string]interface{}{
			{
				"payload_display_name": "WiFi Payload",
				"payload_identifier":   "com.thd.conflict.payload",
				"payload_type":         "com.apple.wifi.managed",
			},
		},
	}
	raw, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/mobile-configs", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	require.Equal(t, http.StatusConflict, res.Code)

	var got createMobileConfigResponse
	err = json.Unmarshal(res.Body.Bytes(), &got)
	require.NoError(t, err)
	require.False(t, got.IsSuccess)
	require.NotNil(t, got.Error)
	require.Equal(t, "CONFLICT", got.Error.Code)
	require.Equal(t, "Payload identifier đã tồn tại", got.Error.Message)
}

func TestUpdateMobileConfig_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	called := false
	stub := &mobileConfigServiceStub{
		updateFn: func(_ context.Context, cmd service.UpdateMobileConfigCommand) (*ent.MobileConfig, error) {
			called = true
			require.Equal(t, uint(100), cmd.ID)
			require.Equal(t, "Updated Profile", cmd.Name)

			now := time.Now().UTC()
			return &ent.MobileConfig{
				ID:                 cmd.ID,
				Name:               cmd.Name,
				PayloadIdentifier:  cmd.PayloadIdentifier,
				PayloadDisplayName: cmd.PayloadDisplayName,
				PayloadType:        cmd.PayloadType,
				PayloadUUID:        "existing-root-uuid",
				PayloadVersion:     1,
				CreatedAt:          now,
				UpdatedAt:          now,
			}, nil
		},
	}

	r := gin.New()
	h := handler.NewMobileConfigHandler(stub)
	r.PUT("/api/mobile-configs/:id", h.Update)

	body := map[string]interface{}{
		"name":                 "Updated Profile",
		"payload_identifier":   "com.thd.portal.mobileconfig.updated",
		"payload_type":         "Configuration",
		"payload_display_name": "Updated THD Profile",
		"payloads": []map[string]interface{}{
			{
				"payload_display_name": "WiFi Payload",
				"payload_identifier":   "com.thd.portal.wifi.payload.updated",
				"payload_type":         "com.apple.wifi.managed",
			},
		},
	}
	raw, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/mobile-configs/100", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, res.Code)

	var got createMobileConfigResponse
	err = json.Unmarshal(res.Body.Bytes(), &got)
	require.NoError(t, err)
	require.True(t, got.IsSuccess)
	require.Equal(t, "Cập nhật mobile config thành công", got.Message)
	require.Equal(t, uint(100), got.Data.ID)
}

func TestDeleteMobileConfig_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	called := false
	stub := &mobileConfigServiceStub{
		deleteFn: func(_ context.Context, id uint) error {
			called = true
			require.Equal(t, uint(100), id)
			return nil
		},
	}

	r := gin.New()
	h := handler.NewMobileConfigHandler(stub)
	r.DELETE("/api/mobile-configs/:id", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/api/mobile-configs/100", nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, res.Code)

	var got createMobileConfigResponse
	err := json.Unmarshal(res.Body.Bytes(), &got)
	require.NoError(t, err)
	require.True(t, got.IsSuccess)
	require.Equal(t, "Xóa mobile config thành công", got.Message)
}
