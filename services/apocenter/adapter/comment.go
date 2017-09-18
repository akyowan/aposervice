package adapter

import (
	"aposervice/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// GetAndLockComment
func GetAndLockComment(param *domain.ApoComment) (*domain.ApoComment, error) {
	var comment domain.ApoComment
	pool := mgoPool.C("apo_comments")
	now := time.Now()
	change := mgo.Change{
		Update:    bson.M{"status": domain.ApoCommentStatusLocked, "ip": param.IP, "update_time": &now},
		ReturnNew: false,
	}
	query := bson.M{
		"status": domain.ApoCommentStatusFree,
		"app_id": param.AppID,
	}
	if _, err := pool.Find(query).Apply(change, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// UpdateComment
func UpdateComment(param *domain.ApoComment) (*domain.ApoComment, error) {
	var comment domain.ApoComment
	pool := mgoPool.C("apo_comments")
	now := time.Now()
	update := bson.M{"status": param.Status, "update_time": &now}
	if param.IP != "" {
		update["ip"] = param.IP
	}
	if param.ApoID != "" {
		update["apo_id"] = param.ApoID
	}
	change := mgo.Change{
		Update:    update,
		ReturnNew: true,
	}
	if _, err := pool.FindId(param.ID).Apply(change, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}
