package shuffle

import (
	"fmt"
	"math/rand"
	"time"
)

type Shuffleable interface {
	length() uint
	swap(i, j uint)
}

func shuffle(s Shuffleable) {
	rand.Seed(time.Now().UnixNano())
	var i uint
	for i = 0; i < s.length(); i++ {
		j := rand.Intn(int(s.length() - i))
		s.swap(i, uint(j))
	}
}

type ShuffleableIntSlice struct {
	intSeq                []uint64
	windowSize            uint
	nextRandomNumberIndex uint   // ranges from 0 - (windowSize-1)
	lastSeqnumber         uint64 // ranges from 0 - 18446744073709551615 and wraparound
}

// windowSize is the length & capacity of underlying array
func NewIntSequence(windowSize uint) *ShuffleableIntSlice {
	return &ShuffleableIntSlice{
		intSeq:                make([]uint64, windowSize),
		windowSize:            windowSize,
		nextRandomNumberIndex: 0,
		lastSeqnumber:         0,
	}
}

func (s *ShuffleableIntSlice) populateElements() {
	value := s.lastSeqnumber
	// fmt.Printf("last seq number: %d\n", lastSeqnumber)
	lengthOfSlice := s.length()
	var i uint
	for i = 0; i < lengthOfSlice; i++ {
		s.insert(i, value)
		value = value + 1
	}
	s.lastSeqnumber = value

}

func (s *ShuffleableIntSlice) insert(index uint, value uint64) error {
	if s == nil {
		return fmt.Errorf("%v object", s)
	}
	s.intSeq[index] = value
	return nil
}

func (s *ShuffleableIntSlice) length() uint {
	return s.windowSize
}

func (s *ShuffleableIntSlice) swap(i, j uint) {
	temp := s.getElement(j)
	s.insert(j, s.getElement(i))
	s.insert(i, temp)
}

func (s *ShuffleableIntSlice) GetRandomNumber() uint64 {
	if s.nextRandomNumberIndex == 0 {
		s.populateElements()
		shuffle(s)
		//fmt.Printf("random slice: %v\n", randomIntSeq)
	}
	randomNumber := s.getElement(s.nextRandomNumberIndex)
	(s.nextRandomNumberIndex)++
	if s.nextRandomNumberIndex == s.windowSize {
		s.nextRandomNumberIndex = 0
	}
	return randomNumber
}

func (s *ShuffleableIntSlice) getElement(index uint) uint64 {
	return s.intSeq[index]
}
