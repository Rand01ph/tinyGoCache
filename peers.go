package tinyGoCache

import pb "github.com/Rand01ph/tinyGoCache/tinygocachepb"

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
