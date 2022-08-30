// Code generated by goctl. DO NOT EDIT!
// Source: chat.proto

package chatclient

import (
	"context"

	"im-center/service/business/chatService/rpc/chat"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AckMsgReq                    = chat.AckMsgReq
	AddFriendReq                 = chat.AddFriendReq
	At                           = chat.At
	BatchChangeFriendRelationReq = chat.BatchChangeFriendRelationReq
	BatchDelMsgReq               = chat.BatchDelMsgReq
	FriendGetsResp               = chat.FriendGetsResp
	FriendItem                   = chat.FriendItem
	GetsQueryItem                = chat.GetsQueryItem
	GetsReq                      = chat.GetsReq
	GroupAddReq                  = chat.GroupAddReq
	GroupBatchDelReq             = chat.GroupBatchDelReq
	GroupGetReq                  = chat.GroupGetReq
	GroupGetsResp                = chat.GroupGetsResp
	GroupItem                    = chat.GroupItem
	GroupMsgItem                 = chat.GroupMsgItem
	GroupUpdateReq               = chat.GroupUpdateReq
	NullResp                     = chat.NullResp
	SendManyMsgReq               = chat.SendManyMsgReq
	SendOneMsgReq                = chat.SendOneMsgReq
	SingleMsgItem                = chat.SingleMsgItem
	SyncMsgReq                   = chat.SyncMsgReq
	SyncMsgResp                  = chat.SyncMsgResp
	UserAddGroupReq              = chat.UserAddGroupReq
	UserExitGroupReq             = chat.UserExitGroupReq
	WithdrawMsgReq               = chat.WithdrawMsgReq

	Chat interface {
		AddFriend(ctx context.Context, in *AddFriendReq, opts ...grpc.CallOption) (*NullResp, error)
		BatchChangeFriendRelation(ctx context.Context, in *BatchChangeFriendRelationReq, opts ...grpc.CallOption) (*NullResp, error)
		FriendGets(ctx context.Context, in *GetsReq, opts ...grpc.CallOption) (*FriendGetsResp, error)
		SendOneMsg(ctx context.Context, in *SendOneMsgReq, opts ...grpc.CallOption) (*NullResp, error)
		SendManyMsg(ctx context.Context, in *SendManyMsgReq, opts ...grpc.CallOption) (*NullResp, error)
		AckMsg(ctx context.Context, in *AckMsgReq, opts ...grpc.CallOption) (*NullResp, error)
		SyncMsg(ctx context.Context, in *SyncMsgReq, opts ...grpc.CallOption) (*SyncMsgResp, error)
		WithdrawMsg(ctx context.Context, in *WithdrawMsgReq, opts ...grpc.CallOption) (*NullResp, error)
		BatchDelMsg(ctx context.Context, in *BatchDelMsgReq, opts ...grpc.CallOption) (*NullResp, error)
		//  --------------------------------------------------------------------------------------------------------
		GroupGet(ctx context.Context, in *GroupGetReq, opts ...grpc.CallOption) (*GroupItem, error)
		GroupGets(ctx context.Context, in *GetsReq, opts ...grpc.CallOption) (*GroupGetsResp, error)
		GroupAdd(ctx context.Context, in *GroupAddReq, opts ...grpc.CallOption) (*NullResp, error)
		GroupUpdate(ctx context.Context, in *GroupUpdateReq, opts ...grpc.CallOption) (*NullResp, error)
		GroupBatchDel(ctx context.Context, in *GroupBatchDelReq, opts ...grpc.CallOption) (*NullResp, error)
		//  用户群组关系
		UserAddGroup(ctx context.Context, in *UserAddGroupReq, opts ...grpc.CallOption) (*NullResp, error)
		UserExitGroup(ctx context.Context, in *UserExitGroupReq, opts ...grpc.CallOption) (*NullResp, error)
	}

	defaultChat struct {
		cli zrpc.Client
	}
)

func NewChat(cli zrpc.Client) Chat {
	return &defaultChat{
		cli: cli,
	}
}

func (m *defaultChat) AddFriend(ctx context.Context, in *AddFriendReq, opts ...grpc.CallOption) (*NullResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.AddFriend(ctx, in, opts...)
}

func (m *defaultChat) BatchChangeFriendRelation(ctx context.Context, in *BatchChangeFriendRelationReq, opts ...grpc.CallOption) (*NullResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.BatchChangeFriendRelation(ctx, in, opts...)
}

func (m *defaultChat) FriendGets(ctx context.Context, in *GetsReq, opts ...grpc.CallOption) (*FriendGetsResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.FriendGets(ctx, in, opts...)
}

func (m *defaultChat) SendOneMsg(ctx context.Context, in *SendOneMsgReq, opts ...grpc.CallOption) (*NullResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.SendOneMsg(ctx, in, opts...)
}

func (m *defaultChat) SendManyMsg(ctx context.Context, in *SendManyMsgReq, opts ...grpc.CallOption) (*NullResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.SendManyMsg(ctx, in, opts...)
}

func (m *defaultChat) AckMsg(ctx context.Context, in *AckMsgReq, opts ...grpc.CallOption) (*NullResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.AckMsg(ctx, in, opts...)
}

func (m *defaultChat) SyncMsg(ctx context.Context, in *SyncMsgReq, opts ...grpc.CallOption) (*SyncMsgResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.SyncMsg(ctx, in, opts...)
}

func (m *defaultChat) WithdrawMsg(ctx context.Context, in *WithdrawMsgReq, opts ...grpc.CallOption) (*NullResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.WithdrawMsg(ctx, in, opts...)
}

func (m *defaultChat) BatchDelMsg(ctx context.Context, in *BatchDelMsgReq, opts ...grpc.CallOption) (*NullResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.BatchDelMsg(ctx, in, opts...)
}

//  --------------------------------------------------------------------------------------------------------
func (m *defaultChat) GroupGet(ctx context.Context, in *GroupGetReq, opts ...grpc.CallOption) (*GroupItem, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.GroupGet(ctx, in, opts...)
}

func (m *defaultChat) GroupGets(ctx context.Context, in *GetsReq, opts ...grpc.CallOption) (*GroupGetsResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.GroupGets(ctx, in, opts...)
}

func (m *defaultChat) GroupAdd(ctx context.Context, in *GroupAddReq, opts ...grpc.CallOption) (*NullResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.GroupAdd(ctx, in, opts...)
}

func (m *defaultChat) GroupUpdate(ctx context.Context, in *GroupUpdateReq, opts ...grpc.CallOption) (*NullResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.GroupUpdate(ctx, in, opts...)
}

func (m *defaultChat) GroupBatchDel(ctx context.Context, in *GroupBatchDelReq, opts ...grpc.CallOption) (*NullResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.GroupBatchDel(ctx, in, opts...)
}

//  用户群组关系
func (m *defaultChat) UserAddGroup(ctx context.Context, in *UserAddGroupReq, opts ...grpc.CallOption) (*NullResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.UserAddGroup(ctx, in, opts...)
}

func (m *defaultChat) UserExitGroup(ctx context.Context, in *UserExitGroupReq, opts ...grpc.CallOption) (*NullResp, error) {
	client := chat.NewChatClient(m.cli.Conn())
	return client.UserExitGroup(ctx, in, opts...)
}
