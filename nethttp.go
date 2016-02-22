import "net"

func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	for {
		rw, e := l.Accept() // HL
		if e != nil {
			//snip
		}
		tempDelay = 0
		c, err := srv.newConn(rw)
		if err != nil {
			continue
		}
		c.setState(c.rwc, StateNew) // before Serve can return
		go c.serve()                // HL
	}
}
