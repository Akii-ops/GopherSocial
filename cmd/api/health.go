package main

import (
	"net/http"
)

// health godoc
//
//	@Summary		测试 应用程序 API 可行
//	@Description	返回json，描述环境信息和版本
//	@Tags			health
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	EnvelopeErr
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}

}
