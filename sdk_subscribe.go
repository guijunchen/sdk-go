/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/api"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/gogo/protobuf/proto"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

func (cc *ChainClient) CreateSubscribeBlockPayload(startBlock, endBlock int64,
	withRWSet, onlyHeader bool) *common.Payload {

	return cc.CreatePayload("", common.TxType_SUBSCRIBE, syscontract.SystemContract_SUBSCRIBE_MANAGE.String(),
		syscontract.SubscribeFunction_SUBSCRIBE_BLOCK.String(), []*common.KeyValuePair{
			{
				Key:   syscontract.SubscribeBlock_START_BLOCK.String(),
				Value: utils.I64ToBytes(startBlock),
			},
			{
				Key:   syscontract.SubscribeBlock_END_BLOCK.String(),
				Value: utils.I64ToBytes(endBlock),
			},
			{
				Key:   syscontract.SubscribeBlock_WITH_RWSET.String(),
				Value: []byte(strconv.FormatBool(withRWSet)),
			},
			{
				Key:   syscontract.SubscribeBlock_ONLY_HEADER.String(),
				Value: []byte(strconv.FormatBool(onlyHeader)),
			},
		}, defaultSeq,
	)
}

func (cc *ChainClient) CreateSubscribeTxPayload(startBlock, endBlock int64,
	contractName string, txIds []string) *common.Payload {

	return cc.CreatePayload("", common.TxType_SUBSCRIBE, syscontract.SystemContract_SUBSCRIBE_MANAGE.String(),
		syscontract.SubscribeFunction_SUBSCRIBE_TX.String(), []*common.KeyValuePair{
			{
				Key:   syscontract.SubscribeTx_START_BLOCK.String(),
				Value: utils.I64ToBytes(startBlock),
			},
			{
				Key:   syscontract.SubscribeTx_END_BLOCK.String(),
				Value: utils.I64ToBytes(endBlock),
			},
			{
				Key:   syscontract.SubscribeTx_CONTRACT_NAME.String(),
				Value: []byte(contractName),
			},
			{
				Key:   syscontract.SubscribeTx_TX_IDS.String(),
				Value: []byte(strings.Join(txIds, ",")),
			},
		}, defaultSeq,
	)
}

func (cc *ChainClient) SubscribeBlock(ctx context.Context, startBlock, endBlock int64, withRWSet,
	onlyHeader bool) (<-chan interface{}, error) {

	payload := cc.CreateSubscribeBlockPayload(startBlock, endBlock, withRWSet, onlyHeader)

	return cc.Subscribe(ctx, payload)
}

func (cc *ChainClient) SubscribeTx(ctx context.Context, startBlock, endBlock int64, contractName string,
	txIds []string) (<-chan interface{}, error) {

	payload := cc.CreateSubscribeTxPayload(startBlock, endBlock, contractName, txIds)

	return cc.Subscribe(ctx, payload)
}

func (cc *ChainClient) SubscribeContractEvent(ctx context.Context, startBlock, endBlock int64,
	contractName, topic string) (<-chan interface{}, error) {

	payload := cc.CreatePayload("", common.TxType_SUBSCRIBE, syscontract.SystemContract_SUBSCRIBE_MANAGE.String(),
		syscontract.SubscribeFunction_SUBSCRIBE_CONTRACT_EVENT.String(), []*common.KeyValuePair{
			{
				Key:   syscontract.SubscribeContractEvent_START_BLOCK.String(),
				Value: utils.I64ToBytes(startBlock),
			},
			{
				Key:   syscontract.SubscribeContractEvent_END_BLOCK.String(),
				Value: utils.I64ToBytes(endBlock),
			},
			{
				Key:   syscontract.SubscribeContractEvent_CONTRACT_NAME.String(),
				Value: []byte(contractName),
			},
			{
				Key:   syscontract.SubscribeContractEvent_TOPIC.String(),
				Value: []byte(topic),
			},
		}, defaultSeq,
	)

	return cc.Subscribe(ctx, payload)
}

