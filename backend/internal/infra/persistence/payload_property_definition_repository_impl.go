package persistence

import (
	"context"
	"sort"
	"strings"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/payloadpropertydefinition"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

var payloadPropertyDefinitionAllowedFields = map[string]bool{
	"id":              true,
	"payload_type":    true,
	"payload_variant": true,
	"key":             true,
	"value_type":      true,
	"deprecated":      true,
	"description":     true,
	"order_index":     true,
	"created_at":      true,
	"updated_at":      true,
	"search":          true,
}

type payloadPropertyDefinitionRepositoryImpl struct {
	client *ent.Client
}

// NewPayloadPropertyDefinitionRepository creates a new payload property definition repository.
func NewPayloadPropertyDefinitionRepository(client *ent.Client) repository.PayloadPropertyDefinitionRepository {
	return &payloadPropertyDefinitionRepositoryImpl{client: client}
}

func (r *payloadPropertyDefinitionRepositoryImpl) Create(ctx context.Context, e *ent.PayloadPropertyDefinition) error {
	created, err := r.client.PayloadPropertyDefinition.Create().
		SetPayloadType(strings.TrimSpace(e.PayloadType)).
		SetPayloadVariant(strings.TrimSpace(e.PayloadVariant)).
		SetKey(strings.TrimSpace(e.Key)).
		SetValueType(strings.TrimSpace(e.ValueType)).
		SetNillableItemsType(e.ItemsType).
		SetDefaultValue(e.DefaultValue).
		SetEnumValues(e.EnumValues).
		SetNillableTitle(e.Title).
		SetDescription(strings.TrimSpace(e.Description)).
		SetPresence(e.Presence).
		SetDeprecated(e.Deprecated).
		SetIsNested(e.IsNested).
		SetNillableNestedReference(e.NestedReference).
		SetNillableItemsReference(e.ItemsReference).
		SetSupportedOs(e.SupportedOs).
		SetConditions(e.Conditions).
		SetOrderIndex(e.OrderIndex).
		SetNillableYamlSourceFile(e.YamlSourceFile).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return apperror.ErrConflict.WithMessage("Định nghĩa thuộc tính payload đã tồn tại").WithError(err)
		}
		return wrapCreateError(err, "định nghĩa thuộc tính payload")
	}

	e.ID = created.ID
	e.CreatedAt = created.CreatedAt
	e.UpdatedAt = created.UpdatedAt
	return nil
}

func (r *payloadPropertyDefinitionRepositoryImpl) FindByID(ctx context.Context, id uint) (*ent.PayloadPropertyDefinition, error) {
	entity, err := r.client.PayloadPropertyDefinition.Query().
		Where(payloadpropertydefinition.IDEQ(int(id))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy định nghĩa thuộc tính payload")
		}
		return nil, wrapFindError(err, "định nghĩa thuộc tính payload")
	}
	return entity, nil
}

func (r *payloadPropertyDefinitionRepositoryImpl) Update(ctx context.Context, e *ent.PayloadPropertyDefinition) error {
	updated, err := r.client.PayloadPropertyDefinition.UpdateOneID(e.ID).
		SetPayloadType(strings.TrimSpace(e.PayloadType)).
		SetPayloadVariant(strings.TrimSpace(e.PayloadVariant)).
		SetKey(strings.TrimSpace(e.Key)).
		SetValueType(strings.TrimSpace(e.ValueType)).
		SetNillableItemsType(e.ItemsType).
		SetDefaultValue(e.DefaultValue).
		SetEnumValues(e.EnumValues).
		SetNillableTitle(e.Title).
		SetDescription(strings.TrimSpace(e.Description)).
		SetPresence(e.Presence).
		SetDeprecated(e.Deprecated).
		SetIsNested(e.IsNested).
		SetNillableNestedReference(e.NestedReference).
		SetNillableItemsReference(e.ItemsReference).
		SetSupportedOs(e.SupportedOs).
		SetConditions(e.Conditions).
		SetOrderIndex(e.OrderIndex).
		SetNillableYamlSourceFile(e.YamlSourceFile).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return apperror.ErrConflict.WithMessage("Định nghĩa thuộc tính payload đã tồn tại").WithError(err)
		}
		return wrapUpdateError(err, "định nghĩa thuộc tính payload")
	}
	e.UpdatedAt = updated.UpdatedAt
	return nil
}

