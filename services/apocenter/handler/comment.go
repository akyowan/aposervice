package handler

import (
	"aposervice/config"
	"aposervice/domain"
	"aposervice/services/apocenter/adapter"
	"fxlibraries/errors"
	"fxlibraries/httpserver"
	"fxlibraries/loggers"
	"strings"
)

func GetAppComments(req *httpserver.Request) *httpserver.Response {
	appID := req.UrlParams["appID"]
	if appID == "" {
		loggers.Warn.Printf("GetAppComments no app id")
		return httpserver.NewResponseWithError(errors.NewBadRequest("No app id input"))
	}
	apoID := req.QueryParams.Get("apo_id")
	if apoID == "" {
		loggers.Warn.Printf("GetAppComments app id:%s no apo id", appID)
		return httpserver.NewResponseWithError(errors.NewBadRequest("No apo id input"))
	}
	param := domain.ApoComment{
		AppID: appID,
		ApoID: apoID,
		IP:    strings.Split(req.RemoteAddr, ":")[0],
	}

	ipUsedCount, err := adapter.IpUsedCount(param.IP, 1)
	if err != nil {
		loggers.Warn.Printf("GetAppComments check app id:%s apo id:%s ip:%s used count error %s", appID, apoID, param.IP, err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}
	if ipUsedCount > config.Conf.CommentDayMaxCount {
		return httpserver.NewResponseWithError(errors.NewBadRequest("Ip used today"))
	}

	apoComment, err := adapter.GetAndLockComment(&param)
	if err != nil {
		adapter.IpUsedCount(param.IP, -1)
		loggers.Warn.Printf("GetAppComments get app id:%s apo id:%s ip:%s comment error %s", appID, apoID, param.IP, err.Error())
		return httpserver.NewResponseWithError(errors.NewNotFound(err.Error()))
	}

	resp := httpserver.NewResponse()
	resp.Data = apoComment
	return resp
}

func ReportComment(req *httpserver.Request) *httpserver.Response {
	var param domain.ApoComment
	if err := req.Parse(&param); err != nil {
		loggers.Warn.Printf("ReportComment error %s", err.Error())
		return httpserver.NewResponseWithError(errors.ParameterError)
	}
	if param.ID == "" {
		loggers.Warn.Printf("ReportComment no comment id")
		return httpserver.NewResponseWithError(errors.NewBadRequest("No comment id input"))
	}
	comment, err := adapter.GetComment(param.ID)
	if err != nil {
		loggers.Warn.Print("ReportComment get comment id:%s error:%s", param.ID.Hex(), err.Error())
		return httpserver.NewResponseWithError(errors.NewNotFound("Comment not found"))
	}
	if comment.Status != domain.ApoCommentStatusLocked {
		loggers.Warn.Printf("ReportComment invalid comment status:%d", comment.Status)
		if comment.Status == domain.ApoCommentStatusFree {
			return httpserver.NewResponseWithError(errors.NewBadRequest("Comment timeout recycle"))
		} else {
			return httpserver.NewResponseWithError(errors.NewBadRequest("Invalid comment status"))
		}
	}
	if param.Errno != 0 {
		param.Status = domain.ApoCommentStatusFree
	} else {
		param.Status = domain.ApoCommentStatusUsed
	}
	param.IP = strings.Split(req.RemoteAddr, ":")[0]

	if _, err := adapter.UpdateComment(&param); err != nil {
		loggers.Warn.Printf("ReportComment error %s", err.Error())
		return httpserver.NewResponseWithError(errors.InternalServerError)
	}

	return httpserver.NewResponse()
}
