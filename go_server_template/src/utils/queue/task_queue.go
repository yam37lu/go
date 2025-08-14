package queue

import "sync"

type Task interface {
	Type() string
	StartTime() int64
}

type TaskQueue interface {
	Push(Task)
	Pop(string) Task
}

type taskNode struct {
	Data Task
	Next *taskNode
}

type taskQueue struct {
	headNode *taskNode
	sync.Mutex
}

type tasksQueue struct {
	tasks map[string]*taskQueue
	sync.RWMutex
}

func (queue *tasksQueue) Push(task Task) {
	var queueItem *taskQueue = nil
	{
		queue.RLock()
		defer queue.RUnlock()
		if item, ok := queue.tasks[task.Type()]; ok {
			queueItem = item
		} else {
			queueItem = &taskQueue{headNode: &taskNode{}}
			queue.tasks[task.Type()] = queueItem
		}
	}
	queueItem.Lock()
	defer queueItem.Unlock()
	node := queueItem.headNode
	for {
		if node.Next == nil {
			node.Next = &taskNode{
				task,
				nil,
			}
			break
		} else {
			if node.Next.Data.StartTime() < task.StartTime() {
				node = node.Next
				continue
			}
			node.Next = &taskNode{
				task,
				node.Next,
			}
			break
		}
	}
}

func (queue *tasksQueue) Pop(taskType string) Task {
	var queueItem *taskQueue = nil
	{
		queue.RLock()
		defer queue.RUnlock()
		item, ok := queue.tasks[taskType]
		if !ok {
			return nil
		}
		queueItem = item
	}

	queueItem.Lock()
	defer queueItem.Unlock()

	if queueItem.headNode.Next == nil {
		return nil
	}
	node := queueItem.headNode.Next
	queueItem.headNode.Next = node.Next
	return node.Data
}

func NewTaskQueue() TaskQueue {
	return &tasksQueue{tasks: make(map[string]*taskQueue)}
}
