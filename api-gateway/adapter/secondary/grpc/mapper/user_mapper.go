package mapper

import (
	userpb "github.com/adityakw90/service-user-proto/gen/go/user"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
)

// UserFromProto converts a gRPC User protobuf to a domain User.
func UserFromProto(u *userpb.User) *model.User {
	if u == nil {
		return nil
	}
	return &model.User{
		UID:       u.Uid,
		Username:  u.Username,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: TimestampToTime(u.CreatedAt),
		UpdatedAt: TimestampToTime(u.UpdatedAt),
		DeletedAt: TimestampToTimePtr(u.DeletedAt),
	}
}

// UsersFromProto converts a gRPC User list to domain Users with metadata.
func UsersFromProto(items []*userpb.User, total int64, page, limit int32) *model.Users {
	users := make([]model.User, 0, len(items))
	for _, u := range items {
		if user := UserFromProto(u); user != nil {
			users = append(users, *user)
		}
	}
	return &model.Users{
		Items: users,
		Meta: model.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: calculatePages(total, limit),
		},
	}
}

// calculatePages calculates total number of pages.
func calculatePages(total int64, limit int32) int32 {
	if limit <= 0 {
		return 0
	}
	pages := int(total) / int(limit)
	if int(total)%int(limit) > 0 {
		pages++
	}
	return int32(pages)
}

// ProfileFromProto converts a gRPC Profile protobuf to a domain Profile.
func ProfileFromProto(p *userpb.Profile) *model.Profile {
	if p == nil {
		return nil
	}
	profile := &model.Profile{
		UID:      p.Uid,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Bio:       p.Bio,
	}
	if p.Avatar != nil {
		profile.Avatar = &model.Avatar{
			UID: p.Avatar.Uid,
			URL: p.Avatar.Url,
		}
	}
	return profile
}

// DeviceFromProto converts a gRPC Device protobuf to a domain Device.
func DeviceFromProto(d *userpb.Device) *model.Device {
	if d == nil {
		return nil
	}
	return &model.Device{
		DeviceUID:    d.DeviceUid,
		DeviceName:   d.DeviceName,
		IPAddress:    d.IpAddress,
		LastActiveAt: TimestampToTime(d.LastActiveAt),
		CreatedAt:    TimestampToTime(d.CreatedAt),
		RevokedAt:    TimestampToTimePtr(d.RevokedAt),
	}
}

// DevicesFromProto converts a gRPC Device list to domain Devices with metadata.
func DevicesFromProto(items []*userpb.Device, total int64, page, limit int32) *model.Devices {
	devices := make([]model.Device, 0, len(items))
	for _, d := range items {
		if device := DeviceFromProto(d); device != nil {
			devices = append(devices, *device)
		}
	}
	return &model.Devices{
		Items: devices,
		Meta: model.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: calculatePages(total, limit),
		},
	}
}
