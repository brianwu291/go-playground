package interview

import (
	"errors"
	"math"
	"strconv"
)

type Interview struct{}

func NewInterview() *Interview {
	return &Interview{}
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

func AddTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	getNodeNext := func(node *ListNode) *ListNode {
		if node == nil {
			return nil
		}
		return node.Next
	}
	getNodeVal := func(node *ListNode) int {
		if node == nil {
			return 0
		}
		return node.Val
	}

	var result *ListNode
	var point *ListNode
	l1Point := l1
	l2Point := l2
	var carry int
	for l1Point != nil || l2Point != nil || carry != 0 {
		curSum := getNodeVal(l1Point) + getNodeVal(l2Point) + carry
		overTen := curSum >= 10
		if overTen {
			curSum = curSum % 10
			carry = 1
		} else {
			carry = 0
		}
		if result == nil {
			firstOne := ListNode{Val: curSum, Next: nil}
			result = &firstOne
			point = result
		} else {
			newNode := ListNode{Val: curSum, Next: nil}
			point.Next = &newNode
			point = point.Next
		}
		l1Point = getNodeNext(l1Point)
		l2Point = getNodeNext(l2Point)
	}
	return result
}

/*
*
[1,0,1]
[1,1,0]
[1,1,0]
*/
func CountSquares(matrix [][]int) int {
	var result int

	rowCount := len(matrix)
	columnCount := len(matrix[0])
	// handle how many row
	for i := 0; i < rowCount; i += 1 {
		ele := matrix[i][0]
		if ele == 1 {
			result += 1
		}
	}

	// handle how many column
	// start from 1 since above calculated
	for i := 1; i < columnCount; i += 1 {
		ele := matrix[0][i]
		if ele == 1 {
			result += 1
		}
	}

	// handle rest one by one
	for i := 1; i < rowCount; i += 1 {
		for j := 1; j < columnCount; j += 1 {
			curEle := matrix[i][j]
			if curEle == 0 {
				continue
			}
			top := matrix[i-1][j]
			left := matrix[i][j-1]
			topLeft := matrix[i-1][j-1]
			if top >= 1 && left >= 1 && topLeft >= 1 {
				// cur is 1
				// get min val from top, left, and topLeft
				// then minVal + cur to assign to cur, for ex, 1 + cur (1) = 2
				// result += cur
				smaller := math.Min(float64(top), float64(left))
				min := int(math.Min(smaller, float64(topLeft)))
				matrix[i][j] = min + curEle
				result += matrix[i][j]
			} else {
				result += curEle
			}
		}
	}

	return result
}

func CoinChange(coins []int, amount int) int {
	// 1, 5, 10
	// 11
	// 2 (1, 10)

	// dp[0] = 0
	// dp[1] = 1
	// dp[2] = dp[1] + 1 --> 2
	// dp[3] = dp[2] + 1 --> 3
	// dp[4] = dp[3] + 1 --> 4
	// dp[5] = 1
	// ...
	// dp[amount] = dp[amount - 1(value of one coin)] + 1 (count of coins)

	// 1, 5, 100
	// 10
	// [0,1,2,3,4,5,6,7,8,9,10]

	if amount <= 0 {
		return 0
	}
	dp := make([]int, amount+1)
	for ac := 1; ac <= amount; ac++ {
		for _, coin := range coins {
			r := ac - coin
			if r == 0 {
				dp[ac] = 1
				continue
			}
			if r > 0 && dp[r] > 0 {
				origin := dp[ac]
				if origin == 0 || dp[r]+1 < origin {
					dp[ac] = dp[r] + 1
				}
			}
		}
	}
	if dp[amount] == 0 {
		return -1
	}
	return dp[amount]
}

// "123" -> (123, 0)
// "0000000123.45" -> (12345, 2)
// "0.123" -> (123, 3)
// 0.
// .2
// "          "
// .
// "        123.1"
// "1.223323" -> (12, 1, nil)
// "-23" -> (-123, 2, nil)
// "1.-23" ->

func ConvertStrToNum(str string) (int, int, error) {
	var num, point int
	length := len(str)
	first := str[0]
	if length == 0 || first == '.' || str[length-1] == '.' {
		return num, point, errors.New("invalid")
	}

	finalIndex := 0
	if first == '-' {
		finalIndex = 1
	}
	var hadPoint bool
	for i := length - 1; i >= finalIndex; i -= 1 {
		cur := str[i]
		if hadPoint && cur == '.' {
			return num, point, errors.New("invalid")
		}
		if (cur > '9' || cur < '0') && cur != '.' {
			return num, point, errors.New("invalid")
		}
		if str[i] == '.' {
			point = length - 1 - i
			hadPoint = true
		} else if cur != '0' {
			stand := length - i - 1
			if hadPoint {
				stand -= 1
			}
			parseInt, _ := strconv.Atoi(string(cur))
			num += parseInt * int(math.Pow10(stand))
		}
	}

	if finalIndex == 1 {
		num = -num
	}
	return num, point, nil
}

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

type NodeWithLevel struct {
	Node  *TreeNode
	Level int
}

func rightSideView(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}
	res := []int{root.Value}
	queue := []*NodeWithLevel{{Node: root, Level: 0}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		curLevel := cur.Level
		if len(res)-1 <= curLevel {
			for i := len(res); i <= curLevel; i += 1 {
				res = append(res, math.MinInt)
			}
		}
		res[cur.Level] = cur.Node.Value
		if cur.Node.Left != nil {
			queue = append(queue, &NodeWithLevel{Node: cur.Node.Left, Level: cur.Level + 1})
		}
		if cur.Node.Right != nil {
			queue = append(queue, &NodeWithLevel{Node: cur.Node.Right, Level: cur.Level + 1})
		}
	}

	return res
}
