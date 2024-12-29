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
	enable_print = 0
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

	// Create visited map to mark areas
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	var checkArea func(row, col int) bool
	checkArea = func(row, col int) bool {
		// Check walls
		if row < 0 || row >= rows || col < 0 || col >= cols ||
			plan[row][col] == '#' || visited[row][col] {
			return false
		}

		visited[row][col] = true
		hasDirty := plan[row][col] == '*'

		// Check all 4 directions and combine results
		hasDirty = checkArea(row-1, col) || hasDirty
		hasDirty = checkArea(row+1, col) || hasDirty
		hasDirty = checkArea(row, col-1) || hasDirty
		hasDirty = checkArea(row, col+1) || hasDirty

		return hasDirty
	}

	// Count areas that have dirty cells
	dirtyAreas := 0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if !visited[i][j] && plan[i][j] != '#' {
				// If this new area contains any dirty cell
				if checkArea(i, j) {
					dirtyAreas++
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
	Node  TreeNode
	Level int
}

func (i *Interview) GetMaxDepth(root *TreeNode) int {

	if root == nil {
		return 0
	}

	result := 1
	var queue []nodeWithLevel
	queue = append(queue, nodeWithLevel{Node: *root, Level: 1})

	for len(queue) > 0 {
		curNodeInfo := queue[0]
		if result < curNodeInfo.Level {
			result = curNodeInfo.Level
		}
		queue = queue[1:]
		if curNodeInfo.Node.Left != nil {
			queue = append(queue, nodeWithLevel{Node: *curNodeInfo.Node.Left, Level: curNodeInfo.Level + 1})
		}
		if curNodeInfo.Node.Right != nil {
			queue = append(queue, nodeWithLevel{Node: *curNodeInfo.Node.Right, Level: curNodeInfo.Level + 1})
		}
	}

	return result
}
