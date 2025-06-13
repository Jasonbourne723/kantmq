package dto

type PublishRequest struct {
	Topic     string
	Partition *int
	Metadata  map[string]string
	Value     []byte
}

type MessageInfo struct {
	Topic     string
	Partition int
	Metadata  map[string]string
	Value     []byte
}
