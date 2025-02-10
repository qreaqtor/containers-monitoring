package models

import "time"

type ContainerInfo struct {
	ID     string
	Name   string
	Image  string
	IP     string
	Ports  []string
	Status string

	UpdatedAt time.Time
}

type ContainersInfo struct {
	Containers []ContainerInfo `json:"containers"`
}

func NewContainersInfo(containers []ContainerInfo) ContainersInfo {
	return ContainersInfo{
		Containers: containers,
	}
}
