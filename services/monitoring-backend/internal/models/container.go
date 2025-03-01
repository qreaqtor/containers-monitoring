package models

import "time"

type ContainerInfo struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Image  string   `json:"image"`
	IP     string   `json:"ip"`
	Ports  []string `json:"ports"`
	Status string   `json:"status"`

	UpdatedAt time.Time `json:"updated_at"`
}

type ContainersInfo struct {
	Containers []ContainerInfo `json:"containers"`
}

func NewContainersInfo(containers []ContainerInfo) ContainersInfo {
	return ContainersInfo{
		Containers: containers,
	}
}
