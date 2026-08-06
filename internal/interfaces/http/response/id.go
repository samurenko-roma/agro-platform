package response

import vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"

type IdResponse struct {
	ID vo.ID `json:"id"`
}

func Id(id vo.ID) IdResponse {
	return IdResponse{ID: id}
}
