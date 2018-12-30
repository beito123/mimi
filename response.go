package mimi

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var NotFoundResponse = &Base{
	Status: http.StatusNotFound,
	Error:  "Not found",
}

type Response interface {
	StatusCode() int
}

type Base struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func (base Base) StatusCode() int {
	return base.Status
}

type Render interface {
	Write(res Response) error
	NotFound() error
}

type BaseRender struct {
	Render

	Writer http.ResponseWriter
}

func NewJSONRender(w http.ResponseWriter) *JSONRender {
	return &JSONRender{
		BaseRender: &BaseRender{
			Writer: w,
		},
	}
}

type JSONRender struct {
	*BaseRender
}

func (render *JSONRender) Write(res Response) error {
	render.Writer.Header().Set("Content-Type", "application/json")
	render.Writer.WriteHeader(res.StatusCode())

	b, err := json.Marshal(res)
	if err != nil {
		return err
	}

	_, err = render.Writer.Write(b)

	return err
}

func (render *JSONRender) NotFound() error {
	return render.Write(NotFoundResponse)
}
