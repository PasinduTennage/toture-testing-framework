package client

import "toture-test/consenbench/common"

type ClientOptions struct {
	NodeInfoFile  string // the yaml file containing the ip address of each node, controller port, client port
	inputMessages chan interface{}
}

type Client struct {
	id      int
	network *common.Network // to communicate with the controller
}

func NewClient(options ClientOptions) *Client {
	return &Client{}
}

// initialize the network layer

func (c *Client) NetworkInit() error {
	return nil
}

// respond to different messages from the controller, periodically send machine stats to the controller

func (c *Client) Run() {

}
