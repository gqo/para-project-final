// Graeme Ferguson | ggf221 | 10/24/18
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"

	"../proto"

	rr "../reqrep"
)

// Review struct holds entry data
type Review struct {
	ID     int32
	Album  string
	Artist string
	Rating int32
	Body   string
}

var rMap = map[int32]*Review{}
var idCount int32
var mutex = &sync.RWMutex{}
var nodeID int32
var started bool
var network *rr.Network

// reportErr prints the error passed if it's not nil
func reportErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

/* rvToRBody (Review to RBody) takes a pointer to a review struct and converts it
into a protocol buffer RBody item for use in responses */
func rvToRBody(data *Review) *pb.PMessage_RBody {
	rBody := &pb.PMessage_RBody{
		RID:    *proto.Int32(data.ID),
		Album:  *proto.String(data.Album),
		Artist: *proto.String(data.Artist),
		Rating: *proto.Int32(data.Rating),
		Body:   *proto.String(data.Body),
	}
	return rBody
}

/* addReview takes review data and adds the review to the map with a unique id
based on a simple int counter */
func addReview(album string, artist string, rating int32, body string) {
	// build review struct
	review := &Review{
		ID:     0,
		Album:  album,
		Artist: artist,
		Rating: rating,
		Body:   body,
	}
	/* 	mutex.Lock()
	   	defer mutex.Unlock() */
	idCount++
	review.ID = idCount
	rMap[idCount] = review
}

// createReview reads sent review data from PMessage and adds it to the map
func createReviewData(data *pb.PMessage) {
	// get review data
	rData := data.GetReviews()[0]
	_, _, err := paxos(toProtoPlusOne(lenLog()), newActionRv(pCreate, toReview(rData)))
	if err != nil {
		log.Println(`backend: createReviewData(): paxos err`)
		log.Println(err)
	}
	// addReview(rData.GetAlbum(), rData.GetArtist(), rData.GetRating(), rData.GetBody())
}

/* readReview reads sent review request data from PMessage, checks if the review
exists within the map, and then responds accordingly with either the packaged
review data in a PMessage or a invalid signal in a PMessage */
func readReviewData(data *pb.PMessage) ([]byte, error) {
	// Build PMessage body
	response := new(pb.PMessage)
	// Get review request data
	rData := data.GetReviews()[0]
	mutex.RLock()
	defer mutex.RUnlock()
	if review, ok := rMap[rData.GetRID()]; ok {
		// Convert review data to protocol buffer form
		body := rvToRBody(review)
		// Build review body data
		response.Reviews = append(response.Reviews, body)
		response.MsgType = *proto.Int32(1) // Sig: review found
	} else {
		response.MsgType = *proto.Int32(2) // Sig: review not found
	}
	return proto.Marshal(response)
}

/* updateReview reads sent review data from PMessage and updates the corresponding
review in the map */
func updateReviewData(data *pb.PMessage) {
	// Get review data from PMessage
	rData := data.GetReviews()[0]
	// Update data based on sent data
	/* 	mutex.Lock()
	   	defer mutex.Unlock() */
	// updateReview(rData.GetRID(), rData.GetAlbum(), rData.GetArtist(), rData.GetRating(),
	// 	rData.GetBody())
	_, _, err := paxos(toProtoPlusOne(lenLog()), newActionRv(pUpdate, toReview(rData)))
	if err != nil {
		log.Println(`backend: updateReviewData(): paxos err`)
		log.Println(err)
	}
}

// Exposed update review point for paxos
func updateReview(id int32, album string, artist string, rating int32, body string) {
	review := rMap[id]
	review.Album = album
	review.Artist = artist
	review.Rating = rating
	review.Body = body
}

/* deleteReview reads sent review data from PMessage and deletes the corresponding
review in the map */
func deleteReviewData(data *pb.PMessage) {
	// Get review data from PMessage
	rData := data.GetReviews()[0]
	// Delete data based on sent ID
	/* 	mutex.Lock()
	   	defer mutex.Unlock() */
	// delete(rMap, rData.GetRID())
	_, _, err := paxos(toProtoPlusOne(lenLog()), newActionRv(pDelete, toReview(rData)))
	if err != nil {
		log.Println(`backend: deleteReviewData(): paxos err`)
		log.Println(err)
	}
}

/* readAllReview responds to a READALL request with a response PMessage containing
all reviews in the map */
func readAllReviewData(data *pb.PMessage) ([]byte, error) {
	// Build PMessage body
	response := new(pb.PMessage)
	var body *pb.PMessage_RBody
	/* Copy all reviews from map, convert them to protobuf format, build review body
	data */
	mutex.RLock()
	defer mutex.RUnlock()
	for _, item := range rMap {
		body = rvToRBody(item)
		response.Reviews = append(response.Reviews, body)
	}
	return proto.Marshal(response)
}

func start(data *pb.PMessage) ([]byte, error) {
	log.Println(`backend`, nodeID, `received start signal.`)
	started = true

	t := new(pb.PMessage)
	t.ClientID = *proto.Int32(nodeID)
	t.MsgType = *proto.Int32(-1) // START
	t.Valid = *proto.Int32(1)    // VALID

	return proto.Marshal(t)
}

