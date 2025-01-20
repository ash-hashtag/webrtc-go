package main

import (
	"fmt"
	"time"

	"github.com/pion/webrtc/v3"
	"math/rand"
)

func main() {

	iceServers := []webrtc.ICEServer{

		{
			URLs: []string{"stun:a.relay.metered.ca:80"},
		},
		{
			URLs:       []string{"turn:a.relay.metered.ca:80"},
			Username:   "6710e0989d4d453b3b1129b0",
			Credential: "IsdQrRbxROrDzxSW",
		},
		{
			URLs:       []string{"turn:a.relay.metered.ca:80?transport=tcp"},
			Username:   "6710e0989d4d453b3b1129b0",
			Credential: "IsdQrRbxROrDzxSW",
		},
		{
			URLs:       []string{"turn:a.relay.metered.ca:443"},
			Username:   "6710e0989d4d453b3b1129b0",
			Credential: "IsdQrRbxROrDzxSW",
		},
		{
			URLs:       []string{"turn:a.relay.metered.ca:443?transport=tcp"},
			Username:   "6710e0989d4d453b3b1129b0",
			Credential: "IsdQrRbxROrDzxSW",
		},
	}

	iceServers = []webrtc.ICEServer{

		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
		{
			URLs:       []string{"turn:192.168.1.196:3478"},
			Username:   "1699024182:earth:user",
			Credential: "Q23fb6lxKwEL+YqGT0/vuJT1Ykc=",
		},
	}
	offerConfig := webrtc.Configuration{
		ICEServers:         iceServers,
		ICETransportPolicy: webrtc.ICETransportPolicyRelay,
	}
	offerApi := webrtc.NewAPI()
	offerConn, err := offerApi.NewPeerConnection(offerConfig)
	if err != nil {
		panic(err)
	}

	offerConn.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		fmt.Printf("offer ice connection state changed %v\n", state)
	})

	channel, err := offerConn.CreateDataChannel("data", &webrtc.DataChannelInit{})
	if err != nil {
		panic(err)
	}

	channel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("offer received data channel message %v\n", string(msg.Data))
	})

	fmt.Printf("offer Created a DataChannel %v\n", channel.Label())

	offer, err := offerConn.CreateOffer(&webrtc.OfferOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("OFFER \n\n\nsdp: %v\ntype: %v\n\n\n", offer.SDP, offer.Type)

	err = offerConn.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	answerConfig := webrtc.Configuration{
		ICEServers: iceServers,
	}
	answerApi := webrtc.NewAPI()
	answerConn, err := answerApi.NewPeerConnection(answerConfig)
	if err != nil {
		panic(err)
	}

	answerConn.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		fmt.Printf("answer ice connection state changed %v\n", state)
	})

	answerConn.OnDataChannel(func(channel *webrtc.DataChannel) {
		if channel == nil {
			return
		}
		fmt.Printf("received data channel %v\n", channel.Label())
		channel.OnOpen(func() {
			fmt.Printf("opened data channel %v\n", channel.Label())
			for i := 0; i < 5; i++ {
				dur := time.Duration(3*i) * time.Second
				time.AfterFunc(dur, func() {
					s := fmt.Sprintf("Hello from Answer %v\n", rand.Int())
					fmt.Printf("Sent %s", s)
					channel.SendText(s)
				})
			}
		})

		channel.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("received data channel message %v\n", string(msg.Data))
		})
	})

	err = answerConn.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	answer, err := answerConn.CreateAnswer(&webrtc.AnswerOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("ANSWER \n\n\nsdp: %v\ntype: %v\n\n\n", answer.SDP, answer.Type)
	err = answerConn.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	err = offerConn.SetRemoteDescription(answer)
	if err != nil {
		panic(err)
	}

	offerConn.OnICECandidate(func(cand *webrtc.ICECandidate) {
		if cand == nil {
			return
		}
		c := cand.ToJSON()
		fmt.Printf("\n\noffer new ice cand: %v index: %v mid: %v ufrag: %v\n\n", c.Candidate, c.SDPMLineIndex, c.SDPMid, c.UsernameFragment)
		err = answerConn.AddICECandidate(c)
		if err != nil {
			panic(err)
		}

	})
	answerConn.OnICECandidate(func(cand *webrtc.ICECandidate) {
		if cand == nil {
			return
		}
		c := cand.ToJSON()
		fmt.Printf("\n\nanswer new ice cand: %v index: %v mid: %v ufrag: %v\n\n", c.Candidate, c.SDPMLineIndex, c.SDPMid, c.UsernameFragment)
		err = offerConn.AddICECandidate(c)
		if err != nil {
			panic(err)
		}
	})

	time.Sleep(time.Second * 15)

	answerConn.Close()
	offerConn.Close()
}
