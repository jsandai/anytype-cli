package core

import (
	"context"
	"fmt"
	"time"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/types"
)

// ObjectInfo represents basic object information for display
type ObjectInfo struct {
	ID        string
	Name      string
	Type      string
	SpaceID   string
	CreatedAt time.Time
	Snippet   string
}

// ObjectDetail represents full object details
type ObjectDetail struct {
	ID        string
	Name      string
	Type      string
	TypeID    string
	SpaceID   string
	CreatedAt time.Time
	Details   map[string]interface{}
}

// CreateObjectRequest contains parameters for creating an object
type CreateObjectRequest struct {
	SpaceID       string
	TypeID        string            // Object type unique key (e.g., "ot-note", "ot-page", or custom type ID)
	Name          string
	Details       map[string]interface{}
}

// CreateObject creates a new object in the specified space
func CreateObject(req CreateObjectRequest) (*ObjectInfo, error) {
	var result *ObjectInfo

	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		// Build details struct
		details := &types.Struct{
			Fields: make(map[string]*types.Value),
		}

		// Add name if provided
		if req.Name != "" {
			details.Fields["name"] = pbtypes.String(req.Name)
		}

		// Add any additional details
		for key, val := range req.Details {
			switch v := val.(type) {
			case string:
				details.Fields[key] = pbtypes.String(v)
			case float64:
				details.Fields[key] = pbtypes.Float64(v)
			case bool:
				details.Fields[key] = pbtypes.Bool(v)
			}
		}

		rpcReq := &pb.RpcObjectCreateRequest{
			SpaceId:            req.SpaceID,
			ObjectTypeUniqueKey: req.TypeID,
			Details:            details,
		}

		resp, err := client.ObjectCreate(ctx, rpcReq)
		if err != nil {
			return fmt.Errorf("failed to create object: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcObjectCreateResponseError_NULL {
			return fmt.Errorf("create object error: %s", resp.Error.Description)
		}

		result = &ObjectInfo{
			ID:      resp.ObjectId,
			Name:    req.Name,
			SpaceID: req.SpaceID,
		}

		return nil
	})

	return result, err
}

// SearchObjectsRequest contains parameters for searching objects
type SearchObjectsRequest struct {
	SpaceID string
	Query   string   // Full-text search query
	Types   []string // Filter by object type IDs
	Limit   int32
}

// SearchObjects searches for objects in a space
func SearchObjects(req SearchObjectsRequest) ([]ObjectInfo, error) {
	var results []ObjectInfo

	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		// Build filters
		var filters []*model.BlockContentDataviewFilter

		// Add type filter if specified
		if len(req.Types) > 0 {
			filters = append(filters, &model.BlockContentDataviewFilter{
				RelationKey: "type",
				Condition:   model.BlockContentDataviewFilter_In,
				Value:       pbtypes.StringList(req.Types),
			})
		}

		rpcReq := &pb.RpcObjectSearchRequest{
			SpaceId:   req.SpaceID,
			Filters:   filters,
			FullText:  req.Query,
			Keys:      []string{"id", "name", "type", "spaceId", "createdDate", "snippet"},
			Limit:     req.Limit,
		}

		resp, err := client.ObjectSearch(ctx, rpcReq)
		if err != nil {
			return fmt.Errorf("failed to search objects: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcObjectSearchResponseError_NULL {
			return fmt.Errorf("search error: %s", resp.Error.Description)
		}

		for _, record := range resp.Records {
			createdDate := pbtypes.GetInt64(record, "createdDate")
			results = append(results, ObjectInfo{
				ID:        pbtypes.GetString(record, "id"),
				Name:      pbtypes.GetString(record, "name"),
				Type:      pbtypes.GetString(record, "type"),
				SpaceID:   pbtypes.GetString(record, "spaceId"),
				CreatedAt: time.Unix(createdDate, 0),
				Snippet:   pbtypes.GetString(record, "snippet"),
			})
		}

		return nil
	})

	return results, err
}