// Pings all nodes on network
func pingNetwork() {
	t := new(pb.PMessage)
	t.ClientID = *proto.Int32(nodeID)
	t.MsgType = *proto.Int32(5)

	testMsg, err := proto.Marshal(t)
	if err != nil {
		log.Println(`backend: handleRequest(): fail marshal test message`)
		log.Println(err)
	}

	go func() {
		log.Println(`Testing other backends.`)
		responses, failures := network.Send(testMsg)
		log.Println(`Received`, failures, `failures`)
		for i := range responses {
			current := new(pb.PMessage)
			err = proto.Unmarshal(responses[i], current)
			if err != nil {
				// log.Println(`Read responses error:`, err)
			}
			// log.Println(`Received:`, current)
		}
	}()
}

func ack(data *pb.PMessage) ([]byte, error) {
	response := new(pb.PMessage)
	response.MsgType = *proto.Int32(1)
	return proto.Marshal(response)
}

func invalidResponse() []byte {
	response := new(pb.PMessage)
	response.ClientID = *proto.Int32(nodeID)
	response.Valid = *proto.Int32(2) // INVALID
	responseBytes, err := proto.Marshal(response)
	if err != nil {
		log.Println(`backend: invalidResponse(): failed marshal invalid response`)
	}
	return responseBytes
}

/* handleRequest takes a connection, reads the request, passes the request data to
the correct function, writes back a response if necessary, and then closes the
connection */
func handleRequest(conn net.Conn) {
	defer conn.Close()

	// Read request
	buf := make([]byte, 2048)
	size, err := conn.Read(buf)
	if err != nil {
		log.Println(`backend: handleRequest(): conn read error`)
		log.Println(err)
	}

	data := new(pb.PMessage)
	err = proto.Unmarshal(buf[:size], data)
	if err != nil {
		log.Println(`backend: handleRequest(): conn message unmarshal error`)
		log.Println(err)
	}
	// log.Println(data)

	switch started {
	// If not started, check if the signal is to start. Otherwise, respond invalid
	case false:
		if data.GetMsgType() == -1 {
			response, err := start(data)
			if err != nil {
				log.Println(`backend: handleRequest(): start error`)
			}

			conn.Write(response)
		} else {
			conn.Write(invalidResponse())
		}
	// If started, handle messages normally
	case true:
		// Check PMessage type and respond/handle accordingly
		switch data.GetMsgType() {
		case pCreate: // CREATE
			// paxos(toProtoPlusOne(lenLog()),)
			createReviewData(data)
		case pRead: // READ
			response, err := readReviewData(data)
			reportErr(err)

			conn.Write(response)
		case pUpdate: // UPDATE
			updateReviewData(data)
		case pDelete: // DELETE
			deleteReviewData(data)
		case pReadAll: // READALL
			pingNetwork()

			response, err := readAllReviewData(data)
			if err != nil {
				log.Println(`backend: handleRequest(): readAllReview error`)
				log.Println(err)
			}

			conn.Write(response)
		case pPingAck: // PING-ACK
			log.Println(`Received ping containing:`, data)
			response, err := ack(data)
			if err != nil {
				log.Println(`backend: handleRequest(): ack error`)
				log.Println(err)
			}

			conn.Write(response)
		case proposeVal:
			log.Println(`Receive proposal:`, data)
			mutex.Lock()
			defer mutex.Unlock()
			response := promise(data)
			conn.Write(response)
			log.Println(`Responded to proposal.`)
		case acceptVal:
			log.Println(`Receive accept!:`, data)
			mutex.Lock()
			defer mutex.Unlock()
			response := accepted(data)
			conn.Write(response)
			log.Println(`Responded to accept!`)
		default:
			conn.Write(invalidResponse())
		}
	}
}

func main() {
	listenPort := flag.Int("listen", 8090, `Specify port for node to listen for requests on`)
	sendPortsList := flag.String("backend", "", `Specify other node ports`)
	id := flag.Int("id", 1, `Specify id of node`)
	flag.Parse()

	listenPortStr := ":" + strconv.Itoa(*listenPort)
	sendPorts := strings.Split(*sendPortsList, ",")
	nodeID = int32(*id)

	server, err := net.Listen("tcp", listenPortStr)
	if err != nil {
		log.Fatalln(`backend: main(): Could not start server`)
	}
	defer server.Close()

	// fmt.Println("Listening for connections on port:", listenPortStr[1:])
	fmt.Println(`Backend`, nodeID, `listening for connections on port:`, listenPortStr[1:])
	fmt.Println(`View includes:`, sendPorts)

	network = rr.StartNetwork(sendPorts, listenPortStr)

	addReview("Grace", "Jeff Buckley", 9, "David Bowie's favorite album!")
	addReview("Exmilitary", "Death Grips", 10, "Zach Hill is good drummer!")
	addReview("Q: Are We Not Men?", "DEVO", 8, "Those are some funny hats!")

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println(`backend: main(): Could not accept connection`)
		}

		go handleRequest(conn)
	}
}
