package handlers

import (
	"net/http"

	"github.com/google/uuid"

	kapyv1 "github.com/kapycluster/corpy/controller/api/v1"
	"github.com/kapycluster/corpy/panel/views/dashboard"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Dashboard struct {
	kc KubeClient
	db DBStore
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
			Server: kapyv1.KapyServer{
				Token: "dummy",
			},
		},
	}

	if err := h.kc.CreateControlPlane(r.Context(), cp); err != nil {
		// TODO: make this a templ error component
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h Dashboard) ShowCreateControlPlaneForm(w http.ResponseWriter, r *http.Request) {
	dashboard.CreateControlPlane().Render(r.Context(), w)
}