func (cc *ChainClient) Subscribe(ctx context.Context, payload *common.Payload) (<-chan interface{}, error) {
	// get stream first time, throw immediately when an error occurs
	stream, err := cc.getSubscribeStream(ctx, payload)
	if err != nil {
		cc.logger.Error(err)
		return nil, err
	}

	var (
		dataC             = make(chan interface{})
		reconnectC, doneC = make(chan struct{}, 1), make(chan struct{}, 1)
	)

	go func() {
		defer func() {
			close(dataC)
			close(reconnectC)
			close(doneC)
		}()

		go cc.subscribe(ctx, payload, dataC, reconnectC, doneC, stream)
		for {
			select {
			case <-reconnectC:
				cc.logger.Debug("[SDK] Subscriber reconnecting...")
				// always get a new stream, set payload.Timestamp to the current time
				payload.Timestamp = time.Now().Unix()
				stream, err := cc.getSubscribeStream(ctx, payload)
				if err != nil {
					cc.logger.Error(err)
					return
				}
				go cc.subscribe(ctx, payload, dataC, reconnectC, doneC, stream)
			case <-doneC:
				cc.logger.Debug("[SDK] Subscriber done")
				return
			}
		}
	}()

	return dataC, nil
}

func (cc *ChainClient) subscribe(ctx context.Context, payload *common.Payload, dataC chan interface{},
	reconnectC, doneC chan struct{}, stream api.RpcNode_SubscribeClient) {

	for {
		select {
		case <-ctx.Done():
			doneC <- struct{}{}
			return
		default:
			result, err := stream.Recv()
			if err == io.EOF {
				cc.logger.Debugf("[SDK] Subscriber got EOF and stop recv msg")
				doneC <- struct{}{}
				return
			}

			if err != nil {
				cc.logger.Errorf("[SDK] Subscriber receive failed, %s", err)
				rpcStatus, ok := grpcstatus.FromError(err)
				if !ok {
					doneC <- struct{}{}
					return
				}

				if rpcStatus.Code() != grpccodes.Unavailable {
					doneC <- struct{}{}
					return
				}

				reconnectC <- struct{}{}
				return
			}

			var ret interface{}
			switch payload.Method {
			case syscontract.SubscribeFunction_SUBSCRIBE_BLOCK.String():
				blockInfo := &common.BlockInfo{}
				if err = proto.Unmarshal(result.Data, blockInfo); err == nil {
					ret = blockInfo
					break
				}

				blockHeader := &common.BlockHeader{}
				if err = proto.Unmarshal(result.Data, blockHeader); err == nil {
					ret = blockHeader
					break
				}

				cc.logger.Error("[SDK] Subscriber receive block failed, %s", err)
				doneC <- struct{}{}
				return
			case syscontract.SubscribeFunction_SUBSCRIBE_TX.String():
				tx := &common.Transaction{}
				if err = proto.Unmarshal(result.Data, tx); err != nil {
					cc.logger.Error("[SDK] Subscriber receive tx failed, %s", err)
					doneC <- struct{}{}
					return
				}
				ret = tx
			case syscontract.SubscribeFunction_SUBSCRIBE_CONTRACT_EVENT.String():
				events := &common.ContractEventInfoList{}
				if err = proto.Unmarshal(result.Data, events); err != nil {
					cc.logger.Error("[SDK] Subscriber receive contract event failed, %s", err)
					doneC <- struct{}{}
					return
				}
				for _, event := range events.ContractEvents {
					dataC <- event
				}
				continue

			default:
				ret = result.Data
			}

			dataC <- ret
		}
	}
}

func (cc *ChainClient) getSubscribeStream(ctx context.Context,
	payload *common.Payload) (api.RpcNode_SubscribeClient, error) {

	req, err := cc.GenerateTxRequest(payload, nil)
	if err != nil {
		return nil, err
	}

	networkCli, err := cc.pool.getClient()
	if err != nil {
		return nil, err
	}

	return networkCli.rpcNode.Subscribe(ctx, req)
}
