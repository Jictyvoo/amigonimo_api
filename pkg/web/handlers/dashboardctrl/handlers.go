package dashboardctrl

import (
	"github.com/go-fuego/fuego"
)

type DashboardHandlers struct {
	// TODO: Add service dependencies
}

func NewDashboardHandlers() *DashboardHandlers {
	return &DashboardHandlers{}
}

// GetDashboard handles GET /dashboard.
func (h *DashboardHandlers) GetDashboard(
	c fuego.ContextNoBody,
) (*DashboardResponse, error) {
	// TODO: Extract userId from JWT token
	// TODO: Implement service call
	return nil, nil
}
