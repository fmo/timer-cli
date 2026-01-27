// Package services contains Store
package services

import "time"

type Storer interface {
	Save(*Task) error
	Update(*Task) error
	LoadData() ([][]string, error)
	ResetData() error
}

type Store struct {
	p Persister
}

func NewStore(persister Persister) *Store {
	return &Store{p: persister}
}

func (s *Store) Save(task *Task) error {
	s.p.CreateHeader()
	startTime := task.StartTime.Format(time.RFC3339)
	var endTime string
	if task.EndTime.IsZero() {
		endTime = ""
	} else {
		endTime = task.EndTime.Format(time.RFC3339)
	}

	row := []string{startTime, endTime, string(task.Status)}

	return s.p.Save(row)
}

func (s *Store) Update(task *Task) error {
	row := []string{
		task.StartTime.Format(time.RFC3339),
		task.EndTime.Format(time.RFC3339),
		string(task.Status),
	}
	return s.p.Update(row)
}

func (s *Store) LoadData() ([][]string, error) {
	if err := s.p.CreateHeader(); err != nil {
		return nil, err
	}
	return s.p.LoadData()
}

func (s *Store) ResetData() error {
	if err := s.p.ResetData(); err != nil {
		return err
	}
	if err := s.p.CreateHeader(); err != nil {
		return err
	}
	return nil
}
