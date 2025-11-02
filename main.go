package main

import (
	"ai-blockchain-shard/shard"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// 全局网络状态变量
var globalNetwork *shard.Network

// NetworkState 网络状态结构，用于序列化保存
type NetworkState struct {
	Nodes []NodeState `json:"nodes"`
}

// NodeState 节点状态结构
type NodeState struct {
	ID      uint64 `json:"id"`
	ShardID uint64 `json:"shard_id"`
	Address string `json:"address"`
}

const stateFile = "network_state.json"

func main() {
	// 尝试加载已有网络状态
	loadNetworkState()
	
	fmt.Println("A Blockchain Sharding System Simulator using LingMa")
	fmt.Println("=============================")
	
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [command]")
		fmt.Println("Available commands:")
		fmt.Println("  init                    - Initialize the blockchain network")
		fmt.Println("  start <shards> <nodes>  - Start a network with specified shards and nodes per shard")
		fmt.Println("  generate <node_id> <count> - Generate sample transactions for a node")
		fmt.Println("  mine <node_id>          - Mine a block on a node")
		fmt.Println("  status                  - Show network status")
		fmt.Println("  run <shards> <nodes> <tx_count> - Run a complete simulation")
		return
	}
	
	command := os.Args[1]
	
	switch command {
	case "init":
		fmt.Println("Initializing blockchain network...")
		initializeNetwork()
		saveNetworkState()
		
	case "start":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run main.go start <shards> <nodes>")
			return
		}
		
		shards, err1 := strconv.Atoi(os.Args[2])
		nodes, err2 := strconv.Atoi(os.Args[3])
		
		if err1 != nil || err2 != nil || shards <= 0 || nodes <= 0 {
			fmt.Println("Invalid parameters. Please provide positive numbers for shards and nodes.")
			return
		}
		
		startNetwork(shards, nodes)
		saveNetworkState()
		
	case "generate":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run main.go generate <node_id> <count>")
			return
		}
		
		nodeID, err1 := strconv.Atoi(os.Args[2])
		count, err2 := strconv.Atoi(os.Args[3])
		
		if err1 != nil || err2 != nil || count <= 0 {
			fmt.Println("Invalid parameters. Please provide positive numbers for node_id and count.")
			return
		}
		
		generateTransactions(nodeID, count)
		
	case "mine":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go mine <node_id>")
			return
		}
		
		nodeID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid node ID.")
			return
		}
		
		mineBlock(nodeID)
		
	case "status":
		showStatus()
		
	case "run":
		if len(os.Args) < 5 {
			fmt.Println("Usage: go run main.go run <shards> <nodes> <tx_count>")
			fmt.Println("Example: go run main.go run 2 2 10")
			return
		}
		
		shards, err1 := strconv.Atoi(os.Args[2])
		nodes, err2 := strconv.Atoi(os.Args[3])
		txCount, err3 := strconv.Atoi(os.Args[4])
		
		if err1 != nil || err2 != nil || err3 != nil || shards <= 0 || nodes <= 0 || txCount <= 0 {
			fmt.Println("Invalid parameters. Please provide positive numbers for shards, nodes, and tx_count.")
			return
		}
		
		runSimulation(shards, nodes, txCount)
		
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

// loadNetworkState 从文件加载网络状态
func loadNetworkState() {
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		return // 文件不存在，无需加载
	}
	
	data, err := ioutil.ReadFile(stateFile)
	if err != nil {
		fmt.Printf("Warning: Failed to read network state file: %v\n", err)
		return
	}
	
	var state NetworkState
	if err := json.Unmarshal(data, &state); err != nil {
		fmt.Printf("Warning: Failed to parse network state file: %v\n", err)
		return
	}
	
	// 重建网络
	globalNetwork = shard.NewNetwork()
	for _, nodeState := range state.Nodes {
		node := shard.NewNode(nodeState.ID, nodeState.ShardID, nodeState.Address)
		globalNetwork.AddNode(node)
	}
	
	// 设置节点中的全局网络引用
	shard.SetGlobalNetwork(globalNetwork)
}

