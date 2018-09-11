package erlang

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewPort(t *testing.T) {
	p := NewPort()

	if p == nil {
		t.Error("NewPort() => nil, wanted Port")
	}
	if p.r == nil {
		t.Error("Port reader is nil")
	}
	if p.w == nil {
		t.Error("Port writer is nil")
	}
}

func TestReadMsg(t *testing.T) {
	tests := []struct {
		name    string
		msg     []byte
		body    []byte
		wantErr bool
	}{
		{
			"empty body",
			append([]byte{0, 0}, []byte("")...),
			[]byte{},
			false,
		},
		{
			"simple body",
			append([]byte{0, 3}, []byte("abc")...),
			[]byte("abc"),
			false,
		},
		{
			"size too large",
			append([]byte{0, 5}, []byte("abc")...),
			append([]byte("abc"), []byte{0, 0}...),
			false,
		},
		{
			"size too small",
			append([]byte{0, 1}, []byte("a")...),
			[]byte("a"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			p := Port{r: buf, w: buf}

			buf.Write(tt.msg)

			body, err := p.ReadMsg()

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error to be non-nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected error to be nil %s", err.Error())
				}
			}
			if !reflect.DeepEqual(body, tt.body) {
				t.Errorf("p.ReadMsg() => %#v, want %#v", body, tt.body)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name string
		msg  []byte
	}{
		{"simple test", []byte("abc")},
		{"simple test", []byte("abcdef")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			p := Port{r: buf, w: buf}
			wantLength := len(tt.msg) + 2
			length, err := p.Write(tt.msg)

			if err != nil {
				t.Errorf("Expected error to be nil %s", err.Error())
			}
			if length != len(tt.msg)+2 {
				t.Errorf("Expected number of bytes to be %d, got %d", wantLength, length)
			}
			got := make([]byte, wantLength)
			want := append(encodeSize(len(tt.msg)), tt.msg...)

			buf.Read(got)

			if !reflect.DeepEqual(got, want) {
				t.Errorf("Wrote %#v, want %#v", got, want)
			}
		})
	}
}

func TestDecodeSize(t *testing.T) {
	tests := []struct {
		name   string
		header []byte
		size   int
	}{
		{"empty header", []byte{}, 0},
		{"small header", []byte{1}, 0},
		{"two-byte header", []byte{1, 2}, 258},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := decodeSize(tt.header)

			if size != tt.size {
				t.Errorf("decodeSize(%#v) => %#v, want %#v", tt.header, size, tt.size)
			}
		})
	}
}

func TestEncodeSize(t *testing.T) {
	tests := []struct {
		name   string
		length int
		header []byte
	}{
		{"empty", 0, []byte{0, 0}},
		{"1 byte", 1, []byte{0, 1}},
		{"258 bytes", 258, []byte{1, 2}},
	}
	for _, tt := range tests {
		header := encodeSize(tt.length)

		if !reflect.DeepEqual(header, tt.header) {
			t.Errorf("encodeSize(%#v) => %#v, want %#v", tt.length, header, tt.header)
		}
	}
}
