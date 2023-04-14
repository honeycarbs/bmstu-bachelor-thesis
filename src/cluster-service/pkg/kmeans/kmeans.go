package kmeans

import (
	"errors"
	"math/rand"
	"time"
)

type Node []float64

func KMeans(Nodes []Node, clusterCount int, maxRounds int) ([]Node, error) {
	if len(Nodes) < clusterCount {
		return nil, errors.New("amount of nodes is smaller than cluster count")
	}

	// Check to make sure everything is consistent, dimension-wise
	stdLen := 0
	for i, Node := range Nodes {
		curLen := len(Node)

		if i > 0 && len(Node) != stdLen {
			return nil, errors.New("data is not consistent dimension-wise")
		}

		stdLen = curLen
	}

	centroids := make([]Node, clusterCount)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Pick centroid starting points from Nodes
	for i := 0; i < clusterCount; i++ {
		srcIndex := r.Intn(len(Nodes))
		srcLen := len(Nodes[srcIndex])
		centroids[i] = make(Node, srcLen)
		copy(centroids[i], Nodes[r.Intn(len(Nodes))])
	}

	return initialCentroids(Nodes, maxRounds, centroids), nil
}

func initialCentroids(Nodes []Node, maxRounds int, centroids []Node) []Node {
	movement := true
	for i := 0; i < maxRounds && movement; i++ {
		movement = false

		groups := make(map[int][]Node)

		for _, Node := range Nodes {
			near := Nearest(Node, centroids)
			groups[near] = append(groups[near], Node)
		}

		for key, group := range groups {
			newNode := meanNode(group)

			if !equal(centroids[key], newNode) {
				centroids[key] = newNode
				movement = true
			}
		}
	}

	return centroids
}

func equal(n1, n2 Node) bool {
	if len(n1) != len(n2) {
		return false
	}

	for i, v := range n1 {
		if v != n2[i] {
			return false
		}
	}

	return true
}

func Nearest(in Node, nodes []Node) int {
	count := len(nodes)

	results := make(Node, count)
	cnt := make(chan int)
	for i, node := range nodes {
		go func(i int, node, cl Node) {
			results[i] = distance(in, node)
			cnt <- 1
		}(i, node, in)
	}

	wait(cnt, results)

	minI := 0
	curDist := results[0]

	for i, dist := range results {
		if dist < curDist {
			curDist = dist
			minI = i
		}
	}

	return minI
}

func distance(node1 Node, node2 Node) float64 {
	length := len(node1)
	squares := make(Node, length, length)

	cnt := make(chan int)

	for i, _ := range node1 {
		go func(i int) {
			diff := node1[i] - node2[i]
			squares[i] = diff * diff
			cnt <- 1
		}(i)
	}

	wait(cnt, squares)

	sum := 0.0
	for _, val := range squares {
		sum += val
	}

	return sum
}

func meanNode(values []Node) Node {
	newNode := make(Node, len(values[0]))

	for _, value := range values {
		for j := 0; j < len(newNode); j++ {
			newNode[j] += value[j]
		}
	}

	for i, value := range newNode {
		newNode[i] = value / float64(len(values))
	}

	return newNode
}

func wait(c chan int, values Node) {
	count := len(values)

	<-c
	for respCnt := 1; respCnt < count; respCnt++ {
		<-c
	}
}
