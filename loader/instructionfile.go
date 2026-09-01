package loader

// instruction files are optional for a zone
type instructionFileEntry struct {
	Type        string `json:"type"`
	ObjectId    string `json:"object_id"`
	MobileId    string `json:"mobile_id"`
	RoomId      string `json:"room_id"`
	InstanceMax int    `json:"instance_max"`
}
