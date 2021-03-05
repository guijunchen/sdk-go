/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/pb/common"
	"context"
	"github.com/golang/protobuf/proto"
	"io"
)

func (cc ChainClient) SubscribeBlock(ctx context.Context, startBlock, endBlock int64, withRwSet bool) (<-chan interface{}, error) {
	payloadBytes, err := constructSubscribeBlockPayload(startBlock, endBlock, withRwSet)
	if err != nil {
		return nil, err
	}

	return cc.Subscribe(ctx, common.TxType_SUBSCRIBE_BLOCK_INFO, payloadBytes)
}

func (cc ChainClient) SubscribeTx(ctx context.Context, startBlock, endBlock int64, txType common.TxType, txIds []string) (<-chan interface{}, error) {
	payloadBytes, err := constructSubscribeTxPayload(startBlock, endBlock, txType, txIds)
	if err != nil {
		return nil, err
	}

	return cc.Subscribe(ctx, common.TxType_SUBSCRIBE_TX_INFO, payloadBytes)
}

func (cc ChainClient) Subscribe(ctx context.Context, txType common.TxType, payloadBytes []byte) (<-chan interface{}, error) {
	txId := GetRandTxId()

	req, err := cc.generateTxRequest(txId, txType, payloadBytes)
	if err != nil {
		return nil, err
	}

	client, err := cc.pool.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.rpcNode.Subscribe(ctx, req)
	if err != nil {
		return nil, err
	}

	c := make(chan interface{})
	go func() {
		defer close(c)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				result, err := resp.Recv()
				if err == io.EOF {
					cc.logger.Debugf("[SDK] Subscirber got EOF and stop recv msg")
					return
				}

				if err != nil {
					cc.logger.Errorf("[SDK] Subscriber receive failed, %s", err)
					return
				}

				var ret interface{}
				switch txType {
				case common.TxType_SUBSCRIBE_BLOCK_INFO:
					blockInfo := &common.BlockInfo{}
					if err = proto.Unmarshal(result.Data, blockInfo); err != nil {
						cc.logger.Error("[SDK] Subscriber receive block failed, %s", err)
						close(c)
						return
					}

					ret = blockInfo
				case common.TxType_SUBSCRIBE_TX_INFO:
					tx := &common.Transaction{}
					if err = proto.Unmarshal(result.Data, tx); err != nil {
						cc.logger.Error("[SDK] Subscriber receive tx failed, %s", err)
						close(c)
						return
					}

					ret = tx
				default:
					ret = result.Data
				}

				c <- ret
			}
		}
	}()

	return c, nil
}
