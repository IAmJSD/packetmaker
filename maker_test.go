package packetmaker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaker_Byte(t *testing.T) {
	assert.Equal(t, "1", string(New().Byte('1').Make()))
}

func TestMaker_String(t *testing.T) {
	assert.Equal(t, "123", string(New().String("123").Make()))
}

func TestMaker_Bytes(t *testing.T) {
	tests := []struct {
		name string

		bytes []byte
	}{
		{
			name: "nil",
		},
		{
			name:  "content",
			bytes: []byte("123"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maker := Maker{}

			// Add an 'a' to stop lazy check.
			maker.Byte('a')

			b := maker.Bytes(tt.bytes).Make()
			if tt.bytes == nil {
				assert.Equal(t, "a", string(b))
			} else {
				assert.Equal(t, string(append([]byte{'a'}, tt.bytes...)), string(b))
			}
		})
	}
}

func TestMaker_UInt16(t *testing.T) {
	tests := []struct {
		name string

		littleEndian bool
		expects      []byte
	}{
		{
			name:         "little endian",
			littleEndian: true,
			expects:      []byte{0x60, 0xea},
		},
		{
			name:    "big endian",
			expects: []byte{0xea, 0x60},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New().Uint16(60000, tt.littleEndian).Make()
			assert.Equal(t, tt.expects, b)
		})
	}
}

func TestMaker_UInt32(t *testing.T) {
	tests := []struct {
		name string

		littleEndian bool
		expects      []byte
	}{
		{
			name:         "little endian",
			littleEndian: true,
			expects:      []byte{0x45, 0x9, 0x61, 0xf4},
		},
		{
			name:    "big endian",
			expects: []byte{0xf4, 0x61, 0x9, 0x45},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New().Uint32(4_100_000_069, tt.littleEndian).Make()
			assert.Equal(t, tt.expects, b)
		})
	}
}

func TestMaker_Uint64(t *testing.T) {
	tests := []struct {
		name string

		littleEndian bool
		expects      []byte
	}{
		{
			name:         "little endian",
			littleEndian: true,
			expects:      []byte{0xff, 0xff, 0x2c, 0xee, 0x5c, 0x7d, 0xcf, 0xf9},
		},
		{
			name:    "big endian",
			expects: []byte{0xf9, 0xcf, 0x7d, 0x5c, 0xee, 0x2c, 0xff, 0xff},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New().Uint64(18_000_744_073_709_551_615, tt.littleEndian).Make()
			assert.Equal(t, tt.expects, b)
		})
	}
}

func TestMaker_Int16(t *testing.T) {
	tests := []struct {
		name string

		littleEndian bool
		expects      []byte
	}{
		{
			name:         "little endian",
			littleEndian: true,
			expects:      []byte{0xbb, 0x82},
		},
		{
			name:    "big endian",
			expects: []byte{0x82, 0xbb},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New().Int16(-32069, tt.littleEndian).Make()
			assert.Equal(t, tt.expects, b)
		})
	}
}

func TestMaker_Int32(t *testing.T) {
	tests := []struct {
		name string

		littleEndian bool
		expects      []byte
	}{
		{
			name:         "little endian",
			littleEndian: true,
			expects:      []byte{0xbb, 0x6b, 0xca, 0x88},
		},
		{
			name:    "big endian",
			expects: []byte{0x88, 0xca, 0x6b, 0xbb},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New().Int32(-2_000_000_069, tt.littleEndian).Make()
			assert.Equal(t, tt.expects, b)
		})
	}
}

func TestMaker_Int64(t *testing.T) {
	tests := []struct {
		name string

		littleEndian bool
		expects      []byte
	}{
		{
			name:         "little endian",
			littleEndian: true,
			expects:      []byte{0xa9, 0x67, 0xdd, 0xf2, 0x30, 0x5c, 0xf1, 0x8c},
		},
		{
			name:    "big endian",
			expects: []byte{0x8c, 0xf1, 0x5c, 0x30, 0xf2, 0xdd, 0x67, 0xa9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New().Int64(-8_290_744_073_709_590_615, tt.littleEndian).Make()
			assert.Equal(t, tt.expects, b)
		})
	}
}

func TestMaker_Make(t *testing.T) {
	m := &Maker{}
	m.add(func(s []byte) {
		s[0] = '1'
	}, 1)
	m.add(func(s []byte) {
		s[0] = '2'
		s[1] = '3'
	}, 2)
	m.add(func(s []byte) {
		s[0] = '4'
		s[1] = '5'
		s[2] = '6'
	}, 3)
	assert.Equal(t, "123456", string(m.Make()))
}

func TestMaker_ManipulatePad(t *testing.T) {
	tests := []struct {
		name string

		pad            []byte
		expectsRealloc bool
	}{
		{
			name:           "nil",
			pad:            nil,
			expectsRealloc: true,
		},
		{
			name:           "length allows",
			pad:            make([]byte, 10),
			expectsRealloc: false,
		},
		{
			name:           "capacity allows",
			pad:            make([]byte, 2, 10),
			expectsRealloc: false,
		},
		{
			name:           "reallocation required",
			pad:            make([]byte, 2, 3),
			expectsRealloc: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Maker{}
			m.add(func(s []byte) {
				s[0] = '1'
			}, 1)
			m.add(func(s []byte) {
				s[0] = '2'
				s[1] = '3'
			}, 2)
			m.add(func(s []byte) {
				s[0] = '4'
				s[1] = '5'
				s[2] = '6'
			}, 3)
			section, realloc := m.ManipulatePad(tt.pad)
			assert.Equal(t, "123456", string(section))
			assert.Equal(t, tt.expectsRealloc, realloc)
		})
	}
}

func TestNew(t *testing.T) {
	assert.Equal(t, &Maker{}, New())
}
