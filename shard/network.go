package shard

import (
	"ai-blockchain-shard/core"
	"fmt"
)

// Network 分片网络
type Network struct {
	Shards map[uint64][]*Node // 分片ID到节点列表的映射
}

// NewNetwork 创建新的分片网络
func NewNetwork() *Network {
	return &Network{
		Shards: make(map[uint64][]*Node),
	}
}

// AddNode 添加节点到网络
func (n *Network) AddNode(node *Node) {
	shardID := node.ShardID
	n.Shards[shardID] = append(n.Shards[shardID], node)
	fmt.Printf("Added node %d to shard %d\n", node.ID, shardID)
}

// GetShardNodes 获取指定分片的节点
func (n *Network) GetShardNodes(shardID uint64) []*Node {
	return n.Shards[shardID]
}

// GetAllShards 获取所有分片ID
func (n *Network) GetAllShards() []uint64 {
	var shards []uint64
	for shardID := range n.Shards {
		shards = append(shards, shardID)
	}
	return shards
}

// BroadcastTransaction 广播交易到指定分片
func (n *Network) BroadcastTransaction(tx *core.Transaction, shardID uint64) {
	nodes := n.Shards[shardID]
	if nodes == nil || len(nodes) == 0 {
		fmt.Printf("No nodes found in shard %d\n", shardID)
		return
	}
	
	for _, node := range nodes {
		node.AddTransactions([]*core.Transaction{tx})
	}
	
	fmt.Printf("Broadcasted transaction to shard %d\n", shardID)
}

// PrintNetworkInfo 打印网络信息
func (n *Network) PrintNetworkInfo() {
	fmt.Println("=== 分片网络信息 ===")
	for shardID, nodes := range n.Shards {
		fmt.Printf("Shard %d: %d nodes\n", shardID, len(nodes))
		for _, node := range nodes {
			fmt.Printf("  Node %d\n", node.ID)
		}
	}
	fmt.Println("==================")
}