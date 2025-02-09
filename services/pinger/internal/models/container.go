package models

type ContainerInfo struct {
	ID     string
	Name   string
	Image  string
	IPv4   string
	Ports  []string
	State  string
	Status string
}
