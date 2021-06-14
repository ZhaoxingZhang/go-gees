package nio


import (
"fmt"
"net"
)

// 定义实体
type MasterReactor struct {
	net.Listener
}

type SlaveReactor struct {
	h func(net.Conn)
	buf chan struct{}
}

func (s *SlaveReactor)Handle(conn net.Conn) {
	s.buf <- struct{}{}
	fmt.Println("go func")
	go func(){
		s.h(conn)
		<-s.buf
	}()
}

// 定义 construction method
func NewMaster(l net.Listener) MasterReactor{
	return MasterReactor{l}
}
func NewSlave(fn func(net.Conn), bufSize int) SlaveReactor {
	return SlaveReactor{h: fn, buf: make(chan struct{}, bufSize)}
}

// 定义 handler method, 这个是一个echo method
func HandleConn(conn net.Conn) {
	defer conn.Close()
	packet := make([]byte, 1024)
	fmt.Println("packet1", packet)
	// 如果没有可读数据，也就是读 buffer 为空，则阻塞
	_, _ = conn.Read(packet)
	fmt.Println("Read", packet)
	// 同理，不可写则阻塞
	_, _ = conn.Write(packet)
	fmt.Println("Write", packet)
}

func Service() {
	listen, err := net.Listen("tcp", ":8089")
	if err != nil {
		fmt.Println("listen error: ", err)
		return
	}

	master := NewMaster(listen)
	fmt.Println("NewMaster")
	slave := NewSlave(HandleConn, 1024)

	for {
		conn, err := master.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			break
		}

		// start a new goroutine to handle the new connection
		fmt.Println("Handle")
		slave.Handle(conn)
	}
}
