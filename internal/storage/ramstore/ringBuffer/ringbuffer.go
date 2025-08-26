package ringBuffer

import "CryptoService/internal/storage"

type RingBuffer struct {
	data  []storage.CryptoVal
	start int
	size  int
}

func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		data: make([]storage.CryptoVal, capacity),
	}
}

func (r *RingBuffer) Add(val storage.CryptoVal) {
	r.data[(r.start+r.size)%len(r.data)] = val
	if r.size < len(r.data) {
		r.size++
	} else {
		r.start = (r.start + 1) % len(r.data)
	}
}

func (r *RingBuffer) Values() []storage.CryptoVal {
	res := make([]storage.CryptoVal, r.size)
	for i := 0; i < r.size; i++ {
		res[i] = r.data[(r.start+i)%len(r.data)]
	}
	return res
}

func (r *RingBuffer) Last() (storage.CryptoVal, bool) {
	if r.size == 0 {
		return storage.CryptoVal{}, false
	}

	idx := (r.start + r.size - 1) % len(r.data)
	return r.data[idx], true
}
