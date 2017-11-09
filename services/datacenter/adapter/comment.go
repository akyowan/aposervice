package adapter

import (
	"aposervice/domain"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"fxlibraries/loggers"
	"fxlibraries/mongo"
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
		contentStr, _ := json.Marshal(c.Content)
		c.MD5 = MD5(string(contentStr))
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
		if err := pool.Insert(c); err != nil {
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

// DeleteComment
func DeleteComment(id string) error {
	pool := mgoPool.C("apo_comments")
	if !bson.IsObjectIdHex(id) {
		loggers.Error.Printf("DeleteComment invalid object id:%s", id)
		return errors.New("NotFound")
	}
	objID := bson.ObjectIdHex(id)
	if err := pool.RemoveId(objID); err != nil {
		if mongo.NotFound(err) {
			return errors.New("NotFound")
		}
		return err
	}
	return nil
}

// DeleteComments
func DeleteComments(ids []string) error {
	pool := mgoPool.C("apo_comments")
	var objIDs []bson.ObjectId
	for i := range ids {
		if !bson.IsObjectIdHex(ids[i]) {
			loggers.Error.Printf("DeleteComment invalid object id:%s", ids[i])
			continue
		}
		objIDs = append(objIDs, bson.ObjectIdHex(ids[i]))
	}
	queryParam := bson.M{
		"_id": bson.M{
			"$in": objIDs,
		},
	}

	info, err := pool.RemoveAll(queryParam)
	if err != nil {
		return err
	}
	loggers.Debug.Printf("DeleteComments match:%d removed:%d updated:%d", info.Matched, info.Removed, info.Updated)
	return nil
}

// DeleteAppComments
func DeleteAppComments(appID string) error {
	pool := mgoPool.C("apo_comments")
	queryParam := bson.M{"app_id": appID}

	info, err := pool.RemoveAll(queryParam)
	if err != nil {
		return err
	}
	loggers.Debug.Printf("DeleteAppComments match:%d removed:%d updated:%d", info.Matched, info.Removed, info.Updated)
	return nil
}

// UpdateComment
func UpdateComment(id string, comment *domain.ApoComment) (*domain.ApoComment, error) {
	if !bson.IsObjectIdHex(id) {
		loggers.Error.Printf("DeleteComment invalid object id:%s", id)
		return nil, errors.New("NotFound")
	}
	comment.ID = bson.ObjectIdHex(id)
	pool := mgoPool.C("apo_comments")
	var oldComment domain.ApoComment
	if err := pool.FindId(comment.ID).One(&oldComment); err != nil {
		loggers.Error.Printf("UpdateComment FindId id:%s error:%s", comment.ID, err.Error())
		if mongo.NotFound(err) {
			return nil, errors.New("NotFound")
		}
		return nil, err
	}
	if comment.AppID != "" {
		oldComment.AppID = comment.AppID
	}

	contentStr, _ := json.Marshal(comment.Content)
	comment.MD5 = MD5(string(contentStr))
	count, err := pool.Find(bson.M{"app_id": oldComment.AppID, "md5": comment.MD5, "_id": bson.M{"$ne": oldComment.ID}}).Count()
	if err != nil {
		loggers.Error.Printf("UpdateComment Find app_id:%s md5:%s error:%s", oldComment.ApoID, comment.MD5, err.Error())
		return nil, err
	}
	if count > 0 {
		loggers.Error.Printf("UpdateComment app_id:%s md5:%s error:content is exist", oldComment.ApoID, comment.MD5)
		return nil, errors.New("Exist")
	}

	now := time.Now()
	oldComment.MD5 = comment.MD5
	oldComment.Content = comment.Content
	oldComment.UpdateAt = &now
	if err := pool.UpdateId(oldComment.ID, oldComment); err != nil {
		loggers.Error.Printf("UpdateComment UpdateId id:%s error:%s", oldComment.ID, err.Error())
		return nil, err
	}

	return &oldComment, nil
}

// Md5
func MD5(content string) string {
	h := md5.New()
	io.WriteString(h, content)
	return fmt.Sprintf("%x", h.Sum(nil))
}
