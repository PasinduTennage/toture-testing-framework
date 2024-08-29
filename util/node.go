package util

type Node struct {
	ID int
	IP string
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) GetID() int {
	return n.ID
}

func (n *Node) GetIP() string {
	return n.IP

}

func (n *Node) SetID(id int) {
	n.ID = id
}

func (n *Node) SetIP(ip string) {
	n.IP = ip
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
