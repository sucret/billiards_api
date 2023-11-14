package cron

import (
	"sync"

	"github.com/robfig/cron/v3"
)

type Server interface {
	Start()
	Add(Task)
	Remove(int32)
	Update(Task)
	Stop()
}

type server struct {
	cronServer *cron.Cron
	mu         sync.Mutex
	task       map[int32][]Task
}

type Task struct {
	TaskId     int32
	TntryId    cron.EntryID
	Spec       string
	ProcessNum int32
	Func       func()
}

var (
	once     sync.Once
	instance *server
)

func New() *server {
	s := new(server)
	s.cronServer = cron.New()
	s.task = make(map[int32][]Task)
	s.cronServer.Start()
	return s
}

// GetCron 单例模式获取
func GetCron() *server {
	once.Do(func() {
		instance = New()
	})

	return instance
}

// Start 启动任务
func (s *server) Start() {
	s.cronServer.Start()
}

// Add 添加任务
func (s *server) Add(task Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task.ProcessNum == 0 {
		return
	}

	for i := int32(0); i < task.ProcessNum; i++ {
		id, err := s.cronServer.AddFunc(task.Spec, task.Func)
		if err == nil {
			task.TntryId = id
			s.task[task.TaskId] = append(s.task[task.TaskId], task)
		}
	}
}

// Remove 移除任务
func (s *server) Remove(taskId int32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	taskList := s.task[taskId]

	for _, task := range taskList {
		s.cronServer.Remove(task.TntryId)
	}

	delete(s.task, taskId)
}

// 修改任务
func (s *server) Update(task Task) {
	// 先移除任务
	s.Remove(task.TaskId)

	// 再添加任务
	s.Add(task)
}

// Stop 停止任务
func (s *server) Stop() {
	s.cronServer.Stop()
}

// TaskMap 获取任务map
func (s *server) TaskMap() map[int32][]Task {
	return s.task
}
