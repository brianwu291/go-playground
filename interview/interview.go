package interview

import (
	"fmt"
)

type Interview struct{}

func NewInterview() *Interview {
	return &Interview{}
}

func (i *Interview) FixBug(N int) {
	var enable_print int
	for N > 0 {
		if enable_print == 0 && N%10 != 0 {
			enable_print = 1
		}
		if enable_print == 1 {
			fmt.Print(N % 10)
		}
		N = N / 10
	}
}

func (i *Interview) CountDirtyAreas(plan []string) int {
	if len(plan) == 0 || len(plan[0]) == 0 {
		return 0
	}

	rows := len(plan)
	cols := len(plan[0])

	// create visited map to mark areas
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	var checkArea func(row, col int) bool
	checkArea = func(row, col int) bool {
		// check walls
		if row < 0 || row >= rows || col < 0 || col >= cols ||
			plan[row][col] == '#' || visited[row][col] {
			return false
		}

		visited[row][col] = true
		hasDirty := plan[row][col] == '*'

		// check all 4 directions and combine results
		hasDirty = checkArea(row-1, col) || hasDirty
		hasDirty = checkArea(row+1, col) || hasDirty
		hasDirty = checkArea(row, col-1) || hasDirty
		hasDirty = checkArea(row, col+1) || hasDirty

		return hasDirty
	}

	// count areas that have dirty cells
	dirtyAreas := 0
	for i := 0; i < rows; i += 1 {
		for j := 0; j < cols; j += 1 {
			if !visited[i][j] && plan[i][j] != '#' {
				// if this new area contains any dirty cell
				if checkArea(i, j) {
					dirtyAreas += 1
				}
			}
		}
	}

	return dirtyAreas
}

type TreeNode struct {
	Left  *TreeNode
	Right *TreeNode
	Value int
}

type nodeWithLevel struct {
	Node  *TreeNode
	Level int
}

func (i *Interview) GetMaxDepth(root *TreeNode) int {

	if root == nil {
		return 0
	}

	result := 1
	var queue []nodeWithLevel
	queue = append(queue, nodeWithLevel{Node: root, Level: 1})

	for len(queue) > 0 {
		curNodeInfo := queue[0]
		if result < curNodeInfo.Level {
			result = curNodeInfo.Level
		}
		queue = queue[1:]
		if curNodeInfo.Node.Left != nil {
			queue = append(queue, nodeWithLevel{Node: curNodeInfo.Node.Left, Level: curNodeInfo.Level + 1})
		}
		if curNodeInfo.Node.Right != nil {
			queue = append(queue, nodeWithLevel{Node: curNodeInfo.Node.Right, Level: curNodeInfo.Level + 1})
		}
	}

	return result
}

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

type ListNode struct {
	Val  int
	Next *ListNode
}

func (i *Interview) RemoveNthFromEnd(head *ListNode, n int) *ListNode {
	var listInSlice []*ListNode
	curNode := head
	for curNode != nil {
		listInSlice = append(listInSlice, curNode)
		curNode = curNode.Next
	}
	posFromStart := len(listInSlice) - n
	if posFromStart == 0 {
		newHead := head.Next
		head.Next = nil
		return newHead
	}
	removedNode := listInSlice[posFromStart]
	newNextNode := removedNode.Next
	preNode := listInSlice[posFromStart-1]
	preNode.Next = newNextNode
	removedNode.Next = nil
	return head
}

// [3,9,20,null,null,15,7]

func (i *Interview) LevelOrder(root *TreeNode) [][]int {
	if root == nil {
		empty := [][]int{}
		return empty
	}
	var result [][]int
	var queue []nodeWithLevel
	queue = append(queue, nodeWithLevel{Level: 0, Node: root})
	firstLevelOrder := []int{root.Value}
	result = append(result, firstLevelOrder)
	for len(queue) > 0 {
		cur := queue[0]
		previousLevel := cur.Level
		nextLevel := previousLevel + 1
		left := cur.Node.Left
		right := cur.Node.Right
		if left != nil {
			if nextLevel > len(result)-1 {
				var empty []int
				result = append(result, empty)
			}
			result[nextLevel] = append(result[nextLevel], left.Value)
			queue = append(queue, nodeWithLevel{Level: nextLevel, Node: left})
		}
		if right != nil {
			if nextLevel > len(result)-1 {
				var empty []int
				result = append(result, empty)
			}
			result[nextLevel] = append(result[nextLevel], right.Value)
			queue = append(queue, nodeWithLevel{Level: nextLevel, Node: right})
		}
		queue = queue[1:]
	}
	return result
}
