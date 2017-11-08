package handler

import (
	"aposervice/domain"
	"aposervice/services/datacenter/adapter"
	"fxlibraries/errors"
	"fxlibraries/httpserver"
	"fxlibraries/loggers"
)

// AddComments 添加评论到数据中心
func AddComments(req *httpserver.Request) *httpserver.Response {
	var comments []domain.ApoComment
	if err := req.Parse(&comments); err != nil {
		loggers.Warn.Printf("AddComments Parse comments error %s", err.Error())
		return httpserver.NewResponseWithError(errors.NewBadRequest("INVALID COMMENTS INPUT"))
	}

	result, err := adapter.AddComments(comments)
	if err != nil {
		loggers.Error.Printf("AddComments Add comments error %s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}

	resp := httpserver.NewResponse()
	resp.Data = result

	return resp
}

// GetComments 从数据中心获取评论
func GetComments(req *httpserver.Request) *httpserver.Response {
	query := adapter.GetCommentsQuery{
		Status: 0,
	}

	query.ApoID = req.QueryParams.Get("apo_id")
	query.AppID = req.QueryParams.Get("app_id")

	comments, err := adapter.GetComments(&query)
	if err != nil {
		loggers.Error.Printf("GetComments error %s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}

	resp := httpserver.NewResponse()
	resp.Data = comments
	return resp
}

// DeleteComment 删除单个评论
func DeleteComment(req *httpserver.Request) *httpserver.Response {
	id := req.UrlParams["id"]
	if err := adapter.DeleteComment(id); err != nil {
		if err.Error() == "NotFound" {
			return httpserver.NewResponseWithError(errors.NotFound)
		}
		loggers.Error.Printf("DeleteComment id:%s error:%s", id, err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	return httpserver.NewResponse()
}

// DeleteComments 批量删除评论
func DeleteComments(req *httpserver.Request) *httpserver.Response {
	var ids []string
	if err := req.Parse(&ids); err != nil {
		loggers.Warn.Printf("DeleteComments invalid param")
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	if err := adapter.DeleteComments(ids); err != nil {
		loggers.Error.Printf("DeleteComments error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	return httpserver.NewResponse()
}

// DeleteAppComments 删除App所有评论
func DeleteAppComments(req *httpserver.Request) *httpserver.Response {
	appID := req.UrlParams["appID"]
	if appID == "" {
		loggers.Warn.Printf("DeleteAppComments no appID input")
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	if err := adapter.DeleteAppComments(appID); err != nil {
		loggers.Warn.Printf("DeleteAppComments error:%s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	return httpserver.NewResponse()
}

// UpdateComment
func UpdateComment(req *httpserver.Request) *httpserver.Response {
	id := req.UrlParams["id"]
	var comment domain.ApoComment
	if err := req.Parse(&comment); err != nil {
		loggers.Warn.Printf("UpdateComment invalid param")
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	newComment, err := adapter.UpdateComment(id, &comment)
	if err != nil {
		if err.Error() == "Exist" {
			loggers.Warn.Printf("UpdateComment error:%s", err.Error())
			return httpserver.NewResponseWithError(errors.NewNotFound("Content is exist"))
		}
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	resp := httpserver.NewResponse()
	resp.Data = newComment
	return resp
}
