package routes

import (
	"net/http"

	"github.com/mhmdfathan/funsies-be/handlers"
	"github.com/mhmdfathan/funsies-be/utils"
)

func UserRoutes(h *http.ServeMux) {
	h.Handle("POST /api/register", utils.WithMiddleware(handlers.Register, utils.CheckKey))
}