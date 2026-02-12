package mapper

import (
	"time"

	"github.com/adityakw90/service-user-proto/gen/go/user_file"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
)

// UserFileFromProto converts proto UserFile to domain model.
func UserFileFromProto(pb *user_file.UserFile) *model.UserFile {
	f := &model.UserFile{
		UID:        pb.Uid,
		UserUID:    pb.UserUid,
		FileType:   pb.FileType,
		FileName:   pb.FileName,
		FilePath:   pb.FilePath,
		MimeType:   pb.MimeType,
		SizeBytes:  pb.SizeBytes,
		Visibility: visibilityToInt32(pb.Visibility),
		CreatedAt:  time.Time{},
	}

	if pb.CreatedAt != nil {
		f.CreatedAt = pb.CreatedAt.AsTime()
	}

	if pb.Thumbnail != nil && pb.Thumbnail.Fields != nil {
		if uid, ok := pb.Thumbnail.Fields["uid"]; ok && uid != nil {
			if url, ok := pb.Thumbnail.Fields["url"]; ok && url != nil {
				f.Thumbnail = &model.Thumbnail{
					UID: uid.GetStringValue(),
					URL:  url.GetStringValue(),
				}
			}
		}
	}

	return f
}

// visibilityToInt32 converts string visibility to int32.
func visibilityToInt32(v string) int32 {
	if v == "public" || v == "1" {
		return model.UserFileVisibilityPublic
	}
	return model.UserFileVisibilityPrivate
}

// UserFilesFromProto converts proto list response to domain model.
func UserFilesFromProto(items []*user_file.UserFile, total int64, page, limit int32) *model.UserFiles {
	files := make([]model.UserFile, len(items))
	for i, item := range items {
		files[i] = *UserFileFromProto(item)
	}

	return &model.UserFiles{
		Items: files,
		Meta: model.PaginationMeta{
			Total: total,
			Page:  page,
			Limit: limit,
			Pages: calculatePages(total, limit),
		},
	}
}
