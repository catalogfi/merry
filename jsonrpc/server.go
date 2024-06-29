package jsonrpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/catalogfi/merry"
	"github.com/fvbock/endless"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	logger    *zap.Logger
	engine    *gin.Engine
	closeFunc func() error
}

func NewServer(logger *zap.Logger, engine *gin.Engine, handler Handler) *Server {
	engine.POST("/", func(context *gin.Context) {
		var jsonrpcReq merry.Request
		if err := context.ShouldBindJSON(&jsonrpcReq); err != nil {
			context.JSON(http.StatusBadRequest, merry.NewResponse("-1", nil, merry.NewInvalidRequest(err)))
			return
		}
		switch jsonrpcReq.Method {
		case MethodMerryVersion:
			data, err := json.Marshal(handler.MerryVersion())
			if err != nil {
				context.JSON(http.StatusBadRequest, merry.NewResponse(jsonrpcReq.ID, nil, merry.NewInvalidRequest(err)))
				return
			}
			context.JSON(http.StatusOK, merry.NewResponse(jsonrpcReq.ID, data, nil))
		case MethodMerryUpdate:
			if err := handler.MerryUpdate(context.Request.Context()); err != nil {
				context.JSON(http.StatusServiceUnavailable, merry.NewResponse(jsonrpcReq.ID, nil, merry.NewInvalidParams(err)))
				return
			}
			context.JSON(http.StatusCreated, merry.NewResponse(jsonrpcReq.ID, nil, nil))
		case MethodMerryFund:
			var req RequestFund
			if err := json.Unmarshal(jsonrpcReq.Params, &req); err != nil {
				context.JSON(http.StatusBadRequest, merry.NewResponse(jsonrpcReq.ID, nil, merry.NewInvalidParams(err)))
				return
			}
			if err := handler.MerryFund(context.Request.Context(), req); err != nil {
				context.JSON(http.StatusBadRequest, merry.NewResponse(jsonrpcReq.ID, nil, merry.NewInvalidRequest(err)))
				return
			}
			context.JSON(http.StatusCreated, merry.NewResponse(jsonrpcReq.ID, nil, nil))
		case MethodMerryArbitrum:
			var req RequestRelay
			if err := json.Unmarshal(jsonrpcReq.Params, &req); err != nil {
				context.JSON(http.StatusBadRequest, merry.NewResponse(jsonrpcReq.ID, nil, merry.NewInvalidParams(err)))
				return
			}
			resp, err := handler.MerryRelay(context.Request.Context(), "arbitrum", req)
			if err != nil {
				context.JSON(http.StatusBadRequest, merry.NewResponse(jsonrpcReq.ID, nil, merry.NewInvalidRequest(err)))
				return
			}
			context.JSON(http.StatusOK, resp)
		case MethodMerryEthereum:
			var req RequestRelay
			if err := json.Unmarshal(jsonrpcReq.Params, &req); err != nil {
				context.JSON(http.StatusBadRequest, merry.NewResponse(jsonrpcReq.ID, nil, merry.NewInvalidParams(err)))
				return
			}
			resp, err := handler.MerryRelay(context.Request.Context(), "ethereum", req)
			if err != nil {
				context.JSON(http.StatusBadRequest, merry.NewResponse(jsonrpcReq.ID, nil, merry.NewInvalidRequest(err)))
				return
			}
			context.JSON(http.StatusOK, resp)
		case MethodMerryBitcoin:
			var req RequestRelay
			if err := json.Unmarshal(jsonrpcReq.Params, &req); err != nil {
				context.JSON(http.StatusBadRequest, merry.NewResponse(jsonrpcReq.ID, nil, merry.NewInvalidParams(err)))
				return
			}
			resp, err := handler.MerryRelay(context.Request.Context(), "bitcoin", req)
			if err != nil {
				context.JSON(http.StatusBadRequest, merry.NewResponse(jsonrpcReq.ID, nil, merry.NewInvalidRequest(err)))
				return
			}
			context.JSON(http.StatusOK, resp)
		default:
			context.JSON(http.StatusNotFound, merry.NewResponse(jsonrpcReq.ID, nil, merry.NewMethodNotFound()))
			return
		}
	})
	return &Server{logger: logger, engine: engine}
}

func (server *Server) Start(port int) error {
	server.engine.Use(cors.Default())
	httpServer := endless.NewServer(fmt.Sprintf(":%v", port), server.engine)
	server.closeFunc = httpServer.Close
	return httpServer.ListenAndServe()
}

func (server *Server) Close() error {
	if server.closeFunc != nil {
		return server.closeFunc()
	}
	return nil
}
