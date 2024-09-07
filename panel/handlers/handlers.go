package handlers

import (
	"net/http"

	"github.com/kapycluster/corpy/panel/views/dashboard"
)

type Dashboard struct{}

func (h Dashboard) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	dashboard.ControlPlanes().Render(r.Context(), w)
}

func (h Dashboard) CreateControlPlane(w http.ResponseWriter, r *http.Request) {
	dashboard.CreateControlPlane().Render(r.Context(), w)
}
