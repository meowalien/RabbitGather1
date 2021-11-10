package broker

import (
	"fmt"
	"log"
)

const DefaultChanSize = 1
const DefaultSubscribeChanSize = 1
const DefaultUnSubscribeChanSize = 1
const DefaultClientQueueSize = 5
const MaximumBroadcastThreads = 1000

type BrokerOptions struct {
	PublishChanSize     int
	SubscribeChanSize   int
	UnSubscribeChanSize int
	ClientQueueSize     int
}

var threadLimiter = make(chan struct{}, MaximumBroadcastThreads)

// broker is a bridge between multiple BrokerClient, it will transfer data between them.
type broker struct {
	// signal to stop broker
	stopCh chan struct{}
	// sent message to all BrokerClient
	publishChan chan [2]interface{}
	// subscribe new *BrokerClient
	subscribeChan chan *BrokerClient
	// Unsubscribe BrokerClient
	unSubscribeChan chan *BrokerClient

	isActive        bool
	clientQueueSize int
}

// NewBroker create a new broker according to given option, will create a default broker if the given option is nil
func NewBroker(option *BrokerOptions) *broker {
	if option == nil {
		option = &BrokerOptions{
			PublishChanSize:     DefaultChanSize,
			SubscribeChanSize:   DefaultSubscribeChanSize,
			UnSubscribeChanSize: DefaultUnSubscribeChanSize,
			ClientQueueSize:     DefaultClientQueueSize,
		}
	}
	return &broker{
		stopCh:          make(chan struct{}),
		publishChan:     make(chan [2]interface{}, option.PublishChanSize),
		subscribeChan:   make(chan *BrokerClient, option.SubscribeChanSize),
		unSubscribeChan: make(chan *BrokerClient, option.UnSubscribeChanSize),
		clientQueueSize: option.ClientQueueSize,
	}
}

// Start should be called before broker use, it starts up the broker
func (b *broker) Start() {
	subs := map[*BrokerClient]struct{}{}
	b.isActive = true
	defer func() { b.isActive = false }()
	for {
		select {
		case <-b.stopCh:
			for msgCh := range subs {
				err := msgCh.Close()
				if err != nil {
					fmt.Printf("error when close %s BrokerClient: %s\n", msgCh.Name, err.Error())
				}
			}
			return
		case msgCh := <-b.subscribeChan:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unSubscribeChan:
			delete(subs, msgCh)
		case m := <-b.publishChan:
			msg := m[0]
			allExcept := m[1].([]*BrokerClient)
			for msgCh := range subs {
				doTransfer := func(bk *BrokerClient) {
					if allExcept != nil {
						for _, exceptMsgCh := range allExcept {
							if exceptMsgCh == bk {
								return
							}
						}
					}
					if bk.Filter == nil || bk.Filter(msg) {
						if !b.isActive {
							return
						}
						select {
						case bk.C <- msg:
						default:
						}
					}
				}
				threadLimiter <- struct{}{}
				doTransfer(msgCh)
				<-threadLimiter
			}
		}
	}
}

// Stop will stop the broker
func (b *broker) Stop() {
	b.isActive = false
	close(b.stopCh)
}

// Subscribe will create a new BrokerClient which Listen on new published message
func (b *broker) Subscribe(name string, filter Filter) *BrokerClient {
	if !b.isActive {
		panic("the broker is not active, please start it up.")
	}
	msgCh := &BrokerClient{C: make(chan interface{}, b.clientQueueSize), Filter: filter, Name: name}
	b.subscribeChan <- msgCh
	return msgCh
}

// Unsubscribe will make broker stop sending new message to the given BrokerClient cnd close the C channel.
func (b *broker) Unsubscribe(msgCh *BrokerClient) {
	if !b.isActive {
		panic("the broker is not active, please start it up.")
	}
	b.unSubscribeChan <- msgCh
	err := msgCh.Close()
	if err != nil {
		log.Printf("error when close BrokerClient: %s\n",err.Error())
	}
}

// Publish will broadcast the message to all subscribed BrokerClient.
func (b *broker) Publish(msg interface{}, except ...*BrokerClient) {
	if !b.isActive {
		panic("the broker is not active, please start it up.")
	}
	//fmt.Println("Publish msg, except : ",msg," , ", except)
	b.publishChan <- [2]interface{}{msg, except}
}

// Filter will filter the messages input and return true if the message you want to pickup.
type Filter func(interface{}) bool

type BrokerClient struct {
	// New messages will be received through C
	C      chan interface{}
	Filter Filter
	Name   string
}

func (b *BrokerClient) Close() error {
	close(b.C)
	return nil
}
