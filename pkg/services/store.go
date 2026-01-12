// Package services contains Store
package services

type Store struct {
	ts TaskStorer
}

func NewStore(ts TaskStorer) *Store {
	return &Store{ts: ts}
}

func (s *Store) Save(task Task) error {
	return s.ts.Save(task)
}

func (s *Store) Update(task Task) error {
	return s.ts.Update(task)
}
