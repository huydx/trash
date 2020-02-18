package trash

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

type Spans []Span

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
