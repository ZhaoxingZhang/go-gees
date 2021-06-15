package geecache

import (
	"fmt"
	"github.com/ZhaoxingZhang/go-gees/geecache/consistenthash"
	pb "github.com/ZhaoxingZhang/go-gees/geecache/geecachepb"
	"github.com/ZhaoxingZhang/go-gees/geecommon/log"
	"github.com/gogo/protobuf/proto"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const defaultBasePath = "/_geecache/"
const defaultReplicas = 3

// HTTPPool
//	implements PeerPicker for a pool of HTTP peers.
// 	implement ServeHTTP interface as HTTPServer

// ServeHTTP() 中使用 proto.Marshal() 编码 HTTP 响应。
// Get() 中使用 proto.Unmarshal() 解码 HTTP 响应

type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	self     string
	basePath string

	mu          sync.Mutex // guards peers and httpGetters
	// 一致性哈希算法的 Map,用来根据具体的 key 选择节点
	peers       *consistenthash.Map
	// 映射远程节点与对应的 httpGetter 1v1
	httpGetters map[string]*httpGetter // keyed by e.g. "http://10.0.0.2:8008"
}

// HTTP server object, implement ServeHTTP interface
// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Info(fmt.Sprintf("[Server %s] %s", p.self, fmt.Sprintf(format, v...)))
}

// ServeHTTP handle all http requests
/*
1. 首先判断访问路径的前缀是否是 basePath，不是返回错误。
2. 约定访问路径格式为 /<basepath>/<groupname>/<key>，
	通过 groupname 得到 group 实例，再使用 group.Get(key) 获取缓存数据。
	最终使用 w.Write() 将缓存值作为 httpResponse 的 body 返回
*/
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		// http: panic serving 127.0.0.1:57604, serve will recover in goroutine
		http.Error(w, "unexpected path: " + r.URL.Path, http.StatusBadRequest)
		log.Error(fmt.Sprintf("Serving %v, unexpected path: %v" , r.Host, r.URL.Path))
		return
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)

}

// Set updates the pool's list of peers.
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// PickPeer picks a peer according to key
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)

//  HTTP 客户端类 httpGetter，实现 PeerGetter 接口
type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	) // e.g. http://localhost:9999/geecache/scores/Tom
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	err = proto.Unmarshal(bytes, out)
	if err != nil {
		return fmt.Errorf("unmarshal res error: %v", err)
	}

	return nil
}

var _ PeerGetter = (*httpGetter)(nil)