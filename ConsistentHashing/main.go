package main

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strconv"
)

// HashRing 表示一致性哈希环
type HashRing struct {
	nodes       []int          // 哈希环上的节点
	nodeMap     map[int]string // 哈希值到节点的映射
	virtualNode int            // 每个物理节点对应的虚拟节点数
}

// NewHashRing 创建一个新的哈希环
func NewHashRing(virtualNode int) *HashRing {
	return &HashRing{
		nodeMap:     make(map[int]string),
		virtualNode: virtualNode,
	}
}

// AddNode 添加一个节点到哈希环
func (h *HashRing) AddNode(node string) {
	for i := 0; i < h.virtualNode; i++ {
		// 在环中新增虚拟节点
		virtualNodeKey := node + "#" + strconv.Itoa(i)
		// 虚拟节点hash
		hash := int(h.hash(virtualNodeKey))
		// 增加虚拟节点
		h.nodes = append(h.nodes, hash)
		// 虚拟节点 -> 物理节点映射
		h.nodeMap[hash] = node
	}

	// 对所有虚拟节点排序
	sort.Ints(h.nodes)
}

// RemoveNode 从哈希环中移除一个节点
func (h *HashRing) RemoveNode(node string) {
	// 在虚拟节点中删掉该node
	for i := 0; i < h.virtualNode; i++ {
		virtualNodeKey := node + "#" + strconv.Itoa(i)
		hash := int(h.hash(virtualNodeKey))
		index := sort.SearchInts(h.nodes, hash)
		h.nodes = append(h.nodes[:index], h.nodes[index+1:]...)
		delete(h.nodeMap, hash)
	}
}

// GetNode 获取给定键对应的节点
func (h *HashRing) GetNode(key string) string {
	// 计算该值的hash
	hash := int(h.hash(key))
	// 在虚拟节点中查找最接近该hash的
	index := sort.SearchInts(h.nodes, hash)
	if index >= len(h.nodes) {
		index = 0
	}
	return h.nodeMap[h.nodes[index]]
}

// hash 计算字符串的哈希值
func (h *HashRing) hash(key string) uint32 {
	hash := md5.Sum([]byte(key))
	return uint32(hash[0])<<24 | uint32(hash[1])<<16 | uint32(hash[2])<<8 | uint32(hash[3])
}

func main() {
	// 创建一个哈希环，每个物理节点对应3个虚拟节点
	hashRing := NewHashRing(3)

	// 添加节点
	hashRing.AddNode("Node1")
	hashRing.AddNode("Node2")
	hashRing.AddNode("Node3")

	// 查找键对应的节点
	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	for _, key := range keys {
		node := hashRing.GetNode(key)
		fmt.Printf("Key: %s, Node: %s\n", key, node)
	}

	// 移除一个节点
	hashRing.RemoveNode("Node2")

	// 再次查找键对应的节点
	for _, key := range keys {
		node := hashRing.GetNode(key)
		fmt.Printf("Key: %s, Node: %s\n", key, node)
	}
}
