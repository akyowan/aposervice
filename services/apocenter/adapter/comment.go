package adapter

import (
	"aposervice/domain"
	"fmt"
	"fxlibraries/loggers"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var IPCommentCountKey string

func init() {
	IPCommentCountKey = "APOCENTER.IPCommentCount"
}

// GetAndLockComment
func GetAndLockComment(param *domain.ApoComment) (*domain.ApoComment, error) {
	var comment domain.ApoComment

	pool := mgoPool.C("apo_comments")
	now := time.Now()
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"status": domain.ApoCommentStatusLocked, "ip": param.IP, "update_time": &now}},
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
		Update:    bson.M{"$set": update},
		ReturnNew: true,
	}
	if _, err := pool.FindId(param.ID).Apply(change, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// IpUsedCount
func IpUsedCount(IP string, count int64) (int64, error) {
	key := fmt.Sprintf("%s.%s", IPCommentCountKey, IP)
	result := redisPool.IncrBy(key, count)
	if result.Err() != nil {
		loggers.Error.Printf("IpUsedCount error %s", result.Err().Error())
		return 0, result.Err()
	}
	if result.Val() <= 1 {
		if err := redisPool.Expire(key, time.Hour*24).Err(); err != nil {
			return 0, err
		}
	}
	return result.Val(), nil
}
