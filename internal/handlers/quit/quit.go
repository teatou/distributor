package quit

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/teatou/distributor/pkg/mylogger"
)

type Quitter interface {
	Quit(port int) error
}

type Request struct {
	Port int `json:"port"`
}

func New(q Quitter, logger mylogger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("failed to decode body")

			render.JSON(w, r, fmt.Errorf("failed to decode body"))

			return
		}

		if err := q.Quit(req.Port); err != nil {
			logger.Error(err.Error())

			render.JSON(w, r, fmt.Errorf(err.Error()))

			return
		}

		logger.Info("quitted successfully")
		render.JSON(w, r, "quitted successfully")
	}
}
