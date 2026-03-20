package serviceimpl

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
)

type mobileConfigRepoStub struct {
	capturedEntity   *ent.MobileConfig
	capturedPayloads []*ent.Payload
	conflict         *repository.UniqueFieldConflict
	updateConflict   *repository.UniqueFieldConflict
}

func (s *mobileConfigRepoStub) GetFullForExport(ctx context.Context, id uint) (*ent.MobileConfig, error) {
	return nil, nil
}

func (s *mobileConfigRepoStub) Create(ctx context.Context, entity *ent.MobileConfig, payload []*ent.Payload) (*ent.MobileConfig, error) {
	s.capturedEntity = entity
	s.capturedPayloads = payload
	entity.Edges = ent.MobileConfigEdges{Payloads: payload}
	return entity, nil
}

func (s *mobileConfigRepoStub) FindCreateUniqueFieldConflict(ctx context.Context, name string, payloadIdentifier string, payloadIdentifiers []string) (*repository.UniqueFieldConflict, error) {
	return s.conflict, nil
}

func (s *mobileConfigRepoStub) GetByID(ctx context.Context, id uint) (*ent.MobileConfig, error) {
	return &ent.MobileConfig{ID: id}, nil
}

func (s *mobileConfigRepoStub) Update(ctx context.Context, id uint, entity *ent.MobileConfig, payload []*ent.Payload) (*ent.MobileConfig, error) {
	entity.ID = id
	entity.Edges = ent.MobileConfigEdges{Payloads: payload}
	return entity, nil
}

func (s *mobileConfigRepoStub) Delete(ctx context.Context, id uint) error {
	return nil
}

func (s *mobileConfigRepoStub) FindUpdateUniqueFieldConflict(ctx context.Context, id uint, name string, payloadIdentifier string, payloadIdentifiers []string) (*repository.UniqueFieldConflict, error) {
	return s.updateConflict, nil
}

