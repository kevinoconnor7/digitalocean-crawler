package queue

type Node struct {
	Value interface{}
	Next  *Node
}

type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}

func (q *Queue) Push(value interface{}) {
	node := &Node{
		Value: value,
	}

	if q.Tail != nil {
		q.Tail.Next = node
	}

	q.Tail = node

	if q.Head == nil {
		q.Head = node
	}

	q.Length++
}

func (q *Queue) Pop() interface{} {
	node := q.Head

	if node == nil {
		return nil
	}

	q.Head = node.Next
	if q.Head == nil {
		q.Tail = nil
	}

	q.Length--

	return node.Value
}
