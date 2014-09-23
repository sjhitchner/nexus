package server

import (
	"errors"
	"net"
	"net/http"
	"time"
)

var StoppedError = errors.New("Listener Stopped")

type StoppableListener struct {
	*net.TCPListener
	stop chan int
}

func NewListener(addr string) (*StoppableListener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	tcpListener, ok := listener.(*net.TCPListener)

	if !ok {
		return nil, errors.New("Cannot wrap listener")
	}

	retval := &StoppableListener{}
	retval.TCPListener = tcpListener
	retval.stop = make(chan int)

	return retval, nil
}

func (t *StoppableListener) Accept() (net.Conn, error) {
	for {
		t.SetDeadline(time.Now().Add(time.Second))

		newConn, err := t.TCPListener.Accept()

		select {
		case <-t.stop:
			return nil, StoppedError
		default:
		}

		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}
		return newConn, err
	}
}

func (t *StoppableListener) Stop() {
	close(t.stop)
}

func Serve(stoppableListener *StoppableListener) error {
	return http.Serve(
		//stoppableListener.(*net.TCPListener),
		stoppableListener,
		nil,
	)
}
