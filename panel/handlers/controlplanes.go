package handlers

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	kapyv1 "github.com/kapycluster/corpy/controller/api/v1"
	"github.com/kapycluster/corpy/panel/views/dashboard"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Dashboard struct {
	kc  KubeClient
	db  DBStore
	log *slog.Logger
}

func (h Dashboard) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	dashboard.ControlPlanes().Render(r.Context(), w)
}

func (h Dashboard) HandleCreateControlPlaneForm(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	namespace := uuid.New().String()

	cp := kapyv1.ControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: kapyv1.ControlPlaneSpec{
			Version: "v1.30",
			Server: kapyv1.KapyServer{
				Token: "dummy",
			},
		},
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

func (h Dashboard) ShowCreateControlPlaneForm(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("hx-request") != "" {
		dashboard.CreateControlPlaneForm().Render(r.Context(), w)
	} else {
		http.Redirect(w, r, "/controlplanes", http.StatusSeeOther)
	}
}
