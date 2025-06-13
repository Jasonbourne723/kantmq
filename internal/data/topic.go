package data

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"kantmq/internal/conf"
	"kantmq/internal/dto"
	"kantmq/internal/mapping"
)

const (
	FileModePerm = 0644
)

type TopicStorage struct {
	path  string
	data  map[string]dto.TopicInfo
	mutex *sync.RWMutex
}

func NewTopicStorage(bs *conf.Bootstrap) *TopicStorage {

	s := &TopicStorage{
		path:  bs.Storage.Metadata,
		data:  make(map[string]dto.TopicInfo),
		mutex: &sync.RWMutex{},
	}
	s.load()
	return s
}

func (s *TopicStorage) Add(req dto.CreateTopicRequest) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exist := s.data[req.Name]; !exist {
		topicInfo := dto.TopicInfo{}
		if err := mapping.Copy(&topicInfo, &req); err != nil {
			return err
		}
		s.data[req.Name] = topicInfo
		return s.sync(topicInfo)

	} else {
		return nil
	}
}

func (s *TopicStorage) GetTotalTopics() (map[string]dto.TopicInfo, error) {

	s.mutex.RLocker().Lock()
	defer s.mutex.RLocker().Unlock()
	return s.data, nil
}

func (s *TopicStorage) sync(req dto.TopicInfo) error {
	file, err := os.OpenFile(s.path, os.O_RDWR|os.O_APPEND|os.O_CREATE, FileModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	topicBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	lBytes := make([]byte, 2, 3+len(topicBytes))
	binary.BigEndian.PutUint16(lBytes, uint16(len(topicBytes)))
	lBytes = append(lBytes, byte(0))
	lBytes = append(lBytes, topicBytes...)
	_, err = file.Write(lBytes)
	return err
}

func (s *TopicStorage) load() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(s.path, os.O_RDWR|os.O_CREATE, FileModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		lBytes := make([]byte, 2)
		size, err := file.Read(lBytes)
		if err != nil {
			return err
		}
		if size < 2 {
			return nil
		}

		tombstoneByte := make([]byte, 1)
		tombstoneSize, err := file.Read(tombstoneByte)
		if err != nil {
			return err
		}
		if tombstoneSize < 1 {
			return nil
		}
		tombstone := uint8(tombstoneByte[0])
		l := binary.BigEndian.Uint16(lBytes)
		dataBytes := make([]byte, l)
		size, err = file.Read(dataBytes)
		if err != nil {
			return err
		}
		if size < int(l) {
			return nil
		}
		if tombstone == 0 {
			var data dto.TopicInfo
			err = json.Unmarshal(dataBytes, &data)
			if err != nil {
				return err
			}
			s.data[data.Name] = data
		}
	}

}
