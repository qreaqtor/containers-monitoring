package models

type ContainerInfo struct {
	ID     string
	Name   string
	Image  string
	IP     string
	Ports  []string
	Status string
}

type ContainersMsg struct {
	Containers []ContainerInfo `json:"containers"`
}

func NewContainersMsg(containers []ContainerInfo) ContainersMsg {
	return ContainersMsg{
		Containers: containers,
	}
}
