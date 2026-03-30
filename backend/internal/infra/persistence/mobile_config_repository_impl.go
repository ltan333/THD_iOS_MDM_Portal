package persistence

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/mobileconfig"
	"github.com/thienel/go-backend-template/internal/ent/payload"
	"github.com/thienel/go-backend-template/internal/ent/payloadproperty"
	"github.com/thienel/go-backend-template/internal/ent/payloadpropertydefinition"
	"github.com/thienel/go-backend-template/internal/infra/database"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type mobileConfigRepositoryImpl struct {
	client *ent.Client
}

func NewMobileConfigRepository(client *ent.Client) repository.MobileConfigRepository {
	return &mobileConfigRepositoryImpl{client: client}
}

func (m *mobileConfigRepositoryImpl) GetByID(ctx context.Context, id uint) (*ent.MobileConfig, error) {
	mc, err := m.client.MobileConfig.Query().
		Where(mobileconfig.IDEQ(id), mobileconfig.DeletedAtIsNil()).
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("MobileConfig không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất MobileConfig").WithError(err)
	}

	return mc, nil
}

func (m *mobileConfigRepositoryImpl) GetByIDWithPayloads(ctx context.Context, id uint) (*ent.MobileConfig, error) {
	mc, err := m.client.MobileConfig.Query().
		Where(mobileconfig.IDEQ(id), mobileconfig.DeletedAtIsNil()).
		WithPayloads(func(q *ent.PayloadQuery) {
			q.Where(payload.DeletedAtIsNil()).
				WithProperties(func(pq *ent.PayloadPropertyQuery) {
					pq.Where(payloadproperty.DeletedAtIsNil()).WithDefinition()
				})
		}).
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("MobileConfig không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất MobileConfig").WithError(err)
	}

	return mc, nil
}

func (m *mobileConfigRepositoryImpl) FindCreateUniqueFieldConflict(ctx context.Context, name string, payloadIdentifier string, payloadIdentifiers []string) (*repository.UniqueFieldConflict, error) {
	exists, err := m.client.MobileConfig.Query().Where(mobileconfig.NameEQ(name)).Exist(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra unique field name").WithError(err)
	}
	if exists {
		return &repository.UniqueFieldConflict{Field: "name", Value: name}, nil
	}

	exists, err = m.client.MobileConfig.Query().Where(mobileconfig.PayloadIdentifierEQ(payloadIdentifier)).Exist(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra unique field payload_identifier").WithError(err)
	}
	if exists {
		return &repository.UniqueFieldConflict{Field: "payload_identifier", Value: payloadIdentifier}, nil
	}

	if len(payloadIdentifiers) == 0 {
		return nil, nil
	}

	dupPayload, err := m.client.Payload.Query().Where(payload.PayloadIdentifierIn(payloadIdentifiers...)).First(ctx)
	if ent.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra unique field payloads[].payload_identifier").WithError(err)
	}

	return &repository.UniqueFieldConflict{Field: "payloads[].payload_identifier", Value: dupPayload.PayloadIdentifier}, nil
}

func (m *mobileConfigRepositoryImpl) FindUpdateUniqueFieldConflict(ctx context.Context, id uint, name string, payloadIdentifier string, payloadIdentifiers []string) (*repository.UniqueFieldConflict, error) {
	exists, err := m.client.MobileConfig.Query().
		Where(
			mobileconfig.NameEQ(name),
			mobileconfig.IDNEQ(id),
		).
		Exist(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra unique field name").WithError(err)
	}
	if exists {
		return &repository.UniqueFieldConflict{Field: "name", Value: name}, nil
	}

	exists, err = m.client.MobileConfig.Query().
		Where(
			mobileconfig.PayloadIdentifierEQ(payloadIdentifier),
			mobileconfig.IDNEQ(id),
		).
		Exist(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra unique field payload_identifier").WithError(err)
	}
	if exists {
		return &repository.UniqueFieldConflict{Field: "payload_identifier", Value: payloadIdentifier}, nil
	}

	if len(payloadIdentifiers) == 0 {
		return nil, nil
	}

	dupPayload, err := m.client.Payload.Query().
		Where(
			payload.PayloadIdentifierIn(payloadIdentifiers...),
			payload.HasMobileConfigWith(mobileconfig.IDNEQ(id)),
		).
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra unique field payloads[].payload_identifier").WithError(err)
	}

	return &repository.UniqueFieldConflict{Field: "payloads[].payload_identifier", Value: dupPayload.PayloadIdentifier}, nil
}

