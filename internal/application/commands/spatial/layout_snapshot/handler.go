package layoutsnapshot

import "github.com/samurenkoroma/agro-platform/internal/application/uow"

type Handler struct {
	uow uow.UnitOfWork
}

func NewHandler(uow uow.UnitOfWork) *Handler {
	return &Handler{uow: uow}
}
