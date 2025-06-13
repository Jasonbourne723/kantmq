package service

import (
	"context"

	pb "kantmq/api/kantmq/v1"
	"kantmq/internal/data"
	"kantmq/internal/dto"
	"kantmq/internal/mapping"
)

type ManagementService struct {
	pb.UnimplementedManagementServer
	topicStorage *data.TopicStorage
}

func NewManagementService(topicStorage *data.TopicStorage) *ManagementService {
	return &ManagementService{
		topicStorage: topicStorage,
	}
}

func (s *ManagementService) CreateTopic(ctx context.Context, req *pb.CreateTopicRequest) (*pb.Empty, error) {

	createReq := dto.CreateTopicRequest{}
	if err := mapping.Copy(&createReq, req); err != nil {
		return nil, err
	}
	if err := s.topicStorage.Add(createReq); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
func (s *ManagementService) DeleteTopic(ctx context.Context, req *pb.DeleteTopicRequest) (*pb.Empty, error) {

	return &pb.Empty{}, nil
}
func (s *ManagementService) GetTopics(ctx context.Context, req *pb.GetTopicsRequest) (*pb.GetTopicsResponse, error) {

	infos, err := s.topicStorage.GetTotalTopics()
	if err != nil {
		return nil, err
	}
	items := make([]*pb.TopicInfo, 0)
	for _, v := range infos {
		item := pb.TopicInfo{}
		if err := mapping.Copy(&item, &v); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	resp := &pb.GetTopicsResponse{
		Items: items,
	}
	return resp, nil
}
