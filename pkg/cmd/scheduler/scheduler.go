package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Name       string
	Interval   time.Duration
	MaxRetries int
	Body       func() error
}

type Scheduler struct {
	tasks   map[string]*Task
	control chan string
	results chan error
	logger  func(string)
	running bool
	mu      sync.RWMutex
}

func NewScheduler(logger func(string)) *Scheduler {
	return &Scheduler{
		tasks:   make(map[string]*Task),
		control: make(chan string),
		results: make(chan error, 10), // buffered channel to prevent blocking
		logger:  logger,
	}
}

func (s *Scheduler) AddTask(t *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[t.Name] = t
}

func (s *Scheduler) RemoveTask(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks, name)
}

func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		s.logger("scheduler already running")
		return
	}
	s.running = true
	go s.run()
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		s.logger("scheduler not running")
		return
	}
	s.control <- "stop"
	s.running = false
}

func (s *Scheduler) Pause() {
	s.control <- "pause"
}

func (s *Scheduler) Resume() {
	s.control <- "resume"
}

func (s *Scheduler) run() {
	ticker := time.NewTicker(time.Second)
	var paused bool

	for s.running {
		select {
		case <-ticker.C:
			if paused {
				continue
			}
			s.mu.RLock()
			for _, task := range s.tasks {
				go func(t *Task) {
					var err error
					for i := 0; i < t.MaxRetries+1; i++ {
						if err = t.Body(); err == nil {
							break
						}
						s.logger(fmt.Sprintf("%s failed: %v", t.Name, err))
					}
					s.results <- err
				}(task)
			}
			s.mu.RUnlock()
		case cmd := <-s.control:
			switch cmd {
			case "stop":
				s.mu.Lock()
				s.running = false
				s.mu.Unlock()
				return
			case "pause":
				paused = true
			case "resume":
				paused = false
			}
		case err := <-s.results:
			if err != nil {
				s.logger(fmt.Sprintf("Task failed: %v", err))
			}
		}
	}
}
