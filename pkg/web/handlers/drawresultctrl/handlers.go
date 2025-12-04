package drawresultctrl

import (
	"github.com/go-fuego/fuego"
)

type DrawResultHandlers struct {
	// TODO: Add service dependencies
}

func NewDrawResultHandlers() *DrawResultHandlers {
	return &DrawResultHandlers{}
}

// GetDrawResult handles GET /secret-friends/{id}/draw-result
func (h *DrawResultHandlers) GetDrawResult(
	c fuego.ContextNoBody,
) (*DrawResultResponse, error) {
	// TODO: Extract secretFriendId from path
	// TODO: Implement service call
	return nil, nil
}
