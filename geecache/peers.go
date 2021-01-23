package geecache

import pb "github.com/ZhaoxingZhang/geecache/geecachepb"
// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.

// 1. PickPeer() 方法用于根据传入的 key 选择相应节点 PeerGetter。
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer.

// 2. 接口 PeerGetter 的 Get() 方法用于从对应 group 查找缓存值
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
