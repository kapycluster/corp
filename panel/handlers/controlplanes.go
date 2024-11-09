package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/kapycluster/corpy/panel/dns"
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
		Network: kube.Network{
			LoadBalancerAddress: h.controlPlaneAddress(namespace),
		},
	}

	if err := kube.ValidateControlPlane(cp); err != nil {
		views.Error(err.Error()).Render(r.Context(), w)
		return
	}

	h.log.Info("creating dns record", "record", h.controlPlaneAddress(namespace))
	if err := h.dns.CreateDNSRecord(r.Context(), dns.Record{
		Name: h.controlPlaneAddress(namespace),
		Type: "A",
		// TODO: this *has* to be parameterized
		Content: "65.109.40.187",
		TTL:     300,
		Proxied: false,
	}); err != nil {
		h.log.Error(err.Error())
		views.Error(err.Error()).Render(r.Context(), w)
		return
	}

	h.log.Info("creating control plane!", slog.String("name", name), slog.String("namespace", namespace))
	if err := h.kc.CreateControlPlane(r.Context(), cp); err == nil {
		w.Header().Set("Hx-Redirect", "/controlplanes")
		w.WriteHeader(http.StatusOK)
	} else {
		h.log.Error(err.Error())
		views.Error(err.Error()).Render(r.Context(), w)
	}

}

func (h Handler) ShowCreateControlPlaneForm(w http.ResponseWriter, r *http.Request) {
	h.RenderOrRedirect(w, r, dashboard.CreateControlPlaneForm(), "/controlplanes")
}

func (h Handler) controlPlaneAddress(ns string) string {
	return fmt.Sprintf("%s.%s", ns, h.c.Server.ControlPlaneBaseURL)
}
