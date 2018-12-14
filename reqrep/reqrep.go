package reqrep

import (
	"errors"
	"io/ioutil"
	"log"
	"math"
	"net"
	"time"
)

// Network defines the data necessary to act on a network of nodes
type Network struct {
	SendPorts []string
	LocalAddr string
}

// StartNetwork takes initial network fails and returns a new network
func StartNetwork(ports []string, local string) *Network {
	return &Network{
		SendPorts: ports,
		LocalAddr: local,
	}
}

// Len returns the # of nodes in the network
func (n *Network) Len() int {
	result := len(n.SendPorts) + 1
	return result
}

// Quorum returns the majority of nodes in the network
func (n *Network) Quorum() int {
	result := int(math.Ceil(float64(n.Len()) / 2))
	return result
}

/*
Send contacts all nodes on the network, compiles responses and number of failures. Failures being
defined as number of replies not received
*/
func (n *Network) Send(msg []byte) ([][]byte, int) {
	var failures int
	var responses [][]byte
	for i := 0; i < n.Len()-1; i++ {
		reply := send(msg, n.SendPorts[i])
		response, err := recv(reply)
		if err != nil {
			failures++
		} else {
			responses = append(responses, response.([]byte))
		}
	}
	return responses, failures
}

// Send data to addr, attempt to receive response
func send(msg []byte, addr string) <-chan []byte {
	reply := make(chan []byte, 1)
	// Send message to other nodes
	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return reply
	}
	defer c.Close()

	c.Write(msg)

	// Receive response
	data, err := ioutil.ReadAll(c)

	if err != nil {
		log.Println(err)
		return reply
	}

	// log.Println(data)

	reply <- data
	close(reply)
	return reply
}

const recvTimeout = 3 * time.Second

// Receive data from channel or timeout and return error
func recv(reply <-chan []byte) (interface{}, error) {
	select {
	case response := <-reply:
		return response, nil
	case <-time.After(recvTimeout):
		err := errors.New("recv timeout")
		return nil, err
	}
}
