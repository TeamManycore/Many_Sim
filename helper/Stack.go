package helper

import (
	"sync"
)

type Stack struct {
	lock     sync.Mutex
	Elements []uint16
	Cap      uint16
	Depth    uint16
}

func NewStack(depth uint16) *Stack {
	return &Stack{sync.Mutex{}, make([]uint16, 0, depth), 0, depth}
}

func (s *Stack) Push(elem uint16) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.Cap < s.Depth {
		s.Elements = append([]uint16{elem}, s.Elements...)
		s.Cap++
	}

}

func (s *Stack) Pop() uint16 {
	s.lock.Lock()
	defer s.lock.Unlock()

	res := s.Elements[0]
	s.Elements = s.Elements[1:]
	s.Cap--
	return res

}

func (s *Stack) Tos() uint16 {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.Elements[0]
}

func (s *Stack) Nos() uint16 {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.Elements[1]
}
