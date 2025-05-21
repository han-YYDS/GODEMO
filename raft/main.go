package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// State 表示节点的状态
type State int

const (
	Follower State = iota
	Candidate
	Leader
)

// Node 表示Raft节点
type Node struct {
	id        int
	state     State
	term      int
	votes     int
	timeout   time.Duration
	election  chan bool
	heartbeat chan bool
	mu        sync.Mutex
	peers     []*Node
}

// NewNode 创建一个新的Raft节点
func NewNode(id int) *Node {
	node := &Node{
		id:        id,
		state:     Follower,
		term:      0,
		votes:     0,
		timeout:   time.Duration(rand.Intn(300)+150) * time.Millisecond,
		election:  make(chan bool),
		heartbeat: make(chan bool),
		peers:     []*Node{},
	}
	go node.run()
	return node
}

// run 节点的主循环
func (n *Node) run() {
	for {
		switch n.state {
		case Follower:
			n.runFollower()
		case Candidate:
			n.runCandidate()
		case Leader:
			n.runLeader()
		}
	}
}

// runFollower 追随者状态的逻辑
func (n *Node) runFollower() {
	select {
	case <-n.heartbeat:
		// 收到心跳，重置选举超时
		n.resetTimeout()
	case <-time.After(n.timeout):
		// 选举超时，转换为候选人
		n.mu.Lock()
		n.state = Candidate
		n.term++
		n.mu.Unlock()
	}
}

// runCandidate 候选人状态的逻辑
func (n *Node) runCandidate() {
	n.mu.Lock()
	n.term++
	n.votes = 1 // 给自己投票
	n.mu.Unlock()

	// 向其他节点发送投票请求
	for _, peer := range n.peers {
		// 模拟其他节点投票
		peer.election <- true
	}

	// 等待投票结果
	for i := 0; i < len(n.peers); i++ {
		if <-n.election {
			n.mu.Lock()
			n.votes++
			n.mu.Unlock()
		}
	}

	if n.votes > len(n.peers)/2 {
		n.mu.Lock()
		n.state = Leader
		n.mu.Unlock()
	} else {
		// 选举失败，重新开始
		n.resetTimeout()
		n.mu.Lock()
		n.state = Follower
		n.mu.Unlock()
	}
}

// runLeader 领导者状态的逻辑
func (n *Node) runLeader() {
	for {
		// 发送心跳
		for _, peer := range n.peers {
			peer.heartbeat <- true
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// resetTimeout 重置选举超时
func (n *Node) resetTimeout() {
	n.timeout = time.Duration(rand.Intn(300)+150) * time.Millisecond
}

// AddPeer 添加一个新的节点到集群
func (n *Node) AddPeer(peer *Node) {
	n.mu.Lock()
	n.peers = append(n.peers, peer)
	n.mu.Unlock()
}

// RemovePeer 从集群中移除一个节点
func (n *Node) RemovePeer(peer *Node) {
	n.mu.Lock()
	for i, p := range n.peers {
		if p.id == peer.id {
			n.peers = append(n.peers[:i], n.peers[i+1:]...)
			break
		}
	}
	n.mu.Unlock()
}

// TriggerElection 触发选举
func (n *Node) TriggerElection() {
	n.mu.Lock()
	n.state = Candidate
	n.term++
	n.mu.Unlock()
}

// GetAllNodes 获取当前集群中的所有节点
func (n *Node) GetAllNodes() []*Node {
	n.mu.Lock()
	defer n.mu.Unlock()
	nodes := []*Node{n}
	for _, peer := range n.peers {
		found := false
		for _, existing := range nodes {
			if existing.id == peer.id {
				found = true
				break
			}
		}
		if !found {
			nodes = append(nodes, peer)
		}
	}
	return nodes
}

func main() {
	// 创建三个节点
	node1 := NewNode(1)
	node2 := NewNode(2)
	node3 := NewNode(3)

	// 配置节点的集群信息
	node1.AddPeer(node2)
	node1.AddPeer(node3)
	node2.AddPeer(node1)
	node2.AddPeer(node3)
	node3.AddPeer(node1)
	node3.AddPeer(node2)

	// 持续打印节点状态
	go func() {
		for {
			allNodes := node1.GetAllNodes()
			for _, node := range allNodes {
				fmt.Printf("Node %d: State = %v, Term = %d\n", node.id, node.state, node.term)
			}
			fmt.Println("-------------------")
			time.Sleep(1 * time.Second)
		}
	}()

	// 模拟一段时间后移除一个节点
	time.Sleep(5 * time.Second)
	fmt.Println("Removing node 3 from the cluster...")
	node1.RemovePeer(node3)
	node2.RemovePeer(node3)
	// 触发选举
	node1.TriggerElection()

	// 模拟一段时间后添加一个新节点
	time.Sleep(5 * time.Second)
	newNode := NewNode(4)
	fmt.Println("Adding node 4 to the cluster...")
	node1.AddPeer(newNode)
	node2.AddPeer(newNode)
	newNode.AddPeer(node1)
	newNode.AddPeer(node2)
	// 触发选举
	node1.TriggerElection()

	select {}
}