// saveNetworkState 保存网络状态到文件
func saveNetworkState() {
	if globalNetwork == nil {
		return
	}
	
	// 设置节点中的全局网络引用
	shard.SetGlobalNetwork(globalNetwork)
	
	var state NetworkState
	for _, nodes := range globalNetwork.Shards {
		for _, node := range nodes {
			nodeState := NodeState{
				ID:      node.ID,
				ShardID: node.ShardID,
				Address: node.Address,
			}
			state.Nodes = append(state.Nodes, nodeState)
		}
	}
	
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		fmt.Printf("Error: Failed to serialize network state: %v\n", err)
		return
	}
	
	if err := ioutil.WriteFile(stateFile, data, 0644); err != nil {
		fmt.Printf("Error: Failed to write network state file: %v\n", err)
		return
	}
}

// initializeNetwork 初始化网络
func initializeNetwork() {
	globalNetwork = shard.NewNetwork()
	fmt.Println("Blockchain network initialized successfully!")
}

// startNetwork 启动网络
func startNetwork(shards, nodes int) {
	fmt.Printf("Starting network with %d shards and %d nodes per shard\n", shards, nodes)
	
	// 初始化网络（如果尚未初始化）
	if globalNetwork == nil {
		initializeNetwork()
	}
	
	// 清空现有节点
	globalNetwork.Shards = make(map[uint64][]*shard.Node)
	
	// 创建分片和节点
	nodeID := 1
	for s := 0; s < shards; s++ {
		for n := 0; n < nodes; n++ {
			address := fmt.Sprintf("192.168.1.%d", nodeID)
			node := shard.NewNode(uint64(nodeID), uint64(s), address)
			globalNetwork.AddNode(node)
			nodeID++
		}
	}
	
	// 设置节点中的全局网络引用
	shard.SetGlobalNetwork(globalNetwork)
	
	fmt.Println("Network started successfully!")
	globalNetwork.PrintNetworkInfo()
}

// generateTransactions 生成交易
func generateTransactions(nodeID, count int) {
	fmt.Printf("Generating %d transactions for node %d\n", count, nodeID)
	
	if globalNetwork == nil {
		fmt.Println("Network not initialized. Please run 'start' command first.")
		return
	}
	
	// 查找指定节点
	var targetNode *shard.Node
	for _, nodes := range globalNetwork.Shards {
		for _, node := range nodes {
			if int(node.ID) == nodeID {
				targetNode = node
				break
			}
		}
		if targetNode != nil {
			break
		}
	}
	
	if targetNode == nil {
		fmt.Printf("Node with ID %d not found\n", nodeID)
		return
	}
	
	// 生成交易
	transactions := targetNode.GenerateSampleTransactions(count)
	targetNode.AddTransactions(transactions)
	
	fmt.Printf("Generated %d transactions for node %d\n", len(transactions), nodeID)
	
	// 显示前几个交易
	for i, tx := range transactions {
		if i < 5 { // 只打印前5个交易
			fmt.Printf("Transaction %d: %s -> %s, Amount: %s, Shard: %d -> %d\n", 
				i, tx.Sender, tx.Recipient, tx.Amount.String(), tx.FromShard, tx.ToShard)
		}
	}
	
	// 广播跨分片交易
	for _, tx := range transactions {
		if tx.IsCrossShard() {
			globalNetwork.BroadcastTransaction(tx, tx.ToShard)
		}
	}
}