func TestMobileConfigService_Create_AutoGeneratesUUIDAndDefaultVersions(t *testing.T) {
	repo := &mobileConfigRepoStub{}
	svc := NewMobileConfigService(repo)

	created, err := svc.Create(context.Background(), service.CreateMobileConfigCommand{
		Name:               "Corp Profile",
		PayloadIdentifier:  "com.thd.mobileconfig",
		PayloadType:        "Configuration",
		PayloadDisplayName: "THD Profile",
		Payloads: []service.CreateMobileConfigPayloadCommand{
			{
				PayloadDisplayName: "WiFi",
				PayloadIdentifier:  "com.thd.wifi",
				PayloadType:        "com.apple.wifi.managed",
				Properties: []service.CreateMobileConfigPropertyCommand{
					{Key: "wifiSSID_STR", ValueJSON: map[string]interface{}{"value": "THD"}},
				},
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, repo.capturedEntity)
	require.Len(t, repo.capturedPayloads, 1)

	_, err = uuid.Parse(repo.capturedEntity.PayloadUUID)
	require.NoError(t, err)
	require.Equal(t, 1, repo.capturedEntity.PayloadVersion)

	_, err = uuid.Parse(repo.capturedPayloads[0].PayloadUUID)
	require.NoError(t, err)
	require.Equal(t, 1, repo.capturedPayloads[0].PayloadVersion)
}

func TestMobileConfigService_Create_RequiresAtLeastOnePayload(t *testing.T) {
	repo := &mobileConfigRepoStub{}
	svc := NewMobileConfigService(repo)

	created, err := svc.Create(context.Background(), service.CreateMobileConfigCommand{
		Name:               "Corp Profile",
		PayloadIdentifier:  "com.thd.mobileconfig",
		PayloadType:        "Configuration",
		PayloadDisplayName: "THD Profile",
		Payloads:           []service.CreateMobileConfigPayloadCommand{},
	})

	require.Nil(t, created)
	require.Error(t, err)
}

func TestMobileConfigService_Create_RejectsWhitespaceRequiredFields(t *testing.T) {
	repo := &mobileConfigRepoStub{}
	svc := NewMobileConfigService(repo)

	created, err := svc.Create(context.Background(), service.CreateMobileConfigCommand{
		Name:               "   ",
		PayloadIdentifier:  "com.thd.mobileconfig",
		PayloadType:        "Configuration",
		PayloadDisplayName: "THD Profile",
		Payloads: []service.CreateMobileConfigPayloadCommand{
			{
				PayloadDisplayName: "WiFi",
				PayloadIdentifier:  "com.thd.wifi",
				PayloadType:        "com.apple.wifi.managed",
			},
		},
	})

	require.Nil(t, created)
	require.Error(t, err)
	appErr, ok := err.(*apperror.AppError)
	require.True(t, ok)
	require.Equal(t, apperror.ErrValidation.Code, appErr.Code)
	require.Equal(t, "name là bắt buộc", appErr.Message)
}

func TestMobileConfigService_Create_RejectsEmptyPropertyKeyInNestedPayload(t *testing.T) {
	repo := &mobileConfigRepoStub{}
	svc := NewMobileConfigService(repo)

	created, err := svc.Create(context.Background(), service.CreateMobileConfigCommand{
		Name:               "Corp Profile",
		PayloadIdentifier:  "com.thd.mobileconfig",
		PayloadType:        "Configuration",
		PayloadDisplayName: "THD Profile",
		Payloads: []service.CreateMobileConfigPayloadCommand{
			{
				PayloadDisplayName: "WiFi",
				PayloadIdentifier:  "com.thd.wifi",
				PayloadType:        "com.apple.wifi.managed",
				Properties: []service.CreateMobileConfigPropertyCommand{
					{Key: "   "},
				},
			},
		},
	})

	require.Nil(t, created)
	require.Error(t, err)
	appErr, ok := err.(*apperror.AppError)
	require.True(t, ok)
	require.Equal(t, apperror.ErrValidation.Code, appErr.Code)
	require.Equal(t, "payloads[0].properties[0].key là bắt buộc", appErr.Message)
}

func TestMobileConfigService_Create_RejectsDuplicatePayloadIdentifierInRequest(t *testing.T) {
	repo := &mobileConfigRepoStub{}
	svc := NewMobileConfigService(repo)

	created, err := svc.Create(context.Background(), service.CreateMobileConfigCommand{
		Name:               "Corp Profile",
		PayloadIdentifier:  "com.thd.mobileconfig",
		PayloadType:        "Configuration",
		PayloadDisplayName: "THD Profile",
		Payloads: []service.CreateMobileConfigPayloadCommand{
			{PayloadDisplayName: "WiFi 1", PayloadIdentifier: "com.thd.wifi", PayloadType: "com.apple.wifi.managed"},
			{PayloadDisplayName: "WiFi 2", PayloadIdentifier: "com.thd.wifi", PayloadType: "com.apple.wifi.managed"},
		},
	})

	require.Nil(t, created)
	require.Error(t, err)
	appErr, ok := err.(*apperror.AppError)
	require.True(t, ok)
	require.Equal(t, apperror.ErrConflict.Code, appErr.Code)
	require.Equal(t, "payloads[1].payload_identifier bị trùng trong request: com.thd.wifi", appErr.Message)
}

func TestMobileConfigService_Create_RejectsExistingUniqueFieldConflict(t *testing.T) {
	repo := &mobileConfigRepoStub{conflict: &repository.UniqueFieldConflict{Field: "name", Value: "Corp Profile"}}
	svc := NewMobileConfigService(repo)

	created, err := svc.Create(context.Background(), service.CreateMobileConfigCommand{
		Name:               "Corp Profile",
		PayloadIdentifier:  "com.thd.mobileconfig",
		PayloadType:        "Configuration",
		PayloadDisplayName: "THD Profile",
		Payloads: []service.CreateMobileConfigPayloadCommand{
			{PayloadDisplayName: "WiFi", PayloadIdentifier: "com.thd.wifi", PayloadType: "com.apple.wifi.managed"},
		},
	})

	require.Nil(t, created)
	require.Error(t, err)
	appErr, ok := err.(*apperror.AppError)
	require.True(t, ok)
	require.Equal(t, apperror.ErrConflict.Code, appErr.Code)
	require.Equal(t, "name đã tồn tại: Corp Profile", appErr.Message)
}
