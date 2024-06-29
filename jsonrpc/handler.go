package jsonrpc

import (
	"context"
	"encoding/json"

	"github.com/catalogfi/merry"
)

const (
	MethodMerryFund     = "merry_fund"
	MethodMerryUpdate   = "merry_update"
	MethodMerryVersion  = "merry_version"
	MethodMerryBitcoin  = "merry_bitcoin"
	MethodMerryArbitrum = "merry_arbitrum"
	MethodMerryEthereum = "merry_ethereum"
)

type Handler interface {
	MerryFund(ctx context.Context, Request RequestFund) error
	MerryUpdate(ctx context.Context) error
	MerryVersion() ResponseVersion
	MerryRelay(ctx context.Context, service string, Request RequestRelay) (merry.Response, error)
}

type merryHandler struct {
	merry *merry.Merry
}

func NewMerryHandler(merry *merry.Merry) Handler {
	return merryHandler{merry}
}

func (handler merryHandler) MerryFund(ctx context.Context, request RequestFund) error {
	return handler.merry.Fund(request.To)
}

func (handler merryHandler) MerryUpdate(ctx context.Context) error {
	return merry.Update()
}

func (handler merryHandler) MerryVersion() ResponseVersion {
	return ResponseVersion{Version: handler.merry.Version, Commit: handler.merry.Commit, Date: handler.merry.Date}
}

func (handler merryHandler) MerryRelay(ctx context.Context, service string, request RequestRelay) (merry.Response, error) {
	resp, err := handler.merry.Proxy(request.Request, service, request.Method)
	return merry.Response(resp), err
}

// ------------------------------
// merry_fund
// ------------------------------

type RequestFund struct {
	To string `json:"to" binding:"required"`
}

// ------------------------------
// merry_version
// ------------------------------

type ResponseVersion struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// ------------------------------
// merry_relay
// ------------------------------

type RequestRelay struct {
	Method  string          `json:"method" binding:"required"`
	Request json.RawMessage `json:"request" binding:"required"`
}

type ResponseRelay interface{}
