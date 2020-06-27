package data

import (
	"github.com/shamhub/threadapp/data/shuffle"
)

type Object struct {
	id   string
	seq  uint64
	data []byte
}

const windowSize uint = 10 // size of random sequence

type RandomSequence struct {
	intSequence *shuffle.ShuffleableIntSlice
}

func NewRandomSequence() *RandomSequence {
	return &RandomSequence{
		intSequence: shuffle.NewIntSequence(windowSize),
	}
}

func CreateObject(id string, dataByte []byte, randomSequence *RandomSequence) *Object {

	return &Object{
		id:   id,
		seq:  randomSequence.intSequence.GetRandomNumber(),
		data: dataByte,
	}
}

func (obj *Object) GetSequenceNumber() uint64 {
	return obj.seq
}
