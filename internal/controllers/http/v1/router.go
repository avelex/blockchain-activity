package v1

import (
	"context"
	"log/slog"
	"net/http"

	httpHelper "github.com/avelex/blockchain-activity/internal/controllers/http"
	"github.com/avelex/blockchain-activity/internal/types"
)

type BlockchainUsecases interface {
	TopAddressesByActivity(ctx context.Context) (types.BlockchainActivity, error)
}

type handler struct {
	usecases BlockchainUsecases
}

func New(usecases BlockchainUsecases) *handler {
	return &handler{
		usecases: usecases,
	}
}

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/api/v1/top", h.showTopAddresses)
}

func (h *handler) showTopAddresses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	activity, err := h.usecases.TopAddressesByActivity(r.Context())
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	httpHelper.RenderJSON(w, activity)
}
