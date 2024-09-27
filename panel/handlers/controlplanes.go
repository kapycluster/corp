package handlers

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/kapycluster/corpy/panel/kube"
	"github.com/kapycluster/corpy/panel/views"
	"github.com/kapycluster/corpy/panel/views/dashboard"
)

func (h Handler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	u := h.MustGetUser(w, r)
	list, err := h.kc.ListControlPlanes(r.Context(), u.UserID)
	if err != nil {
		h.log.Error(err.Error())
		views.Error("failed to get control plane list").Render(r.Context(), w)
		return
	}

	dashboard.ControlPlanes(u, list).Render(r.Context(), w)
}

func (h Handler) HandleCreateControlPlaneForm(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	namespace := uuid.New().String()

	user := h.MustGetUser(w, r)
	if user.UserID == "" {
		return
	}

	cp := kube.ControlPlane{
		Name:   name,
		ID:     namespace,
		UserID: user.UserID,
	}

	h.log.Info("creating control plane!", slog.String("name", name), slog.String("namespace", namespace))
	if err := h.kc.CreateControlPlane(r.Context(), cp); err == nil {
		w.Header().Set("Hx-Redirect", "/controlplanes")
		w.WriteHeader(http.StatusOK)
	} else {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (h Handler) ShowCreateControlPlaneForm(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("hx-request") != "" {
		dashboard.CreateControlPlaneForm().Render(r.Context(), w)
	} else {
		http.Redirect(w, r, "/controlplanes", http.StatusSeeOther)
	}
}
