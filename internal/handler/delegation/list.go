package delegation

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	delegationlist "kiln-exercice/internal/usecase/delegation/list"
	"kiln-exercice/pkg/http/api"
)

type DelegationUseCase interface {
	ListDelegations(ctx context.Context, input delegationlist.Input) (delegationlist.Output, error)
}

type DelegationListHandler struct {
	useCase DelegationUseCase
}

func NewDelegationHandler(useCase DelegationUseCase) *DelegationListHandler {
	return &DelegationListHandler{
		useCase: useCase,
	}
}

func (h *DelegationListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.Handle(h.Handle)(w, r)
}

func (h *DelegationListHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var (
		ctx   = r.Context()
		input delegationlist.Input
		err   error
	)

	yearQuery := r.URL.Query().Get("year")
	if yearQuery != "" {
		input.Year, err = strconv.Atoi(yearQuery)
		if err != nil {
			return api.BadRequestError(fmt.Sprintf("invalid year format: %s", yearQuery), err)
		}
	}

	pagination, err := api.PaginationFromRequest(r)
	if err != nil {
		return api.BadRequestError(fmt.Sprintf("Pagination error: %s.", err.Error()))
	}

	input.Pagination = pagination

	delegations, err := h.useCase.ListDelegations(ctx, input)
	if err != nil {
		return err
	}

	return api.JSONResponse(w, http.StatusOK, delegations)
}