func (m *mobileConfigRepositoryImpl) Create(ctx context.Context, entity *ent.MobileConfig, payload []*ent.Payload) (*ent.MobileConfig, error) {
	if entity == nil {
		return nil, apperror.ErrBadRequest.WithMessage("MobileConfig is required")
	}

	var mc *ent.MobileConfig
	err := database.WithTx(ctx, func(tx *ent.Tx) error {
		var txErr error
		// Create MobileConfig.
		mc, txErr = tx.MobileConfig.Create().
			SetName(entity.Name).
			SetPayloadIdentifier(entity.PayloadIdentifier).
			SetPayloadType(entity.PayloadType).
			SetPayloadDisplayName(entity.PayloadDisplayName).
			SetPayloadDescription(entity.PayloadDescription).
			SetPayloadOrganization(entity.PayloadOrganization).
			SetPayloadUUID(entity.PayloadUUID).
			SetPayloadVersion(entity.PayloadVersion).
			SetPayloadRemovalDisallowed(entity.PayloadRemovalDisallowed).
			Save(ctx)
		if txErr != nil {
			if ent.IsConstraintError(txErr) {
				return apperror.ErrConflict.WithMessage("Tên hoặc PayloadIdentifier đã tồn tại")
			}
			return apperror.ErrInternalServerError.WithMessage("Failed to create MobileConfig").WithError(txErr)
		}

		// Create Payloads and their Properties.
		for _, p := range payload {
			if p == nil {
				return apperror.ErrBadRequest.WithMessage("Payload item is required")
			}

			pl, txErr := tx.Payload.Create().
				SetMobileConfigID(mc.ID).
				SetPayloadDescription(p.PayloadDescription).
				SetPayloadDisplayName(p.PayloadDisplayName).
				SetPayloadIdentifier(p.PayloadIdentifier).
				SetPayloadOrganization(p.PayloadOrganization).
				SetPayloadType(p.PayloadType).
				SetPayloadUUID(p.PayloadUUID).
				SetPayloadVersion(p.PayloadVersion).
				Save(ctx)
			if txErr != nil {
				return apperror.ErrInternalServerError.WithMessage("Failed to create Payload").WithError(txErr)
			}

			for _, prop := range p.Edges.Properties {
				if prop == nil || prop.Edges.Definition == nil {
					return apperror.ErrBadRequest.WithMessage("Payload property definition is required")
				}

				// Validate definition by payload type and key.
				def, txErr := tx.PayloadPropertyDefinition.Query().
					Where(
						payloadpropertydefinition.PayloadTypeEQ(p.PayloadType),
						payloadpropertydefinition.KeyEQ(prop.Edges.Definition.Key),
					).
					First(ctx)
				if ent.IsNotFound(txErr) {
					return apperror.ErrBadRequest.WithMessage("Invalid property definition for payload type").WithError(txErr)
				}
				if txErr != nil {
					return apperror.ErrInternalServerError.WithMessage("Error validating property definition").WithError(txErr)
				}

				if _, txErr := tx.PayloadProperty.Create().
					SetPayloadID(pl.ID).
					SetDefinitionID(def.ID).
					SetValueJSON(prop.ValueJSON).
					Save(ctx); txErr != nil {
					return apperror.ErrInternalServerError.WithMessage("Failed to create payload property").WithError(txErr)
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return m.client.MobileConfig.Query().
		Where(mobileconfig.IDEQ(mc.ID)).
		WithPayloads(func(q *ent.PayloadQuery) {
			q.WithProperties(func(pq *ent.PayloadPropertyQuery) {
				pq.WithDefinition()
			})
		}).
		First(ctx)
}

func (m *mobileConfigRepositoryImpl) GetFullForExport(ctx context.Context, id uint) (*ent.MobileConfig, error) {
	return m.client.MobileConfig.Query().
		Where(mobileconfig.IDEQ(id), mobileconfig.DeletedAtIsNil()).
		WithPayloads(func(q *ent.PayloadQuery) {
			q.Where(payload.DeletedAtIsNil()).
				WithProperties(func(pq *ent.PayloadPropertyQuery) {
					pq.Where(payloadproperty.DeletedAtIsNil()).WithDefinition()
				})
		}).
		First(ctx)
}
func (m *mobileConfigRepositoryImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.MobileConfig, int64, error) {
	q := m.client.MobileConfig.Query()

	// Apply search filter
	if searchFilter, ok := opts.Filters["search"]; ok {
		searchValue := searchFilter.Value.(string)
		q = q.Where(mobileconfig.Or(
			mobileconfig.NameContainsFold(searchValue),
			mobileconfig.PayloadTypeContainsFold(searchValue),
			mobileconfig.PayloadDisplayNameContainsFold(searchValue),
			mobileconfig.PayloadIdentifierContainsFold(searchValue),
		))
	}

	// Apply other filters
	if nameFilter, ok := opts.Filters["name"]; ok {
		nameValue := nameFilter.Value.(string)
		q = q.Where(mobileconfig.NameEQ(nameValue))
	}
	if payloadTypeFilter, ok := opts.Filters["payload_type"]; ok {
		payloadTypeValue := payloadTypeFilter.Value.(string)
		q = q.Where(mobileconfig.PayloadTypeEQ(payloadTypeValue))
	}
	if payloadDisplayNameFilter, ok := opts.Filters["payload_display_name"]; ok {
		payloadDisplayNameValue := payloadDisplayNameFilter.Value.(string)
		q = q.Where(mobileconfig.PayloadDisplayNameEQ(payloadDisplayNameValue))
	}

	// Get total count before pagination
	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Failed to count mobile configs").WithError(err)
	}

	// Apply sorting with default
	if len(opts.Sort) > 0 {
		for _, sort := range opts.Sort {
			if sort.Desc {
				q = q.Order(ent.Desc(sort.Field))
			} else {
				q = q.Order(ent.Asc(sort.Field))
			}
		}
	} else {
		q = q.Order(ent.Desc(mobileconfig.FieldCreatedAt))
	}

	// Apply pagination
	items, err := q.Offset(offset).Limit(limit).All(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Failed to list mobile configs").WithError(err)
	}

	return items, int64(total), nil
}
func (m *mobileConfigRepositoryImpl) Update(ctx context.Context, id uint, entity *ent.MobileConfig, payloads []*ent.Payload) (*ent.MobileConfig, error) {
	if entity == nil {
		return nil, apperror.ErrBadRequest.WithMessage("MobileConfig is required")
	}

	var mc *ent.MobileConfig
	err := database.WithTx(ctx, func(tx *ent.Tx) error {
		var txErr error
		mc, txErr = tx.MobileConfig.UpdateOneID(id).
			SetName(entity.Name).
			SetPayloadIdentifier(entity.PayloadIdentifier).
			SetPayloadType(entity.PayloadType).
			SetPayloadDisplayName(entity.PayloadDisplayName).
			SetPayloadDescription(entity.PayloadDescription).
			SetPayloadOrganization(entity.PayloadOrganization).
			SetPayloadVersion(entity.PayloadVersion).
			SetPayloadRemovalDisallowed(entity.PayloadRemovalDisallowed).
			Save(ctx)
		if ent.IsNotFound(txErr) {
			return apperror.ErrNotFound.WithMessage("MobileConfig không tồn tại")
		}
		if ent.IsConstraintError(txErr) {
			return apperror.ErrConflict.WithMessage("Tên hoặc PayloadIdentifier đã tồn tại")
		}
		if txErr != nil {
			return apperror.ErrInternalServerError.WithMessage("Failed to update MobileConfig").WithError(txErr)
		}

		existingPayloads, txErr := tx.Payload.Query().
			Where(payload.HasMobileConfigWith(mobileconfig.IDEQ(id))).
			All(ctx)
		if txErr != nil {
			return apperror.ErrInternalServerError.WithMessage("Failed to query existing payloads").WithError(txErr)
		}

		payloadIDs := make([]uint, 0, len(existingPayloads))
		for _, p := range existingPayloads {
			payloadIDs = append(payloadIDs, p.ID)
		}

		if len(payloadIDs) > 0 {
			if _, txErr := tx.PayloadProperty.Delete().Where(payloadproperty.HasPayloadWith(payload.IDIn(payloadIDs...))).Exec(ctx); txErr != nil {
				return apperror.ErrInternalServerError.WithMessage("Failed to delete old payload properties").WithError(txErr)
			}
			if _, txErr := tx.Payload.Delete().Where(payload.IDIn(payloadIDs...)).Exec(ctx); txErr != nil {
				return apperror.ErrInternalServerError.WithMessage("Failed to delete old payloads").WithError(txErr)
			}
		}

		for _, p := range payloads {
			if p == nil {
				return apperror.ErrBadRequest.WithMessage("Payload item is required")
			}

			pl, txErr := tx.Payload.Create().
				SetMobileConfigID(mc.ID).
				SetPayloadDescription(p.PayloadDescription).
				SetPayloadDisplayName(p.PayloadDisplayName).
				SetPayloadIdentifier(p.PayloadIdentifier).
				SetPayloadOrganization(p.PayloadOrganization).
				SetPayloadType(p.PayloadType).
				SetPayloadUUID(p.PayloadUUID).
				SetPayloadVersion(p.PayloadVersion).
				Save(ctx)
			if txErr != nil {
				return apperror.ErrInternalServerError.WithMessage("Failed to create Payload").WithError(txErr)
			}

			for _, prop := range p.Edges.Properties {
				if prop == nil || prop.Edges.Definition == nil {
					return apperror.ErrBadRequest.WithMessage("Payload property definition is required")
				}

				def, txErr := tx.PayloadPropertyDefinition.Query().
					Where(
						payloadpropertydefinition.PayloadTypeEQ(p.PayloadType),
						payloadpropertydefinition.KeyEQ(prop.Edges.Definition.Key),
					).
					First(ctx)
				if ent.IsNotFound(txErr) {
					return apperror.ErrBadRequest.WithMessage("Invalid property definition for payload type").WithError(txErr)
				}
				if txErr != nil {
					return apperror.ErrInternalServerError.WithMessage("Error validating property definition").WithError(txErr)
				}

				if _, txErr := tx.PayloadProperty.Create().
					SetPayloadID(pl.ID).
					SetDefinitionID(def.ID).
					SetValueJSON(prop.ValueJSON).
					Save(ctx); txErr != nil {
					return apperror.ErrInternalServerError.WithMessage("Failed to create payload property").WithError(txErr)
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return m.client.MobileConfig.Query().
		Where(mobileconfig.IDEQ(mc.ID)).
		WithPayloads(func(q *ent.PayloadQuery) {
			q.WithProperties(func(pq *ent.PayloadPropertyQuery) {
				pq.WithDefinition()
			})
		}).
		First(ctx)
}

func (m *mobileConfigRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.WithTx(ctx, func(tx *ent.Tx) error {
		existingPayloads, txErr := tx.Payload.Query().
			Where(payload.HasMobileConfigWith(mobileconfig.IDEQ(id))).
			All(ctx)
		if txErr != nil {
			return apperror.ErrInternalServerError.WithMessage("Failed to query existing payloads").WithError(txErr)
		}

		payloadIDs := make([]uint, 0, len(existingPayloads))
		for _, p := range existingPayloads {
			payloadIDs = append(payloadIDs, p.ID)
		}

		if len(payloadIDs) > 0 {
			if _, txErr := tx.PayloadProperty.Delete().Where(payloadproperty.HasPayloadWith(payload.IDIn(payloadIDs...))).Exec(ctx); txErr != nil {
				return apperror.ErrInternalServerError.WithMessage("Failed to delete payload properties").WithError(txErr)
			}
			if _, txErr := tx.Payload.Delete().Where(payload.IDIn(payloadIDs...)).Exec(ctx); txErr != nil {
				return apperror.ErrInternalServerError.WithMessage("Failed to delete payloads").WithError(txErr)
			}
		}

		txErr = tx.MobileConfig.DeleteOneID(id).Exec(ctx)
		if ent.IsNotFound(txErr) {
			return apperror.ErrNotFound.WithMessage("MobileConfig không tồn tại")
		}
		if txErr != nil {
			return apperror.ErrInternalServerError.WithMessage("Failed to delete MobileConfig").WithError(txErr)
		}

		return nil
	})
}
