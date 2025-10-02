package twitter

import "container/heap"

type tweet struct {
	id   int
	time int
}

type Twitter struct {
	timestamp int
	tweets    map[int][]tweet
	following map[int]map[int]bool
}

func Constructor() Twitter {
	return Twitter{
		tweets:    make(map[int][]tweet),
		following: make(map[int]map[int]bool),
	}
}

func (t *Twitter) PostTweet(userId int, tweetId int) {
	t.timestamp++
	t.tweets[userId] = append(t.tweets[userId], tweet{tweetId, t.timestamp})
}

func (t *Twitter) GetNewsFeed(userId int) []int {
	h := &maxHeap{}
	heap.Init(h)

	if tweets := t.tweets[userId]; len(tweets) > 0 {
		heap.Push(h, &heapNode{
			tweet: tweets[len(tweets)-1],
			user:  userId,
			idx:   len(tweets) - 1,
		})
	}

	for followeeId := range t.following[userId] {
		if tweets := t.tweets[followeeId]; len(tweets) > 0 {
			heap.Push(h, &heapNode{
				tweet: tweets[len(tweets)-1],
				user:  followeeId,
				idx:   len(tweets) - 1,
			})
		}
	}

	feed := make([]int, 0, 10)
	for h.Len() > 0 && len(feed) < 10 {
		node := heap.Pop(h).(*heapNode)
		feed = append(feed, node.tweet.id)

		if node.idx > 0 {
			heap.Push(h, &heapNode{
				tweet: t.tweets[node.user][node.idx-1],
				user:  node.user,
				idx:   node.idx - 1,
			})
		}
	}

	return feed
}

func (t *Twitter) Follow(followerId int, followeeId int) {
	if followerId == followeeId {
		return
	}
	if t.following[followerId] == nil {
		t.following[followerId] = make(map[int]bool)
	}
	t.following[followerId][followeeId] = true
}

func (t *Twitter) Unfollow(followerId int, followeeId int) {
	delete(t.following[followerId], followeeId)
}

type heapNode struct {
	tweet tweet
	user  int
	idx   int
}

type maxHeap []*heapNode

func (h maxHeap) Len() int {
	return len(h)
}
func (h maxHeap) Less(i, j int) bool {
	return h[i].tweet.time > h[j].tweet.time
}
func (h maxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *maxHeap) Push(x interface{}) {
	*h = append(*h, x.(*heapNode))
}

func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
