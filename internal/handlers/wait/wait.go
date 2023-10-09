package wait

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/teatou/distributor/pkg/mylogger"
)

type Waiter interface {
	Wait() error
}

func New(waiter Waiter, logger mylogger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := waiter.Wait(); err != nil {
			logger.Error(err.Error())

			render.JSON(w, r, fmt.Errorf(err.Error()))

			return
		}

		logger.Info("waited successfully")
		render.JSON(w, r, "waited successfully")
	}
}