// GetObject retrieves full details of an object
func GetObject(spaceID string, objectID string) (*ObjectDetail, error) {
	var result *ObjectDetail

	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		// First open the object to get its details
		openReq := &pb.RpcObjectOpenRequest{
			SpaceId:  spaceID,
			ObjectId: objectID,
		}

		openResp, err := client.ObjectOpen(ctx, openReq)
		if err != nil {
			return fmt.Errorf("failed to open object: %w", err)
		}
		if openResp.Error != nil && openResp.Error.Code != pb.RpcObjectOpenResponseError_NULL {
			return fmt.Errorf("open object error: %s", openResp.Error.Description)
		}

		// Extract details from the response
		if openResp.ObjectView != nil && openResp.ObjectView.Details != nil {
			for _, detail := range openResp.ObjectView.Details {
				if detail.Id == objectID {
					result = &ObjectDetail{
						ID:      objectID,
						SpaceID: spaceID,
						Details: make(map[string]interface{}),
					}

					if detail.Details != nil && detail.Details.Fields != nil {
						for key, val := range detail.Details.Fields {
							result.Details[key] = extractValue(val)
						}
						if name, ok := result.Details["name"].(string); ok {
							result.Name = name
						}
						if typeVal, ok := result.Details["type"].(string); ok {
							result.Type = typeVal
						}
					}
					break
				}
			}
		}

		// Close the object
		closeReq := &pb.RpcObjectCloseRequest{
			SpaceId:  spaceID,
			ObjectId: objectID,
		}
		_, _ = client.ObjectClose(ctx, closeReq)

		return nil
	})

	return result, err
}

// UpdateObjectDetails updates an object's details/relations
func UpdateObjectDetails(spaceID string, objectID string, details map[string]interface{}) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		// Build details list
		var detailsList []*model.Detail

		for key, val := range details {
			var pbVal *types.Value
			switch v := val.(type) {
			case string:
				pbVal = pbtypes.String(v)
			case float64:
				pbVal = pbtypes.Float64(v)
			case bool:
				pbVal = pbtypes.Bool(v)
			case nil:
				pbVal = pbtypes.Null()
			default:
				continue
			}

			detailsList = append(detailsList, &model.Detail{
				Key:   key,
				Value: pbVal,
			})
		}

		req := &pb.RpcObjectSetDetailsRequest{
			ContextId: objectID,
			Details:   detailsList,
		}

		resp, err := client.ObjectSetDetails(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to update object: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcObjectSetDetailsResponseError_NULL {
			return fmt.Errorf("update error: %s", resp.Error.Description)
		}

		return nil
	})
}

// DeleteObjects deletes one or more objects
func DeleteObjects(objectIDs []string) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcObjectListDeleteRequest{
			ObjectIds: objectIDs,
		}

		resp, err := client.ObjectListDelete(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to delete objects: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcObjectListDeleteResponseError_NULL {
			return fmt.Errorf("delete error: %s", resp.Error.Description)
		}

		return nil
	})
}

// extractValue converts a protobuf Value to a Go interface{}
func extractValue(val *types.Value) interface{} {
	if val == nil {
		return nil
	}
	switch v := val.Kind.(type) {
	case *types.Value_StringValue:
		return v.StringValue
	case *types.Value_NumberValue:
		return v.NumberValue
	case *types.Value_BoolValue:
		return v.BoolValue
	case *types.Value_NullValue:
		return nil
	case *types.Value_ListValue:
		var list []interface{}
		for _, item := range v.ListValue.Values {
			list = append(list, extractValue(item))
		}
		return list
	case *types.Value_StructValue:
		m := make(map[string]interface{})
		for k, val := range v.StructValue.Fields {
			m[k] = extractValue(val)
		}
		return m
	default:
		return nil
	}
}
