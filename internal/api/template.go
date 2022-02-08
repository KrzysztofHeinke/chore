package api

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/registry"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/server/middleware"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/utils"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
)

type TemplatePureID struct {
	models.TemplatePure
	apimodels.ID
}

type MetaFolder struct {
	Folder string `json:"folder" query:"folder" example:"folderx"`
	apimodels.Meta
}

type ItemName struct {
	Item string `json:"item" example:"template1"`
	Name string `json:"name" example:"deepcore/template1"`
}

// @Summary List templates
// @Tags template
// @Description Get list of the templates, specify key query to get inner paths
// @Security ApiKeyAuth
// @Router /templates [get]
// @Param folder query string false "set the limit, default is empty"
// @Param limit query int false "set the limit, default is 20"
// @Param offset query int false "set the offset, default is 0"
// @Success 200 {object} apimodels.DataMeta{data=[]ItemName{}}
// @failure 400 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func listTemplates(c *fiber.Ctx) error {
	items := []ItemName{}

	meta := &MetaFolder{Meta: apimodels.Meta{Limit: apimodels.Limit}}

	if err := c.QueryParser(meta); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	// Table(reg.DB.Config.NamingStrategy.JoinTableName("folders"))

	result := reg.DB.WithContext(c.UserContext()).Model(&models.Folder{}).Select("item", "name").Limit(meta.Limit).Offset(meta.Offset).Where(
		"folder = ?", meta.Folder,
	).Find(&items)

	// check write error
	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	return c.Status(http.StatusOK).JSON(
		apimodels.DataMeta{
			Meta: meta,
			Data: apimodels.Data{Data: items},
		},
	)
}

// @Summary Get template
// @Tags template
// @Description Get one template with id
// @Security ApiKeyAuth
// @Router /template [get]
// @Param id query string false "get by id"
// @Param name query string false "get by name"
// @Success 200 {object} apimodels.Data{data=TemplatePureID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 404 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func getTemplate(c *fiber.Ctx) error {
	id := c.Query("id")
	name := c.Query("name")

	if id == "" && name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: apimodels.ErrRequiredIDName.Error(),
			},
		)
	}

	getData := TemplatePureID{}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.Template{})
	if id != "" {
		query = query.Where("id = ?", id)
	}

	if name != "" {
		query = query.Where("name = ?", name)
	}

	result := query.First(&getData)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: getData,
		},
	)
}

// @Summary New template
// @Tags template
// @Description Send and record new template
// @Security ApiKeyAuth
// @Router /template [post]
// @Param name query string true "name of file 'deepcore/template1'"
// @Param payload body string false "send template object"
// @Accept plain
// @Success 200 {object} apimodels.Data{data=apimodels.ID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func postTemplate(c *fiber.Ctx) error {
	template := new(models.Template)

	body := c.Body()
	template.Content = base64.StdEncoding.EncodeToString(body)

	name := c.Query("name")
	if name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: "name is required",
			},
		)
	}

	// trim slash
	template.Name = strings.Trim(name, "/")

	reg := registry.Reg().Get(c.Locals("registry").(string))

	var err error

	template.ID.ID, err = uuid.NewUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	result := reg.DB.WithContext(c.UserContext()).Create(template)

	// check write error
	var pErr *pgconn.PgError

	errors.As(result.Error, &pErr)

	if pErr != nil && pErr.Code == "23505" {
		return c.Status(http.StatusConflict).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	// create folder
	folderMap := utils.FolderFile(template.Name)

	// on conflict do nothing
	reg.DB.WithContext(c.UserContext()).Model(models.Folder{}).Clauses(
		clause.OnConflict{DoNothing: true},
	).Create(folderMap)

	// return recorded data's id
	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: apimodels.ID{ID: template.ID.ID},
		},
	)
}

// TODO: currently just changeable inside of the data.

// @Summary Replace template
// @Tags template
// @Description Replace with new data, id or name must exist in request
// @Security ApiKeyAuth
// @Router /template [patch]
// @Param name query string false "get by name"
// @Param payload body string false "send template object"
// @Accept plain
// @Success 200 {object} apimodels.Data{data=apimodels.ID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func patchTemplate(c *fiber.Ctx) error {
	name := c.Query("name")

	if name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: "name is required and cannot be empty",
			},
		)
	}

	body := c.Body()

	// fix parameter
	name = strings.Trim(name, "/")

	reg := registry.Reg().Get(c.Locals("registry").(string))

	data := models.Template{
		TemplatePure: models.TemplatePure{
			Name:    name,
			Content: base64.StdEncoding.EncodeToString(body),
		},
	}

	// save new value
	result := reg.DB.WithContext(c.UserContext()).Where("name = ?", name).Updates(&data)

	// check write error
	var pErr *pgconn.PgError

	errors.As(result.Error, &pErr)

	if pErr != nil && pErr.Code == "23505" {
		return c.Status(http.StatusConflict).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	// // update from folder table
	// if prevValues.Name != body["name"].(string) {
	// 	reg.DB.WithContext(c.UserContext()).Where("name = ?", prevValues.Name).Delete(&models.Folder{})

	// 	// create folder
	// 	folderMap := utils.FolderFile(body["name"].(string))

	// 	// on conflict do nothing
	// 	reg.DB.WithContext(c.UserContext()).Model(models.Folder{}).Clauses(
	// 		clause.OnConflict{DoNothing: true},
	// 	).Create(folderMap)
	// }

	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: fiber.Map{"id": data.ID},
		},
	)
}

// @Summary Delete template
// @Tags template
// @Description Delete with id, name
// @Security ApiKeyAuth
// @Router /template [delete]
// @Param id query string false "get by id"
// @Param name query string false "get by name"
// @Success 204 "No Content"
// @failure 400 {object} apimodels.Error{}
// @failure 404 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func deleteTemplate(c *fiber.Ctx) error {
	id := c.Query("id")
	name := c.Query("name")

	if id == "" && name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: apimodels.ErrRequiredIDName,
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext())
	if id != "" {
		query = query.Where("id = ?", id)
	}

	if name != "" {
		if name[len(name)-1] == '/' {
			query = query.Where("name LIKE ?", name+"%")
		} else {
			query = query.Where("name = ?", name)
		}
	}

	// delete directly in DB
	result := query.Unscoped().Delete(&models.Template{})

	if result.RowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(
			apimodels.Error{
				Error: "not found any releated data",
			},
		)
	}

	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	// delete from folder table
	query = reg.DB.WithContext(c.UserContext())
	if name[len(name)-1] == '/' {
		query = query.Where("name LIKE ?", name+"%")
	} else {
		query = query.Where("name = ?", name)
	}

	query.Delete(&models.Folder{})

	//nolint:wrapcheck // checking before
	return c.SendStatus(http.StatusNoContent)
}

func Template(router fiber.Router) {
	router.Get("/templates", middleware.JWTCheck(""), listTemplates)
	router.Get("/template", middleware.JWTCheck(""), getTemplate)
	router.Post("/template", middleware.JWTCheck(""), postTemplate)
	router.Patch("/template", middleware.JWTCheck(""), patchTemplate)
	router.Delete("/template", middleware.JWTCheck(""), deleteTemplate)
}
