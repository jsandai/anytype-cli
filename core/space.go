package core

import (
	"context"
	"fmt"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
)

func ApproveJoinRequest(spaceID, identity string, permissions model.ParticipantPermissions) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcSpaceRequestApproveRequest{
			SpaceId:     spaceID,
			Identity:    identity,
			Permissions: permissions,
		}
		_, err := client.SpaceRequestApprove(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to approve join request: %w", err)
		}
		return nil
	})
}

func JoinSpace(networkID, spaceID, inviteCID, inviteFileKey string) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcSpaceJoinRequest{
			NetworkId:     networkID,
			SpaceId:       spaceID,
			InviteCid:     inviteCID,
			InviteFileKey: inviteFileKey,
		}
		_, err := client.SpaceJoin(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to join space: %w", err)
		}
		return nil
	})
}

func LeaveSpace(spaceID string) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcSpaceDeleteRequest{
			SpaceId: spaceID,
		}
		_, err := client.SpaceDelete(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to leave space: %w", err)
		}
		return nil
	})
}

type SpaceInviteInfo struct {
	SpaceID           string
	SpaceName         string
	SpaceIconCID      string
	CreatorName       string
	IsGuestUserInvite bool
	InviteType        model.InviteType
}

func ViewSpaceInvite(inviteCID, inviteFileKey string) (*SpaceInviteInfo, error) {
	var info *SpaceInviteInfo
	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcSpaceInviteViewRequest{
			InviteCid:     inviteCID,
			InviteFileKey: inviteFileKey,
		}
		resp, err := client.SpaceInviteView(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to view space invite: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcSpaceInviteViewResponseError_NULL {
			return fmt.Errorf("space invite view error: %s", resp.Error.Description)
		}
		info = &SpaceInviteInfo{
			SpaceID:           resp.SpaceId,
			SpaceName:         resp.SpaceName,
			SpaceIconCID:      resp.SpaceIconCid,
			CreatorName:       resp.CreatorName,
			IsGuestUserInvite: resp.IsGuestUserInvite,
			InviteType:        resp.InviteType,
		}
		return nil
	})
	return info, err
}

type SpaceListItem struct {
	SpaceID string
	Name    string
	Status  model.SpaceStatus
}

// ListSpaces returns a list of all available spaces
func ListSpaces() ([]SpaceListItem, error) {
	techSpaceID, err := GetStoredTechSpaceID()
	if err != nil {
		return nil, fmt.Errorf("tech space ID not found in config - please login first: %w", err)
	}

	var spaces []SpaceListItem
	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcObjectSearchRequest{
			SpaceId: techSpaceID,
			Filters: []*model.BlockContentDataviewFilter{
				{
					RelationKey: bundle.RelationKeyResolvedLayout.String(),
					Condition:   model.BlockContentDataviewFilter_Equal,
					Value:       pbtypes.Int64(int64(model.ObjectType_spaceView)),
				},
				{
					RelationKey: bundle.RelationKeySpaceLocalStatus.String(),
					Condition:   model.BlockContentDataviewFilter_In,
					Value:       pbtypes.IntList(int(model.SpaceStatus_Unknown), int(model.SpaceStatus_Ok)),
				},
				{
					RelationKey: bundle.RelationKeySpaceAccountStatus.String(),
					Condition:   model.BlockContentDataviewFilter_In,
					Value:       pbtypes.IntList(int(model.SpaceStatus_Unknown), int(model.SpaceStatus_SpaceActive)),
				},
			},
			Sorts: []*model.BlockContentDataviewSort{
				{
					RelationKey:    bundle.RelationKeySpaceOrder.String(),
					Type:           model.BlockContentDataviewSort_Asc,
					NoCollate:      true,
					EmptyPlacement: model.BlockContentDataviewSort_End,
				},
			},
			Keys: []string{
				bundle.RelationKeyTargetSpaceId.String(),
				bundle.RelationKeyName.String(),
				bundle.RelationKeySpaceLocalStatus.String(),
			},
		}

		resp, err := client.ObjectSearch(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to search spaces: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcObjectSearchResponseError_NULL {
			return fmt.Errorf("object search error: %s", resp.Error.Description)
		}

		for _, record := range resp.Records {
			item := SpaceListItem{}

			// Get space ID
			if spaceIDVal := pbtypes.GetString(record, bundle.RelationKeyTargetSpaceId.String()); spaceIDVal != "" {
				item.SpaceID = spaceIDVal
			}

			// Get name
			if nameVal := pbtypes.GetString(record, bundle.RelationKeyName.String()); nameVal != "" {
				item.Name = nameVal
			}

			// Get status
			if statusVal := pbtypes.GetInt64(record, bundle.RelationKeySpaceLocalStatus.String()); statusVal != 0 {
				item.Status = model.SpaceStatus(statusVal)
			}

			if item.SpaceID != "" {
				spaces = append(spaces, item)
			}
		}

		return nil
	})

	return spaces, err
}