func (r *payloadPropertyDefinitionRepositoryImpl) Delete(ctx context.Context, id uint) error {
	if err := r.client.PayloadPropertyDefinition.DeleteOneID(int(id)).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return apperror.ErrNotFound.WithMessage("Không tìm thấy định nghĩa thuộc tính payload")
		}
		return wrapDeleteError(err, "định nghĩa thuộc tính payload")
	}
	return nil
}

func (r *payloadPropertyDefinitionRepositoryImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]ent.PayloadPropertyDefinition, int64, error) {
	entities, total, err := r.ListWithQuery(ctx, offset, limit, opts)
	if err != nil {
		return nil, 0, err
	}

	result := make([]ent.PayloadPropertyDefinition, len(entities))
	for i, item := range entities {
		result[i] = *item
	}
	return result, total, nil
}

func (r *payloadPropertyDefinitionRepositoryImpl) Exists(ctx context.Context, id uint) (bool, error) {
	count, err := r.client.PayloadPropertyDefinition.Query().Where(payloadpropertydefinition.IDEQ(int(id))).Count(ctx)
	if err != nil {
		return false, wrapFindError(err, "định nghĩa thuộc tính payload")
	}
	return count > 0, nil
}

func (r *payloadPropertyDefinitionRepositoryImpl) ListWithQuery(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.PayloadPropertyDefinition, int64, error) {
	q := r.client.PayloadPropertyDefinition.Query()

	if filter, ok := opts.Filters["payload_type"]; ok {
		if value, ok := filter.Value.(string); ok && strings.TrimSpace(value) != "" {
			q = q.Where(payloadpropertydefinition.PayloadTypeContainsFold(strings.TrimSpace(value)))
		}
	}

	if filter, ok := opts.Filters["payload_variant"]; ok {
		if value, ok := filter.Value.(string); ok {
			trimmed := strings.TrimSpace(value)
			if trimmed == "" {
				q = q.Where(payloadpropertydefinition.PayloadVariantEQ(""))
			} else {
				q = q.Where(payloadpropertydefinition.PayloadVariantContainsFold(trimmed))
			}
		}
	}

	if filter, ok := opts.Filters["key"]; ok {
		if value, ok := filter.Value.(string); ok && strings.TrimSpace(value) != "" {
			q = q.Where(payloadpropertydefinition.KeyContainsFold(strings.TrimSpace(value)))
		}
	}

	if filter, ok := opts.Filters["value_type"]; ok {
		if value, ok := filter.Value.(string); ok && strings.TrimSpace(value) != "" {
			q = q.Where(payloadpropertydefinition.ValueTypeContainsFold(strings.TrimSpace(value)))
		}
	}

	if filter, ok := opts.Filters["search"]; ok {
		if value, ok := filter.Value.(string); ok && strings.TrimSpace(value) != "" {
			searchValue := strings.TrimSpace(value)
			q = q.Where(payloadpropertydefinition.Or(
				payloadpropertydefinition.PayloadTypeContainsFold(searchValue),
				payloadpropertydefinition.PayloadVariantContainsFold(searchValue),
				payloadpropertydefinition.KeyContainsFold(searchValue),
				payloadpropertydefinition.DescriptionContainsFold(searchValue),
			))
		}
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, wrapListError(err, "định nghĩa thuộc tính payload")
	}

	if len(opts.Sort) > 0 {
		for _, sort := range opts.Sort {
			if !payloadPropertyDefinitionAllowedFields[sort.Field] {
				continue
			}
			if sort.Desc {
				q = q.Order(ent.Desc(sort.Field))
			} else {
				q = q.Order(ent.Asc(sort.Field))
			}
		}
	} else {
		q = q.Order(ent.Desc(payloadpropertydefinition.FieldCreatedAt))
	}

	items, err := q.Offset(offset).Limit(limit).All(ctx)
	if err != nil {
		return nil, 0, wrapListError(err, "định nghĩa thuộc tính payload")
	}

	return items, int64(total), nil
}

