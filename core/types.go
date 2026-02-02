package core

import (
	"context"
	"fmt"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/types"
)

// TypeInfo represents an object type
type TypeInfo struct {
	ID          string
	UniqueKey   string
	Name        string
	Description string
	IconEmoji   string
}

// CreateTypeRequest contains parameters for creating a custom object type
type CreateTypeRequest struct {
	SpaceID     string
	Name        string
	PluralName  string
	Description string
	IconEmoji   string
}

// CreateObjectType creates a custom object type
func CreateObjectType(req CreateTypeRequest) (*TypeInfo, error) {
	var result *TypeInfo

	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		details := &types.Struct{
			Fields: map[string]*types.Value{
				"name": pbtypes.String(req.Name),
			},
		}

		if req.PluralName != "" {
			details.Fields["pluralName"] = pbtypes.String(req.PluralName)
		}
		if req.Description != "" {
			details.Fields["description"] = pbtypes.String(req.Description)
		}
		if req.IconEmoji != "" {
			details.Fields["iconEmoji"] = pbtypes.String(req.IconEmoji)
		}

		rpcReq := &pb.RpcObjectCreateObjectTypeRequest{
			SpaceId: req.SpaceID,
			Details: details,
		}

		resp, err := client.ObjectCreateObjectType(ctx, rpcReq)
		if err != nil {
			return fmt.Errorf("failed to create type: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcObjectCreateObjectTypeResponseError_NULL {
			return fmt.Errorf("create type error: %s", resp.Error.Description)
		}

		result = &TypeInfo{
			ID:        resp.ObjectId,
			UniqueKey: pbtypes.GetString(resp.Details, "uniqueKey"),
			Name:      req.Name,
		}

		return nil
	})

	return result, err
}

// ListObjectTypes lists all object types in a space
func ListObjectTypes(spaceID string) ([]TypeInfo, error) {
	var results []TypeInfo

	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		// Search for objects of type "objectType"
		filters := []*model.BlockContentDataviewFilter{
			{
				RelationKey: "layout",
				Condition:   model.BlockContentDataviewFilter_Equal,
				Value:       pbtypes.Int64(int64(model.ObjectType_objectType)),
			},
		}

		req := &pb.RpcObjectSearchRequest{
			SpaceId: spaceID,
			Filters: filters,
			Keys:    []string{"id", "uniqueKey", "name", "description", "iconEmoji"},
			Limit:   100,
		}

		resp, err := client.ObjectSearch(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to list types: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcObjectSearchResponseError_NULL {
			return fmt.Errorf("list types error: %s", resp.Error.Description)
		}

		for _, record := range resp.Records {
			results = append(results, TypeInfo{
				ID:          pbtypes.GetString(record, "id"),
				UniqueKey:   pbtypes.GetString(record, "uniqueKey"),
				Name:        pbtypes.GetString(record, "name"),
				Description: pbtypes.GetString(record, "description"),
				IconEmoji:   pbtypes.GetString(record, "iconEmoji"),
			})
		}

		return nil
	})

	return results, err
}

// RelationInfo represents a relation (field/property)
type RelationInfo struct {
	ID        string
	UniqueKey string
	Name      string
	Format    string // text, number, select, multi-select, date, etc.
}

// RelationFormat maps format names to protobuf values
var RelationFormats = map[string]model.RelationFormat{
	"text":         model.RelationFormat_longtext,
	"number":       model.RelationFormat_number,
	"select":       model.RelationFormat_status,
	"multi-select": model.RelationFormat_tag,
	"date":         model.RelationFormat_date,
	"checkbox":     model.RelationFormat_checkbox,
	"url":          model.RelationFormat_url,
	"email":        model.RelationFormat_email,
	"phone":        model.RelationFormat_phone,
	"object":       model.RelationFormat_object,
	"file":         model.RelationFormat_file,
}

// CreateRelationRequest contains parameters for creating a relation
type CreateRelationRequest struct {
	SpaceID     string
	Name        string
	Format      string // text, number, select, etc.
	Description string
}

// CreateRelation creates a custom relation
func CreateRelation(req CreateRelationRequest) (*RelationInfo, error) {
	var result *RelationInfo

	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		format, ok := RelationFormats[req.Format]
		if !ok {
			return fmt.Errorf("unknown relation format: %s (valid: text, number, select, multi-select, date, checkbox, url, email, phone, object, file)", req.Format)
		}

		details := &types.Struct{
			Fields: map[string]*types.Value{
				"name":           pbtypes.String(req.Name),
				"relationFormat": pbtypes.Float64(float64(format)),
			},
		}

		if req.Description != "" {
			details.Fields["description"] = pbtypes.String(req.Description)
		}

		rpcReq := &pb.RpcObjectCreateRelationRequest{
			SpaceId: req.SpaceID,
			Details: details,
		}

		resp, err := client.ObjectCreateRelation(ctx, rpcReq)
		if err != nil {
			return fmt.Errorf("failed to create relation: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcObjectCreateRelationResponseError_NULL {
			return fmt.Errorf("create relation error: %s", resp.Error.Description)
		}

		result = &RelationInfo{
			ID:        resp.ObjectId,
			UniqueKey: resp.Key,
			Name:      req.Name,
			Format:    req.Format,
		}

		return nil
	})

	return result, err
}

// ListRelations lists all relations in a space
func ListRelations(spaceID string) ([]RelationInfo, error) {
	var results []RelationInfo

	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		// Search for objects of type "relation"
		filters := []*model.BlockContentDataviewFilter{
			{
				RelationKey: "layout",
				Condition:   model.BlockContentDataviewFilter_Equal,
				Value:       pbtypes.Int64(int64(model.ObjectType_relation)),
			},
		}

		req := &pb.RpcObjectSearchRequest{
			SpaceId: spaceID,
			Filters: filters,
			Keys:    []string{"id", "uniqueKey", "name", "relationKey", "relationFormat"},
			Limit:   200,
		}

		resp, err := client.ObjectSearch(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to list relations: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcObjectSearchResponseError_NULL {
			return fmt.Errorf("list relations error: %s", resp.Error.Description)
		}

		// Build reverse format map
		formatNames := make(map[model.RelationFormat]string)
		for name, val := range RelationFormats {
			formatNames[val] = name
		}

		for _, record := range resp.Records {
			formatVal := model.RelationFormat(pbtypes.GetInt64(record, "relationFormat"))
			formatName := formatNames[formatVal]
			if formatName == "" {
				formatName = fmt.Sprintf("unknown(%d)", formatVal)
			}

			key := pbtypes.GetString(record, "relationKey")
			if key == "" {
				key = pbtypes.GetString(record, "uniqueKey")
			}

			results = append(results, RelationInfo{
				ID:        pbtypes.GetString(record, "id"),
				UniqueKey: key,
				Name:      pbtypes.GetString(record, "name"),
				Format:    formatName,
			})
		}

		return nil
	})

	return results, err
}

// AddRelationToType adds a relation to an object type's recommended relations
func AddRelationToType(spaceID string, typeID string, relationKey string) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcObjectTypeRelationAddRequest{
			ObjectTypeUrl: typeID,
			RelationKeys:  []string{relationKey},
		}

		resp, err := client.ObjectTypeRelationAdd(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to add relation to type: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcObjectTypeRelationAddResponseError_NULL {
			return fmt.Errorf("add relation error: %s", resp.Error.Description)
		}

		return nil
	})
}
