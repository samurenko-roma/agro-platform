package swaggerdocs

import (
	"github.com/samurenkoroma/agro-platform/internal/application/commands/account/organization"
	"github.com/samurenkoroma/agro-platform/internal/application/queries/account/dto"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
)

var (
	_ = organization.CreateOrganizationCmd{}
	_ = dto.UserOrganizationInfo{}
	_ = organization.SwitchOrganizationResult{}
	_ = response.CommandResponse[any]{}
)

// docCreateOrganization godoc
// @Summary Создать организацию (ферму/хозяйство), пользователь становится владельцем
// @Description Можно вызвать напрямую этим путём (тело = data), либо через POST /api/commands с конвертом {"command": "account.create_organization", "data": {...}}
// @Tags account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body organization.CreateOrganizationCmd true "Название организации"
// @Success 200 {object} response.CommandResponse[dto.UserOrganizationInfo]
// @Router /api/commands/account.create_organization [post]
func docCreateOrganization() {}

// docSwitchOrganization godoc
// @Summary Переключиться на другую организацию (перевыпускает JWT с новым organization_id)
// @Description Можно вызвать напрямую этим путём (тело = data), либо через POST /api/commands с конвертом {"command": "account.switch_organization", "data": {...}}
// @Tags account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body organization.SwitchOrganizationCmd true "ID организации, в которой есть членство"
// @Success 200 {object} response.CommandResponse[organization.SwitchOrganizationResult]
// @Failure 403 {object} response.CommandResponse[any] "нет членства в организации"
// @Router /api/commands/account.switch_organization [post]
func docSwitchOrganization() {}

// docMe godoc
// @Summary Текущий пользователь + список организаций, в которых он состоит
// @Description Можно вызвать напрямую этим путём (тело = data), либо через POST /api/queries с конвертом {"query": "account.me", "data": {...}}
// @Tags account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.CommandResponse[any]
// @Router /api/queries/account.me [post]
func docMe() {}
