package dto

type CreateTopicRequest struct {
	Name              string `json:"Name"`
	NumberOfPartition int    `json:"NumberOfPartition"`
	NumberOfReplicas  int    `json:"NumberOfReplicas"`
}

type TopicInfo struct {
	Name              string `json:"Name"`
	NumberOfPartition int    `json:"NumberOfPartition"`
	NumberOfReplicas  int    `json:"NumberOfReplicas"`
}
