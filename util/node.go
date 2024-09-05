package util

type Node interface {
	ExecCmd(cmd string)

	Get_Load(remote_location string, local_location string) error

	Put_Load(local_location string, remote_location string)

	Shut_Down() error
}

type NodeImpl struct {
	Id       int
	Ip       string
	Username string
	HomeDir  string
}

func NewNode() *NodeImpl {
	return &NodeImpl{}
}

func (n *NodeImpl) ExecCmd(cmd string) error {
	// run the given command on the node
	return nil
}

func (n *NodeImpl) Get_Load(remote_location string, local_location string) error {
	// download the file from the remote location to the local location
	return nil
}

func (n *NodeImpl) Put_Load(local_location string, remote_location string) error {
	// upload the file from the local location to the remote location
	return nil
}

func (n *NodeImpl) Shut_Down() error {
	// shut down the node
	return nil
}
