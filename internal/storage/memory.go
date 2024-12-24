package storage

// MemoryStore provides an in-memory implementation of Store
type MemoryStore struct {
	tasks []Task
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks: make([]Task, 0),
	}
}

func (s *MemoryStore) Add(task Task) error {
	s.tasks = append(s.tasks, task)
	return nil
}

func (s *MemoryStore) List() ([]Task, error) {
	return s.tasks, nil
}

func (s *MemoryStore) Update(task Task) error {
	for i, t := range s.tasks {
		if t.ID == task.ID {
			s.tasks[i] = task
			return nil
		}
	}
	return ErrTaskNotFound
}

func (s *MemoryStore) Delete(id string) error {
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}
	return ErrTaskNotFound
}
