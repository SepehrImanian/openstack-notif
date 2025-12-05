package openstack

type ServerListResponse struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	TaskState  string `json:"task_state"`
	VMState    string `json:"vm_state"`
	PowerState int    `json:"power_state"`
	Progress   int    `json:"progress"`
}
