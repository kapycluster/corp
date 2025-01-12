package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"kapycluster.com/corp/panel/dns"
	"kapycluster.com/corp/panel/kube"
	"kapycluster.com/corp/panel/views/dashboard"
)

func (h Handler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	u := h.MustGetUser(w, r)
	regions, err := h.db.GetUserRegions(r.Context(), u.UserID)
	if err != nil {
		h.Error(r.Context(), w, "failed to fetch regions", err)
		return
	}
	list, err := h.kc.ListControlPlanes(r.Context(), u.UserID, regions)
	if err != nil {
		h.Error(r.Context(), w, "failed to get control plane list", err)
		return
	}

	dashboard.ControlPlanes(u, list).Render(r.Context(), w)
}

func (h Handler) HandleCreateControlPlaneForm(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	region := r.FormValue("region")
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
		Region: region,
	}

	if err := h.kc.ValidateControlPlane(cp); err != nil {
		h.Error(r.Context(), w, "failed to validate control plane", err)
		return
	}

	if h.c.Server.LocalDev {
		h.log.Info("skipping dns record creation", "reason", "localdev")
	} else {
		h.log.Info("creating dns record", "record", h.controlPlaneAddress(namespace))
		if err := h.dns.CreateDNSRecord(r.Context(), dns.Record{
			Name: h.controlPlaneAddress(namespace),
			Type: "A",
			// TODO: this *has* to be parameterized
			Content: "65.109.40.187",
			TTL:     300,
			Proxied: false,
		}); err != nil {
			h.Error(r.Context(), w, "failed to create dns record", err)
			return
		}
	}

	h.log.Info("creating control plane!", slog.String("name", name), slog.String("namespace", namespace))
	if err := h.kc.CreateControlPlane(r.Context(), cp); err == nil {
		if err := h.db.CreateControlPlane(r.Context(), &cp); err != nil {
			h.log.Error("failed to store control plane in database", "error", err)
			h.Error(r.Context(), w, "failed to store control plane info", err)
			return
		}
		w.Header().Set("Hx-Redirect", "/controlplanes")
		w.WriteHeader(http.StatusOK)
	} else {
		h.log.Error(err.Error())
		h.Error(r.Context(), w, "failed to create control plane", err)
	}
}

func (h Handler) DownloadKubeconfigStub(w http.ResponseWriter, r *http.Request) {
	cpID := chi.URLParam(r, "id")
	if cpID == "" {
		h.Error(r.Context(), w, "missing control plane id", fmt.Errorf("no id provided"))
		return
	}

	w.Header().Set("Hx-Redirect", "/controlplanes/"+cpID+"/kubeconfig/download")
	w.WriteHeader(http.StatusOK)
}

func (h Handler) DownloadKubeconfig(w http.ResponseWriter, r *http.Request) {
	user := h.MustGetUser(w, r)
	if user.UserID == "" {
		return
	}

	cpID := chi.URLParam(r, "id")
	if cpID == "" {
		h.Error(r.Context(), w, "missing control plane id", fmt.Errorf("no id provided"))
		return
	}

	cpUser, err := h.db.GetControlPlaneUser(r.Context(), cpID)
	if err != nil {
		h.Error(r.Context(), w, "failed to get control plane user", err)
		return
	}

	if cpUser != user.UserID {
		h.Error(r.Context(), w, "not authorized to access control plane", fmt.Errorf("unauthorized"))
		return
	}

	cp, err := h.db.GetControlPlane(r.Context(), cpID)
	if err != nil {
		h.Error(r.Context(), w, "failed to get control plane region", err)
		return
	}

	kubeconfig, err := h.kc.GetKubeconfig(r.Context(), cpID, cp.Region)
	if err != nil {
		h.Error(r.Context(), w, "failed to get kubeconfig", err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(kubeconfig)))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=kubeconfig-%s.yaml", cpID))
	w.Write(kubeconfig)
}

func (h Handler) ShowCreateControlPlaneForm(w http.ResponseWriter, r *http.Request) {
	regions := h.kc.GetRegions()
	h.RenderOrRedirect(w, r, dashboard.CreateControlPlaneForm(regions), "/controlplanes")
}

func (h Handler) controlPlaneAddress(ns string) string {
	return fmt.Sprintf("%s.%s", ns, h.c.Server.ControlPlaneBaseURL)
}
