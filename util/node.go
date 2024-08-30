package util

type Node struct {
	Id       int
	Ip       string
	Username string
	HomeDir  string
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) SSH(cmd string) error {
	// run the given command on the node
	return nil
}

func (n *Node) Get_Load(remote_location string, local_location string) error {
	// download the file from the remote location to the local location
	return nil
}

func (n *Node) Put_Load(local_location string, remote_location string) error {
	// upload the file from the local location to the remote location
	return nil
}

func (n *Node) Shut_Down() error {
	// shut down the node
	return nil
}
