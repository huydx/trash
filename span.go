package trash

import (
	"bytes"
	"encoding/gob"
	"io"
	"sync"
)

/**
[
  {
    "id": "352bff9a74ca9ad2",
    "traceId": "5af7183fb1d4cf5f",
    "parentId": "6b221d5bc9e6496c",
    "name": "get /api",
    "timestamp": 1556604172355737,
    "duration": 1431,
    "kind": "SERVER",
    "localEndpoint": {
      "serviceName": "backend",
      "ipv4": "192.168.99.1",
      "port": 3306
    },
    "remoteEndpoint": {
      "ipv4": "172.19.0.2",
      "port": 58648
    },
    "tags": {
      "http.method": "GET",
      "http.path": "/api"
    }
  }
]
*/

type Spans []*Span

type Span struct {
	ID            string `json:"id"`
	TraceID       string `json:"traceId"`
	ParentID      string `json:"parentId"`
	Name          string `json:"name"`
	Timestamp     int64  `json:"timestamp"`
	Duration      int    `json:"duration"`
	Kind          string `json:"kind"`
	LocalEndpoint struct {
		ServiceName string `json:"serviceName"`
		Ipv4        string `json:"ipv4"`
		Port        int    `json:"port"`
	} `json:"localEndpoint"`
	RemoteEndpoint struct {
		Ipv4 string `json:"ipv4"`
		Port int    `json:"port"`
	} `json:"remoteEndpoint"`
	Tags map[string]string `json:"tags"`
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer([]byte{})
	},
}

func getBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func putBuffer(bs *bytes.Buffer) {
	bufferPool.Put(bs)
}

func (s *Span) Marshal() ([]byte, error) {
	buf := getBuffer()
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(s); err != nil {
		return nil, err
	}
	dup := make([]byte, buf.Len())
	copy(dup, buf.Bytes())
	putBuffer(buf)
	return dup, nil
}

func (s *Span) Unmarshal(bs []byte) error {
	buf := bytes.NewBuffer(bs)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(s)
}

func (s *Spans) Marshal() ([]byte, error) {
	buf := getBuffer()
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(s); err != nil {
		return nil, err
	}
	dup := make([]byte, buf.Len())
	copy(dup, buf.Bytes())
	putBuffer(buf)
	return dup, nil
}

func (s *Spans) Unmarshal(bs []byte) error {
	buf := getBuffer()
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(s); err == io.EOF {
		return nil
	} else {
		return err
	}
}
