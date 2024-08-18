package handlers

import (
	"net/http"

	"github.com/kapycluster/corpy/panel/views/dashboard"
)

type Dashboard struct{}

func (h Dashboard) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	dashboard.Show().Render(r.Context(), w)
}
