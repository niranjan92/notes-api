package notes

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/internal/errors"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"github.com/qiangxue/go-rest-api/pkg/pagination"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, rateLimiter routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Use(authHandler) // the following endpoints require a valid JWT
	r.Use(rateLimiter)
	r.Get("/notes/<id>", res.get)
	r.Get("/notes", res.query)

	r.Post("/notes", res.create)
	r.Put("/notes/<id>", res.update)
	r.Delete("/notes/<id>", res.delete)
	r.Post("/notes/<note_id>/share/<user_id>", res.share)

	r.Get("/search", res.search) // create separate controller later
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {

	note, err := r.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(note)
}

func (r resource) search(c *routing.Context) error {
	ctx := c.Request.Context()

	query := c.Request.URL.Query().Get("q")
	userId := c.Get("user_id").(string)

	notes, err := r.service.SearchNotes(ctx, userId, query)
	if err != nil {
		return err
	}

	return c.Write(notes)
}

func (r resource) query(c *routing.Context) error {
	ctx := c.Request.Context()

	userId, ok := c.Get("user_id").(string)
	if !ok {
		return errors.Unauthorized("user not found")
	}

	count, err := r.service.Count(ctx)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request, count)

	// fetch notes where user_id = userId
	notes, err := r.service.QueryByUser(ctx, userId)
	if err != nil {
		return err
	}

	// fetch shared notes
	sharedNotes, err := r.service.QuerySharedNotes(ctx, userId)
	if err != nil {
		return err
	}
	notes = append(notes, sharedNotes...)

	pages.Items = notes

	return c.Write(pages)
}

func (r resource) share(c *routing.Context) error {
	note_id := c.Param("note_id")
	user_id := c.Param("user_id")

	id := entity.GenerateID()
	input := ShareNoteRequest{
		ID:           id,
		NoteID:       note_id,
		SharedUserID: user_id,
	}
	note, err := r.service.ShareNote(c.Request.Context(), c.Param("note_id"), input)
	if err != nil {
		return err
	}

	return c.Write(note)
}

func (r resource) create(c *routing.Context) error {
	var input CreateNoteRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	userId := c.Get("user_id").(string)
	input.UserID = userId

	note, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(note, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateNoteRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	note, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(note)
}

func (r resource) delete(c *routing.Context) error {
	note, err := r.service.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(note)
}
