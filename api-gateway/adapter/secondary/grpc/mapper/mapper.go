package mapper

import (
	"time"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TimestampToTime converts a protobuf Timestamp to time.Time.
func TimestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

// TimestampToTimePtr converts a protobuf Timestamp to *time.Time.
func TimestampToTimePtr(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

// TimeToTimestamp converts time.Time to protobuf Timestamp.
func TimeToTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