func (r *payloadPropertyDefinitionRepositoryImpl) ListPayloadTypes(ctx context.Context) ([]string, error) {
	types, err := r.client.PayloadPropertyDefinition.Query().
		Where(payloadpropertydefinition.PayloadTypeNEQ("")).
		GroupBy(payloadpropertydefinition.FieldPayloadType).
		Strings(ctx)
	if err != nil {
		return nil, wrapListError(err, "payload type")
	}

	for i := range types {
		types[i] = strings.TrimSpace(types[i])
	}
	sort.Strings(types)

	return types, nil
}

func (r *payloadPropertyDefinitionRepositoryImpl) DeleteAll(ctx context.Context) (int, error) {
	deletedCount, err := r.client.PayloadPropertyDefinition.Delete().Exec(ctx)
	if err != nil {
		return 0, wrapDeleteError(err, "định nghĩa thuộc tính payload")
	}

	return deletedCount, nil
}

func (r *payloadPropertyDefinitionRepositoryImpl) FindByPayloadTypeVariantAndKey(ctx context.Context, payloadType, payloadVariant, key string) (*ent.PayloadPropertyDefinition, error) {
	item, err := r.client.PayloadPropertyDefinition.Query().
		Where(
			payloadpropertydefinition.PayloadTypeEQ(strings.TrimSpace(payloadType)),
			payloadpropertydefinition.PayloadVariantEQ(strings.TrimSpace(payloadVariant)),
			payloadpropertydefinition.KeyEQ(strings.TrimSpace(key)),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy định nghĩa thuộc tính payload")
		}
		return nil, wrapFindError(err, "định nghĩa thuộc tính payload")
	}
	return item, nil
}

func (r *payloadPropertyDefinitionRepositoryImpl) UpsertByPayloadTypeVariantAndKey(ctx context.Context, e *ent.PayloadPropertyDefinition) (*ent.PayloadPropertyDefinition, bool, error) {
	payloadType := strings.TrimSpace(e.PayloadType)
	payloadVariant := strings.TrimSpace(e.PayloadVariant)
	key := strings.TrimSpace(e.Key)

	existing, err := r.client.PayloadPropertyDefinition.Query().
		Where(
			payloadpropertydefinition.PayloadTypeEQ(payloadType),
			payloadpropertydefinition.PayloadVariantEQ(payloadVariant),
			payloadpropertydefinition.KeyEQ(key),
		).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, false, wrapFindError(err, "định nghĩa thuộc tính payload")
	}

	if ent.IsNotFound(err) {
		created := &ent.PayloadPropertyDefinition{
			PayloadType:     payloadType,
			PayloadVariant:  payloadVariant,
			Key:             key,
			ValueType:       strings.TrimSpace(e.ValueType),
			ItemsType:       e.ItemsType,
			DefaultValue:    e.DefaultValue,
			EnumValues:      e.EnumValues,
			Title:           e.Title,
			Description:     strings.TrimSpace(e.Description),
			Presence:        e.Presence,
			Deprecated:      e.Deprecated,
			IsNested:        e.IsNested,
			NestedReference: e.NestedReference,
			ItemsReference:  e.ItemsReference,
			SupportedOs:     e.SupportedOs,
			Conditions:      e.Conditions,
			OrderIndex:      e.OrderIndex,
			YamlSourceFile:  e.YamlSourceFile,
		}
		if err := r.Create(ctx, created); err != nil {
			return nil, false, err
		}
		return created, true, nil
	}

	existing.PayloadType = payloadType
	existing.PayloadVariant = payloadVariant
	existing.Key = key
	existing.ValueType = strings.TrimSpace(e.ValueType)
	existing.ItemsType = e.ItemsType
	existing.DefaultValue = e.DefaultValue
	existing.EnumValues = e.EnumValues
	existing.Title = e.Title
	existing.Description = strings.TrimSpace(e.Description)
	existing.Presence = e.Presence
	existing.Deprecated = e.Deprecated
	existing.IsNested = e.IsNested
	existing.NestedReference = e.NestedReference
	existing.ItemsReference = e.ItemsReference
	existing.SupportedOs = e.SupportedOs
	existing.Conditions = e.Conditions
	existing.OrderIndex = e.OrderIndex
	existing.YamlSourceFile = e.YamlSourceFile
	if err := r.Update(ctx, existing); err != nil {
		return nil, false, err
	}
	return existing, false, nil
}
