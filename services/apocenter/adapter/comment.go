package adapter

import (
	"aposervice/config"
	"aposervice/domain"
	"fmt"
	"fxlibraries/loggers"
	"fxlibraries/mongo"
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
		Update:    bson.M{"$set": update, "$push": bson.M{"errnos": param.Errno}},
		ReturnNew: true,
	}
	if _, err := pool.FindId(param.ID).Apply(change, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetComment
func GetComment(id bson.ObjectId) (*domain.ApoComment, error) {
	var comment domain.ApoComment
	pool := mgoPool.C("apo_comments")
	if err := pool.FindId(id).One(&comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// RecycleComment
func RecycleComment() error {
	pool := mgoPool.C("apo_comments")
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":      domain.ApoCommentStatusFree,
			"update_time": &now,
		},
	}
	query := bson.M{
		"status": domain.ApoCommentStatusLocked,
		"update_time": bson.M{
			"$lt": now.Add(-(config.Conf.CommentTimeout * 1000000000)),
		},
	}
	info, err := pool.UpdateAll(query, update)
	if err != nil {
		if mongo.NotFound(err) {
			loggers.Info.Printf("RecycleComment no comment need recyclee")
			return nil
		}
		loggers.Info.Printf("RecycleComment error:%s", err.Error())
		return nil
	}
	loggers.Info.Printf("RecycleComment count:%d", info.Updated)

	return nil
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
