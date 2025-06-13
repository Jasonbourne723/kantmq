package mapping

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var DefaultOption = copier.Option{Converters: DefaultConverter}

var DefaultConverter = []copier.TypeConverter{
	TimeToUnixConverter,
	UnixToTimeConverter,
	TimeToTimestampConverter,
	TimestampToTimeConverter,
	float64ToString,
}

var float64ToString = copier.TypeConverter{
	SrcType: float64(0),
	DstType: string(""),
	Fn: func(src interface{}) (interface{}, error) {
		s, ok := src.(float64)
		if !ok {
			return nil, errors.New("src type not matching")
		}
		return fmt.Sprintf("%v", s), nil
	},
}

var TimeToUnixConverter = copier.TypeConverter{
	SrcType: time.Time{},
	DstType: int64(0),
	Fn: func(src interface{}) (interface{}, error) {
		s, ok := src.(time.Time)
		if !ok {
			return nil, errors.New("src type not matching")
		}
		return s.Unix() * 1000, nil
	},
}

var UnixToTimeConverter = copier.TypeConverter{
	SrcType: int64(0),
	DstType: time.Time{},
	Fn: func(src interface{}) (interface{}, error) {
		s, ok := src.(int64)
		if !ok {
			return nil, errors.New("src type not matching")
		}
		return time.Unix(s/1000, 0), nil
	},
}

// time to timestamp. standard library time.Time to google.golang.org/protobuf/types/known/timestamppb.Timestamp
var TimeToTimestampConverter = copier.TypeConverter{
	SrcType: time.Time{},
	DstType: &timestamppb.Timestamp{},
	Fn: func(src interface{}) (interface{}, error) {
		s, ok := src.(time.Time)
		if !ok {
			return nil, errors.New("src type not matching")
		}
		if s.IsZero() {
			return nil, nil
		}
		return timestamppb.New(s), nil
	},
}

// timestamp to time. google.golang.org/protobuf/types/known/timestamppb.Timestamp to standard library time.Time
var TimestampToTimeConverter = copier.TypeConverter{
	SrcType: &timestamppb.Timestamp{},
	DstType: time.Time{},
	Fn: func(src interface{}) (interface{}, error) {
		s, ok := src.(*timestamppb.Timestamp)
		if !ok {
			return nil, errors.New("src type not matching")
		}
		if s == nil {
			return time.Time{}, nil
		}
		return s.AsTime(), nil
	},
}

func Copy(toValue interface{}, fromValue interface{}) (err error) {
	return copier.CopyWithOption(toValue, fromValue, DefaultOption)
}
