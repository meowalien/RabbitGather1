package util

const DEFULT_CHAN_SIZE = 1
const MaximumBroadcastThreads = 1000

type BrokerOptions struct {
	PublishChanSize     int
	SubscribeChanSize   int
	UnSubscribeChanSize int
}

var threadLimiter = make(chan struct{}, MaximumBroadcastThreads)

// Broker is a bridge between multiple BrokerClient, it will transfer data between them.
type Broker struct {
	// signal to stop Broker
	stopCh chan struct{}
	// sent message to all BrokerClient
	publishChan chan [2]interface{}
	// subscribe new *BrokerClient
	subscribeChan chan *BrokerClient
	// Unsubscribe BrokerClient
	unSubscribeChan chan *BrokerClient

	isActive bool
}

// NewBroker create a new Broker according to given option, will create a default Broker if the given option is nil
func NewBroker(option *BrokerOptions) *Broker {
	if option == nil {
		option = &BrokerOptions{
			PublishChanSize:     DEFULT_CHAN_SIZE,
			SubscribeChanSize:   1,
			UnSubscribeChanSize: 1,
		}
	}
	return &Broker{
		stopCh:          make(chan struct{}),
		publishChan:     make(chan [2]interface{}, option.PublishChanSize),
		subscribeChan:   make(chan *BrokerClient, option.SubscribeChanSize),
		unSubscribeChan: make(chan *BrokerClient, option.UnSubscribeChanSize),
	}
}

// Start should be called before Broker use, it starts up the Broker
func (b *Broker) Start() {
	subs := map[*BrokerClient]struct{}{}
	b.isActive = true
	defer func() { b.isActive = false }()
	for {
		select {
		case <-b.stopCh:
			for msgCh := range subs {
				close(msgCh.C)
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
							if exceptMsgCh == msgCh {
								return
							}
						}
					}
					if msgCh.Filter == nil || msgCh.Filter(msg) {
						// msgCh is buffered, use non-blocking send to protect the broker:
						select {
						case msgCh.C <- msg:
						default:
						}
					}
				}
				threadLimiter <- struct{}{}
				go func(msgCh *BrokerClient) {
					doTransfer(msgCh)
					<-threadLimiter
				}(msgCh)
			}
		}
	}
}

// Stop will stop the Broker
func (b *Broker) Stop() {
	close(b.stopCh)
	b.isActive = false
}

// Subscribe will create a new BrokerClient which Listen on new published message
func (b *Broker) Subscribe(filter Filter) *BrokerClient {
	if !b.isActive {
		panic("the broker is not active, please start it up.")
	}
	msgCh := &BrokerClient{C: make(chan interface{}, 5), Filter: filter}
	b.subscribeChan <- msgCh
	return msgCh
}

// Unsubscribe will make Broker stop sending new message to the given BrokerClient cnd close the C channel.
func (b *Broker) Unsubscribe(msgCh *BrokerClient) {
	if !b.isActive {
		panic("the broker is not active, please start it up.")
	}
	b.unSubscribeChan <- msgCh
	close(msgCh.C)
}

// Publish will broadcast the message to all subscribed BrokerClient.
func (b *Broker) Publish(msg interface{}, except ...*BrokerClient) {
	if !b.isActive {
		panic("the broker is not active, please start it up.")
	}
	b.publishChan <- [2]interface{}{msg, except}
}

// Filter will filter the messages input and return true if the message you want to pickup.
type Filter func(interface{}) bool

type BrokerClient struct {
	// New messages will be received through C
	C      chan interface{}
	Filter Filter
}
