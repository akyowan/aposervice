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
