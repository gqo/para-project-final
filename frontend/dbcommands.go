// Graeme Ferguson | ggf221 | 10/24/18
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"time"

	"../proto"
	"github.com/golang/protobuf/proto"
)

// Review struct holds entry data
type Review struct {
	ID     int32
	Album  string
	Artist string
	Rating int32
	Body   string
}

// reportErr prints the error passed if it's not nil
func reportErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

// req sends a request without expecting a response
func req(data []byte) {
	// Send request
	c, errD := net.Dial("tcp", randAddr())
	reportErr(errD) // Check dial error
	c.Write(data)
}

// reqRep sends a request and expects a response that it then returns
func reqRep(data []byte) *pb.PMessage {
	// Send request
	c, errD := net.Dial("tcp", randAddr())
	reportErr(errD) // Check dial error
	c.Write(data)

	// Receive response
	// buf := make([]byte, 8096)
	// size, errR := c.Read(buf)
	// var buf bytes.Buffer
	// io.Copy(&buf, c)
	buf, errR := ioutil.ReadAll(c)
	reportErr(errR) // Check socket read error
	response := new(pb.PMessage)
	errR = proto.Unmarshal(buf, response)
	reportErr(errR) // Check byte conv error
	return response
}

// createRBody takes review data and constructs protobuf object that it returns
func createRBody(id int32, album string, artist string, rating int32, body string) *pb.PMessage_RBody {
	r := &pb.PMessage_RBody{
		RID:    *proto.Int32(id),
		Album:  *proto.String(album),
		Artist: *proto.String(artist),
		Rating: *proto.Int32(rating),
		Body:   *proto.String(body),
	}
	return r
}

// bCreate takes review data and sends a create request with that data to the backend
func bCreate(album string, artist string, rating int32, body string) {
	// Build PMessage
	// --- Build PMessage main
	t := new(pb.PMessage)
	t.ClientID = *proto.Int32(-1)
	t.MsgType = *proto.Int32(0)
	// --- Build review body data
	r := createRBody(0, album, artist, rating, body)
	t.Reviews = append(t.Reviews, r)
	// --- Convert to bytes
	data, errW := proto.Marshal(t)
	reportErr(errW) // Check write error

	// Handle PMessage
	req(data)
}

// bRead takes a review id and sends a read request with that data to the backend
// and, upon receival, returns the review data and whether it exists on the backend
func bRead(id int32) (*Review, bool, error) {
	// Build PMessage
	// --- Build PMessage main
	t := new(pb.PMessage)
	t.ClientID = *proto.Int32(-1)
	t.MsgType = *proto.Int32(1)
	// --- Build review body data
	r := new(pb.PMessage_RBody)
	r.RID = *proto.Int32(id)
	t.Reviews = append(t.Reviews, r)
	// --- Convert to bytes
	data, errW := proto.Marshal(t)
	reportErr(errW) // Check write error

	// Handle request/response structure
	responseInteface, err := reqRepTimeout(data)
	if err != nil {
		log.Println(`dbcommands: bRead(): Read failed`)
		log.Println(err)
	}
	response := responseInteface.(*pb.PMessage)

	// Handle response
	var review *Review
	exist := false
	if response.GetMsgType() == 1 {
		exist = true
		rData := response.GetReviews()[0]
		review = &Review{
			ID:     rData.GetRID(),
			Album:  rData.GetAlbum(),
			Artist: rData.GetArtist(),
			Rating: rData.GetRating(),
			Body:   rData.GetBody(),
		}
	}
	return review, exist, nil
}

// bUpdate takes review data and sends an update request with that data to the backend
func bUpdate(id int32, album string, artist string, rating int32, body string) {
	// Build PMessage
	// --- Build PMessage main
	t := new(pb.PMessage)
	t.ClientID = *proto.Int32(-1)
	t.MsgType = *proto.Int32(2)
	// --- Build review body data
	r := createRBody(id, album, artist, rating, body)
	t.Reviews = append(t.Reviews, r)
	// --- Conv to bytes
	data, errW := proto.Marshal(t)
	reportErr(errW) // Check write error

	// Handle PMessage
	req(data)
}

// bDelete takes a review id and sends a delete request with that data to the backend
func bDelete(id int32) {
	// Build PMessage
	// --- Build PMessage main
	t := new(pb.PMessage)
	t.ClientID = *proto.Int32(-1)
	t.MsgType = *proto.Int32(3)
	// --- Build review body data
	r := new(pb.PMessage_RBody)
	r.RID = *proto.Int32(id)
	t.Reviews = append(t.Reviews, r)
	// --- Convert to bytes
	data, errW := proto.Marshal(t)
	reportErr(errW) // Check write error

	// Handle PMessage
	req(data)
}

