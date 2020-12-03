/**
 * @Author: jasonruan
 * @Date:   2020-12-03 11:31:33
 **/
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"context"
	"github.com/golang/protobuf/proto"
	"io"
)

func (cc ChainClient) SubscribeBlock(ctx context.Context, startBlock, endBlock int64, withRwSet bool) (<-chan interface{}, error) {
	payloadBytes, err := constructSubscribeBlockPayload(startBlock, endBlock, withRwSet)
	if err != nil {
		return nil, err
	}

	return cc.Subscribe(ctx, pb.TxType_SUBSCRIBE_BLOCK_INFO, payloadBytes)
}

func (cc ChainClient) SubscribeTx(ctx context.Context, startBlock, endBlock int64, txType pb.TxType, txIds []string) (<-chan interface{}, error) {
	payloadBytes, err := constructSubscribeTxPayload(startBlock, endBlock, txType, txIds)
	if err != nil {
		return nil, err
	}

	return cc.Subscribe(ctx, pb.TxType_SUBSCRIBE_TX_INFO, payloadBytes)
}

func (cc ChainClient) Subscribe(ctx context.Context, txType pb.TxType, payloadBytes []byte) (<-chan interface{}, error) {
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
				switch txType{
				case pb.TxType_SUBSCRIBE_BLOCK_INFO:
					blockInfo := &pb.BlockInfo{}
					if err = proto.Unmarshal(result.Data, blockInfo); err != nil {
						cc.logger.Error("[SDK] Subscriber receive block failed, %s", err)
						close(c)
						return
					}

					ret = blockInfo
				case pb.TxType_SUBSCRIBE_TX_INFO:
					tx := &pb.Transaction{}
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