// mineBlock 挖掘区块
func mineBlock(nodeID int) {
	fmt.Printf("Mining block on node %d\n", nodeID)
	
	if globalNetwork == nil {
		fmt.Println("Network not initialized. Please run 'start' command first.")
		return
	}
	
	// 查找指定节点
	var targetNode *shard.Node
	for _, nodes := range globalNetwork.Shards {
		for _, node := range nodes {
			if int(node.ID) == nodeID {
				targetNode = node
				break
			}
		}
		if targetNode != nil {
			break
		}
	}
	
	if targetNode == nil {
		fmt.Printf("Node with ID %d not found\n", nodeID)
		return
	}
	
	// 挖矿
	block := targetNode.MineBlock()
	if block != nil {
		fmt.Printf("Block #%d mined successfully on node %d!\n", block.Header.Number, nodeID)
	} else {
		fmt.Printf("Failed to mine block on node %d\n", nodeID)
	}
}

// showStatus 显示状态
func showStatus() {
	fmt.Println("Network status:")
	
	if globalNetwork == nil {
		fmt.Println("Network not initialized.")
		return
	}
	
	globalNetwork.PrintNetworkInfo()
	
	// 显示各节点状态
	fmt.Println("\n=== 节点详细状态 ===")
	for shardID, nodes := range globalNetwork.Shards {
		fmt.Printf("Shard %d:\n", shardID)
		for _, node := range nodes {
			fmt.Printf("  Node %d:\n", node.ID)
			fmt.Printf("    Blockchain height: %d\n", node.Blockchain.CurrentHeight())
			fmt.Printf("    Pending transactions: %d\n", node.Blockchain.TxPool.PendingCount())
			fmt.Printf("    Relay transactions: %d\n", node.Blockchain.TxPool.RelayCount())
		}
	}
}

// runSimulation 运行完整模拟
func runSimulation(shards, nodes, txCount int) {
	fmt.Printf("Running blockchain sharding simulation with %d shards, %d nodes per shard, %d transactions per node\n", 
		shards, nodes, txCount)
	
	// 记录开始时间
	startTime := time.Now()
	
	// 1. 初始化网络
	fmt.Println("\n--- Step 1: Initializing network ---")
	startNetwork(shards, nodes)
	saveNetworkState()
	
	// 2. 为每个节点生成交易
	fmt.Println("\n--- Step 2: Generating transactions ---")
	totalNodes := shards * nodes
	totalTransactions := 0
	crossShardTransactions := 0
	
	for i := 1; i <= totalNodes; i++ {
		fmt.Printf("\nGenerating transactions for node %d...\n", i)
		generateTransactions(i, txCount)
		
		// 统计交易数量
		var node *shard.Node
		for _, nodes := range globalNetwork.Shards {
			for _, n := range nodes {
				if int(n.ID) == i {
					node = n
					break
				}
			}
			if node != nil {
				break
			}
		}
		
		if node != nil {
			totalTransactions += txCount
			// 统计跨分片交易（近似值）
			for _, tx := range node.Blockchain.TxPool.PendingTxs {
				if tx.IsCrossShard() {
					crossShardTransactions++
				}
			}
			for _, tx := range node.Blockchain.TxPool.RelayTxs {
				if tx.IsCrossShard() {
					crossShardTransactions++
				}
			}
		}
	}
	
	// 3. 在每个节点上挖矿
	fmt.Println("\n--- Step 3: Mining blocks ---")
	blocksMined := 0
	for i := 1; i <= totalNodes; i++ {
		fmt.Printf("\nMining on node %d...\n", i)
		mineBlock(i)
		blocksMined++
	}
	
	// 4. 输出最终统计信息
	fmt.Println("\n--- Step 4: Final Results ---")
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	
	showStatus()
	
	fmt.Println("\n=== Simulation Summary ===")
	fmt.Printf("Simulation Duration: %v\n", duration)
	fmt.Printf("Shards: %d\n", shards)
	fmt.Printf("Nodes per shard: %d\n", nodes)
	fmt.Printf("Total nodes: %d\n", totalNodes)
	fmt.Printf("Transactions per node: %d\n", txCount)
	fmt.Printf("Total transactions generated: %d\n", totalTransactions)
	fmt.Printf("Approximate cross-shard transactions: %d\n", crossShardTransactions)
	fmt.Printf("Blocks mined: %d\n", blocksMined)
	fmt.Println("=========================")
}