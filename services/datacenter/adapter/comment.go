package adapter

import (
	"aposervice/domain"
	"crypto/md5"
	"fmt"
	"io"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type AddCommentsResult struct {
	Success int                 `json:"sucess"`
	Exists  []domain.ApoComment `json:"exists"`
}

type GetCommentsQuery struct {
	AppID  string
	Status domain.ApoCommentStatus
	ApoID  string
	Limit  int
	Offset int
	Start  *time.Time
	End    *time.Time
}

// AddComments
func AddComments(comments []domain.ApoComment) (*AddCommentsResult, error) {
	var (
		result AddCommentsResult
		err    error
		count  int
	)
	pool := mgoPool.C("apo_comments")
	now := time.Now()
	for _, c := range comments {
		c.MD5 = MD5(c.Content)
		c.ID = bson.NewObjectId()
		count, err = pool.Find(bson.M{"app_id": c.AppID, "md5": c.MD5}).Count()
		if err != nil {
			return nil, err
		}
		if count > 0 {
			result.Exists = append(result.Exists, c)
			continue
		}

		c.CreateAt = &now
		c.UpdateAt = &now
		c.Status = domain.ApoCommentStatusFree
		if _, err = pool.Upsert(bson.M{"md5": c.MD5}, c); err != nil {
			return nil, err
		}
		result.Success += 1
	}
	return &result, nil
}

// GetComments
func GetComments(query *GetCommentsQuery) ([]domain.ApoComment, error) {
	comments := []domain.ApoComment{}
	queryParam := bson.M{}
	if query.AppID != "" {
		queryParam["app_id"] = query.AppID
	}
	if query.ApoID != "" {
		queryParam["apo_id"] = query.ApoID
	}
	if query.Status != 0 {
		queryParam["status"] = query.Status
	}

	if query.Start != nil || query.End != nil {
		timeRange := bson.M{}
		if query.Start != nil {
			timeRange["$gt"] = query.Start
		}
		if query.End != nil {
			timeRange["$lt"] = query.End
		}
		queryParam["update_time"] = timeRange
	}
	pool := mgoPool.C("apo_comments")

	if err := pool.Find(queryParam).Limit(query.Limit).Skip(query.Offset).All(&comments); err != nil {
		return nil, err
	}

	return comments, nil
}

// Md5
func MD5(content string) string {
	h := md5.New()
	io.WriteString(h, content)
	return fmt.Sprintf("%x", h.Sum(nil))
}
