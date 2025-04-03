package packetmaker

import (
	"encoding/binary"
	"unsafe"
)

type queue struct {
	do   func([]byte) int
	next *queue
}

// Maker is used to define the packet maker.
type Maker struct {
	len                   int
	queueFirst, queueLast *queue
}

// Adds to the queue.
func (m *Maker) add(f func([]byte), l int) *Maker {
	m.len += l
	queueItem := &queue{do: func(offsettedChunk []byte) int {
		f(offsettedChunk)
		return l
	}}
	if m.queueLast == nil {
		m.queueFirst = queueItem
		m.queueLast = queueItem
	} else {
		m.queueLast.next = queueItem
		m.queueLast = queueItem
	}
	return m
}

// Len returns the length of the packet maker.
func (m *Maker) Len() int {
	return m.len
}

// Byte is used to add a single byte to the packet.
func (m *Maker) Byte(b uint8) *Maker {
	return m.add(func(offsettedChunk []byte) {
		offsettedChunk[0] = b
	}, 1)
}

// String is used to add a string to the packet.
func (m *Maker) String(s string) *Maker {
	return m.add(func(offsettedChunk []byte) {
		copy(offsettedChunk, s)
	}, len(s))
}

// Bytes is used to add the specified byte slice to the packet. Note the length of the slice CANNOT
// change after this is set. Doing so will result in undefined behaviour.
func (m *Maker) Bytes(b []byte) *Maker {
	return m.add(func(offsettedChunk []byte) {
		copy(offsettedChunk, b)
	}, len(b))
}

// Uint16 is used to add an 16-bit unsigned integer.
func (m *Maker) Uint16(v uint16, littleEndian bool) *Maker {
	hn := binary.BigEndian.PutUint16
	if littleEndian {
		hn = binary.LittleEndian.PutUint16
	}
	return m.add(func(offsettedChunk []byte) {
		hn(offsettedChunk, v)
	}, 2)
}

// Uint32 is used to add an 32-bit unsigned integer.
func (m *Maker) Uint32(v uint32, littleEndian bool) *Maker {
	hn := binary.BigEndian.PutUint32
	if littleEndian {
		hn = binary.LittleEndian.PutUint32
	}
	return m.add(func(offsettedChunk []byte) {
		hn(offsettedChunk, v)
	}, 4)
}

// Uint64 is used to add an 64-bit unsigned integer.
func (m *Maker) Uint64(v uint64, littleEndian bool) *Maker {
	hn := binary.BigEndian.PutUint64
	if littleEndian {
		hn = binary.LittleEndian.PutUint64
	}
	return m.add(func(offsettedChunk []byte) {
		hn(offsettedChunk, v)
	}, 8)
}

// Int16 is used to add an 16-bit signed integer.
func (m *Maker) Int16(v int16, littleEndian bool) *Maker {
	hn := binary.BigEndian.PutUint16
	if littleEndian {
		hn = binary.LittleEndian.PutUint16
	}
	return m.add(func(offsettedChunk []byte) {
		hn(offsettedChunk, *(*uint16)(unsafe.Pointer(&v)))
	}, 2)
}

// Int32 is used to add an 32-bit signed integer.
func (m *Maker) Int32(v int32, littleEndian bool) *Maker {
	hn := binary.BigEndian.PutUint32
	if littleEndian {
		hn = binary.LittleEndian.PutUint32
	}
	return m.add(func(offsettedChunk []byte) {
		hn(offsettedChunk, *(*uint32)(unsafe.Pointer(&v)))
	}, 4)
}

// Int64 is used to add an 64-bit signed integer.
func (m *Maker) Int64(v int64, littleEndian bool) *Maker {
	hn := binary.BigEndian.PutUint64
	if littleEndian {
		hn = binary.LittleEndian.PutUint64
	}
	return m.add(func(offsettedChunk []byte) {
		hn(offsettedChunk, *(*uint64)(unsafe.Pointer(&v)))
	}, 8)
}

// Make is used to make the packet.
func (m *Maker) Make() []byte {
	if m.len == 0 {
		return []byte{}
	}
	offset := 0
	s := make([]byte, m.len)
	q := m.queueFirst
	for q != nil {
		offset += q.do(s[offset:])
		q = q.next
	}
	return s
}

// ManipulatePad is used to try and manipulate a pre-existing pad. This is useful for saving memory.
// If the size of the data is greater than cap, then Make will be called and a new slice will be allocated.
// The part of the slice that was used will be returned and whether it was re-allocated.
func (m *Maker) ManipulatePad(s []byte) (section []byte, reallocated bool) {
	if s != nil {
		if len(s) >= m.len || cap(s) >= m.len {
			// Best case scenario - we do not need to reallocate.
			s = s[:m.len]
			q := m.queueFirst
			offset := 0
			for q != nil {
				offset += q.do(s[offset:])
				q = q.next
			}
			return s, false
		}
	}

	// We need to re-allocate the pad.
	return m.Make(), true
}

// New is used to create a new packet maker.
func New() *Maker {
	return &Maker{}
}
