package main

import (
	"errors"
	"log"
	"sync"

	"../proto"
	"github.com/golang/protobuf/proto"
)

const (
	pCreate = iota
	pRead
	pUpdate
	pDelete
	pReadAll
	pPingAck
	proposeVal
	promiseVal
	oldPromiseVal // if node proposed val for something that already exists, must learn old log value
	nackPromiseVal
	acceptVal
	acceptedVal
	bullyVal
)

// Action represents an action taken upon the database of reviews
type Action struct {
	Op  int32
	Val Review
}

func newAction(op int32, id int32, album string, artist string, rating int32, body string) Action {
	return Action{
		Op: op,
		Val: Review{
			ID:     id,
			Album:  album,
			Artist: artist,
			Rating: rating,
			Body:   body,
		},
	}
}

// Mutex used for locking dataLog and promiseMap (necessary for paxos)
var paxMutex *sync.Mutex

// Log of all actions performed on dataStore
var dataLog []Action

// Key = ArrVal, Val = HighestPropNum for tha ArrVal in dataLog
var promiseMap map[int64]int64

// Leader ID
var currentLeader = int32(-1)

func lenLog() int64 {
	return int64(len(dataLog))
}

// Protocol Buffers don't transfer 0 values so add one to pass val
func toProtoPlusOne(x int64) int64 {
	x++
	return x
}

func fromProtoMinusOne(x int64) int64 {
	x--
	return x
}

// Message represents messages sent for the purpose of paxos
type Message struct {
	From     string
	ClientID int32
	MsgType  int32
	ArrVal   int64
	PropNum  int64
	Vals     []Action
	LogLen   int64
}

// newMessage returns a new message with given values
func newMessage(from string, clientID int32, msgType int32, arrVal int64, propNum int64, logLen int64, vals ...Action) *Message {
	msg := &Message{
		From:     from,
		ClientID: clientID,
		MsgType:  msgType,
		ArrVal:   arrVal,
		PropNum:  propNum,
		LogLen:   logLen,
	}
	if len(vals) != 0 {
		msg.Vals = vals
	}
	return msg
}

/*
Protobuf message to Golang type conversion functions take protobuf message argument
Golang type to Protobuf message conversion functions are member functions of the type
*/

func toReview(rbody *pb.PMessage_RBody) *Review {
	return &Review{
		ID:     rbody.GetRID(),
		Album:  rbody.GetAlbum(),
		Artist: rbody.GetArtist(),
		Rating: rbody.GetRating(),
		Body:   rbody.GetBody(),
	}
}

func (review *Review) toRBody() *pb.PMessage_RBody {
	return &pb.PMessage_RBody{
		RID:    *proto.Int32(review.ID),
		Album:  *proto.String(review.Album),
		Artist: *proto.String(review.Artist),
		Rating: *proto.Int32(review.Rating),
		Body:   *proto.String(review.Body),
	}
}

func toAction(pAction *pb.PMessage_PAction) *Action {
	return &Action{
		Op:  pAction.GetOp(),
		Val: *toReview(pAction.GetVal()),
	}
}

func (action *Action) toPAction() *pb.PMessage_PAction {
	return &pb.PMessage_PAction{
		Op:  *proto.Int32(action.Op),
		Val: action.Val.toRBody(),
	}
}

func toActions(pActions []*pb.PMessage_PAction) []Action {
	var actions []Action
	for i := range pActions {
		actions = append(actions, *toAction(pActions[i]))
	}
	return actions
}

func toPActions(actions []Action) []*pb.PMessage_PAction {
	var pActions []*pb.PMessage_PAction
	for i := range actions {
		pActions = append(pActions, actions[i].toPAction())
	}
	return pActions
}

func toMessage(pMSG *pb.PMessage) *Message {
	return &Message{
		From:     pMSG.GetFrom(),
		ClientID: pMSG.GetClientID(),
		MsgType:  pMSG.GetMsgType(),
		ArrVal:   pMSG.GetArrVal(),
		PropNum:  pMSG.GetPropNum(),
		Vals:     toActions(pMSG.GetVals()),
		LogLen:   pMSG.GetLogLen(),
	}
}

func (msg *Message) toPMessage() *pb.PMessage {
	return &pb.PMessage{
		From:     msg.From,
		ClientID: msg.ClientID,
		MsgType:  msg.MsgType,
		ArrVal:   msg.ArrVal,
		PropNum:  msg.PropNum,
		Vals:     toPActions(msg.Vals),
		LogLen:   msg.LogLen,
	}
}

// paxos runs a round of the paxos algorithm for a specific section of the dataLog
func paxos(arrVal int64, val Action) (int64, interface{}, error) {
	var propNum int64
	responses, err := propose(arrVal, propNum)
	if err != nil {
		return -1, 0, errors.New("propose received too many failures")
	}

beforeCheck:
	propNum, recvAction, err := checkPromises(propNum, responses.([][]byte))
	if err != nil {
		switch err.Error() {
		case "nack":
			// I'm sorry
			goto beforeCheck
		case "old":
			val = recvAction.(Action)
		case "fail":
			return propNum, nil, err
		}
	}

	promiseMap[arrVal] = propNum

	err = accept(arrVal, propNum, val)
	for err != nil {
		err = accept(arrVal, propNum, val)
	}

	// Update self here with function call

	return propNum, nil, nil
}

// propose sends proposals to other nodes and returns responses or an error if a quorum does not respond
func propose(arrVal int64, propNum int64) (interface{}, error) {
	msg := newMessage(network.LocalAddr, nodeID, proposeVal, arrVal, propNum, toProtoPlusOne(lenLog()))
	protoMsg := msg.toPMessage()
	request, err := proto.Marshal(protoMsg)
	if err != nil {
		log.Println(`paxos: propose(): Request marshal fail`)
		log.Println(err)
	}

	responses, failures := network.Send(request)

	if failures >= network.Quorum() {
		return nil, errors.New("too many failures")
	}

	return responses, nil
}