// bReadAll sends a read all request to the backend and, upon receival, converts
// protobuf data to a map of reviews and returns said map
func bReadAll() map[int32]*Review {
	// Build PMessage
	// --- Build PMessage main
	t := new(pb.PMessage)
	t.ClientID = *proto.Int32(-1)
	t.MsgType = *proto.Int32(4)
	// --- Convert to bytes
	data, errW := proto.Marshal(t)
	reportErr(errW) // Check write error

	// Handle request/response structure
	responseInteface, err := reqRepTimeout(data)
	if err != nil {
		log.Println(`dbcommands: bRead(): Read failed`)
		log.Println(err)
	}
	response := responseInteface.(*pb.PMessage)

	// Handle response
	rData := response.GetReviews()
	var rMap = map[int32]*Review{}
	for _, item := range rData {
		rMap[item.GetRID()] = &Review{
			ID:     item.GetRID(),
			Album:  item.GetAlbum(),
			Artist: item.GetArtist(),
			Rating: item.GetRating(),
			Body:   item.GetBody(),
		}
	}
	return rMap
}

func bPing() {
	// Build PMessage
	// --- Build PMessage main
	t := new(pb.PMessage)
	t.ClientID = *proto.Int32(-1)
	t.MsgType = *proto.Int32(5)
	// --- Convert to bytes
	data, errW := proto.Marshal(t)
	reportErr(errW)

	duration := time.Second * 5
	buf := make([]byte, 128)
	var waited bool
	for {
		waited = false
		c, errD := net.Dial("tcp", randAddr())
		if errD != nil {
			fmt.Println("Detected failure on", randAddr(), "at", time.Now().UTC())
		} else {
			_, errWD := c.Write(data)
			if errWD != nil {
				log.Println("Detected failure on", randAddr(), "at", time.Now().UTC())
			} else {
				c.SetReadDeadline(time.Now().Add(duration))
				_, errR := c.Read(buf)
				if errR != nil {
					log.Println("Detected failure on", randAddr(), "at", time.Now().UTC())
					waited = true
				}
			}
			c.Close()
		}
		if !waited {
			time.Sleep(duration)
		}
		_ = bStart()
	}
}

// reqRepTimeout takes data, sends it to a port, waits for a response, and returns
// either the response or an error on timeout (or failure)
func reqRepTimeout(data []byte) (interface{}, error) {
dialAnother:
	randAddr := randAddr()
	var failDial int
	c, err := net.Dial("tcp", randAddr)
	if err != nil {
		log.Println(`Dial`, randAddr, `failed`)
		failDial++
		if failDial == len(backendPorts) {
			return nil, errors.New("dial fail")
		}
		goto dialAnother
	}
	defer c.Close()

	_, err = c.Write(data)
	if err != nil {
		log.Println(`Write error`)
		return nil, errors.New("write fail")
	}

	out := make(chan []byte, 1)
	defer close(out)

	go func() {
		buf, err := ioutil.ReadAll(c)
		if err != nil {
			log.Println(`Read error`)
		}
		out <- buf
	}()

	select {
	case result := <-out:
		response := new(pb.PMessage)
		err := proto.Unmarshal(result, response)
		if err != nil {
			log.Println(`Unmarshal response fail`)
		}
		return response, nil
	case <-time.After(3 * time.Second):
		return nil, errors.New(`Read timeout`)
	}
}

// Starts all backends
func bStart() error {
	// Build PMessage
	// --- Build PMessage main
	t := new(pb.PMessage)
	t.ClientID = *proto.Int32(-1)
	t.MsgType = *proto.Int32(-1)
	// --- Convert to bytes
	data, errW := proto.Marshal(t)
	reportErr(errW)
	log.Println(`Attempting to contact:`, backendPorts)
	for i := range backendPorts {
		c, err := net.Dial("tcp", backendPorts[i])
		if err != nil {
			log.Println(`dbcommands: bStart(): Could not connect to backend`)
			return errors.New("backend failed before start")
		}
		c.Write(data)

		// Receive response
		buf, errR := ioutil.ReadAll(c)
		reportErr(errR) // Check socket read error
		response := new(pb.PMessage)
		errR = proto.Unmarshal(buf, response)
		reportErr(errR) // Check byte conv error

		log.Println(`Backend`, response.GetClientID(), `contacted and started`)
	}
	return nil
}

// Returns a random address from list of ports
func randAddr() string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(backendPorts))
	return backendPorts[i]
}