// checkPromises returns the valid accept message content for an accept given proposal results
func checkPromises(propNum int64, responses [][]byte) (int64, interface{}, error) {
	// if no msg val in response, propNum arrVal pair promised by all nodes
	// if msg val in response, propNum arrVal pair returned that is already promised
	highestPropNum := propNum
	var msgVal Action
	// var setVal bool
	// var setPropNum bool
	var result string
	log.Println(`Checking promises...`)
	for i := range responses {
		var current *pb.PMessage
		log.Println(current)
		err := proto.Unmarshal(responses[i], current)
		if err != nil {
			log.Println(`paxos: propose(): unmarshal response fail`)
		}

		switch current.GetMsgType() {
		case promiseVal:
			// Good to continue with accept phase
			if result == "" {
				result = "promise"
			}
		case nackPromiseVal:
			// Promise val too low, repeat propose with higher number
			if result == "" || result == "nack" {
				result = "nack"
				// If the highest propNum known is lower than currently examined propNum
				if highestPropNum < current.GetPropNum() {
					highestPropNum = current.GetPropNum()
				}
			}
		case oldPromiseVal:
			// Stale log detected, update log, repeat paxos
			if result == "" {
				result = "old"
				highestPropNum = current.GetPropNum()
				msgVal = *toAction(current.GetVals()[0])
			}
		}
	}

	switch result {
	case "promise":
		return highestPropNum, nil, nil
	case "nack":
		return highestPropNum, nil, errors.New("nack")
	case "old":
		return highestPropNum, msgVal, errors.New("old")
	default:
		// This shouldn't happen
		log.Println(`paxos: checkPromises(): checkPromises failed silently`)
		return highestPropNum, nil, errors.New("fail")
	}
	/* lenVals := len(current.GetVals())
		log.Println(`Length of vals:`, lenVals)
		responsePropNum := current.GetPropNum()

		// If node sent val back, it has already accepted a value
		if lenVals != 0 {
			log.Println(`Current propNum:`, propNum, `responsePropNum:`, responsePropNum)
			if propNum > responsePropNum {
				highestPropNum = responsePropNum
				msgVal = *toAction(current.GetVals()[0])
				setVal = true
			}
		// If node sent back nack, proposal # too low but no accepted val
		} else if current.GetMsgType() == nackPromiseVal {
			if propNum > res
		}
	}

	if setVal {
		return highestPropNum, msgVal
	}
	return propNum, nil */
}

// accept sends accept! messages to other nodes and returns an error if quorum does not respond (might remove error)
func accept(arrVal int64, propNum int64, val Action) error {
	msg := newMessage(network.LocalAddr, nodeID, acceptVal, arrVal, propNum, toProtoPlusOne(lenLog()), val)
	protoMsg := msg.toPMessage()
	request, err := proto.Marshal(protoMsg)
	if err != nil {
		log.Println(`paxos: accept(): Request marshal fail`)
		log.Println(err)
	}

	_, failures := network.Send(request)

	if failures >= network.Quorum() {
		return errors.New("too many failures")
	}

	return nil
}

func bully() {
	msg := newMessage(network.LocalAddr, nodeID, bullyVal, 0, 0, toProtoPlusOne(lenLog()))
	protoMsg := msg.toPMessage()
	request, err := proto.Marshal(protoMsg)
	if err != nil {
		log.Println(`paxos: bully(): Request mashal fail`)
		log.Println(err)
	}

	responses, failures := network.Send(request)

	if failures >= network.Quorum() {
		// do something
	}

	chosenID := nodeID
	longestLog := lenLog()
	for i := range responses {
		var current *pb.PMessage
		err := proto.Unmarshal(responses[i], current)
		if err != nil {
			log.Println(`paxos: bully(): Response unmarshal fail`)
			log.Println(err)
		}

		responseID := current.GetClientID()
		responseLogLen := fromProtoMinusOne(current.GetLogLen())

		if responseLogLen > longestLog {
			chosenID = responseID
		} else if responseLogLen == longestLog {
			if chosenID < responseID {
				chosenID = responseID
			}
		}
	}

	currentLeader = chosenID
}

func promise(data *pb.PMessage) []byte {
	arrVal := data.GetArrVal()
	propNum := data.GetPropNum()
	var msg *Message
	// if array val listed (starting at 1) is equal to the length of the node's log, then the node has a value
	// at that array index (when subtracting 1)
	// ex. if arrVal == 1 and lenLog() == 1, then another node is trying to propose a value for a spot in the log
	// that the current node has already accepted a value for
	if arrVal == lenLog() {
		currentPropNum := promiseMap[arrVal]
		msg = newMessage(network.LocalAddr, nodeID, oldPromiseVal, arrVal, currentPropNum,
			toProtoPlusOne(lenLog()), dataLog[fromProtoMinusOne(arrVal)])
	} else if currentPropNum, ok := promiseMap[arrVal]; ok {
		msg = newMessage(network.LocalAddr, nodeID, nackPromiseVal, arrVal, currentPropNum,
			toProtoPlusOne(lenLog()))
	} else {
		msg = newMessage(network.LocalAddr, nodeID, promiseVal, arrVal, propNum,
			toProtoPlusOne(lenLog()))
	}

	protoMsg := msg.toPMessage()
	responseBytes, err := proto.Marshal(protoMsg)
	if err != nil {
		log.Println(`paxos: promise(): Response marshal fail`)
	}

	return responseBytes
}
