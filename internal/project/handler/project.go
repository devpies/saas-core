package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devpies/saas-core/internal/project/fail"
	"github.com/devpies/saas-core/internal/project/model"
	"github.com/devpies/saas-core/pkg/msg"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type projectService interface {
	List(ctx context.Context) ([]model.Project, error)
	Retrieve(ctx context.Context, projectID string) (model.Project, error)
	RetrieveShared(ctx context.Context, projectID string) (model.Project, error)
	Create(ctx context.Context, project model.NewProject, now time.Time) (model.Project, error)
	Update(ctx context.Context, projectID string, update model.UpdateProject, now time.Time) (model.Project, error)
	Delete(ctx context.Context, projectID string) error
}

// ProjectHandler handles the project requests.
type ProjectHandler struct {
	logger         *zap.Logger
	js             publisher
	projectService projectService
	columnService  columnService
	taskService    taskService
}

// NewProjectHandler returns a new project handler.
func NewProjectHandler(
	logger *zap.Logger,
	js publisher,
	projectService projectService,
	columnService columnService,
	taskService taskService,
) *ProjectHandler {
	return &ProjectHandler{
		logger:         logger,
		js:             js,
		projectService: projectService,
		columnService:  columnService,
		taskService:    taskService,
	}
}

// List handles project list requests.
func (ph *ProjectHandler) List(w http.ResponseWriter, r *http.Request) error {
	list, err := ph.projectService.List(r.Context())
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Retrieve handles project retrieval requests.
func (ph *ProjectHandler) Retrieve(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")

	opr, err := ph.projectService.Retrieve(r.Context(), pid)
	if err == nil {
		return web.Respond(r.Context(), w, opr, http.StatusOK)
	}

	spr, err := ph.projectService.RetrieveShared(r.Context(), pid)
	if err != nil {
		switch err {
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("error retrieving project %q: %w", pid, err)
		}
	}

	return web.Respond(r.Context(), w, spr, http.StatusOK)
}

// Create handles project create requests.
func (ph *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var (
		np  model.NewProject
		err error
	)

	values, ok := web.FromContext(r.Context())
	if !ok {
		return web.CtxErr()
	}

	if err = web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	project, err := ph.projectService.Create(r.Context(), np, time.Now())
	if err != nil {
		return err
	}

	titles := [4]string{"To Do", "In Progress", "Review", "Done"}

	ph.logger.Info(project.ID)
	for i, title := range titles {
		nt := model.NewColumn{
			ProjectID:  project.ID,
			Title:      title,
			ColumnName: fmt.Sprintf(`column-%d`, i+1),
		}
		_, err = ph.columnService.Create(r.Context(), nt, time.Now())
		if err != nil {
			return err
		}
	}

	e := msg.ProjectCreatedEvent{
		Data: msg.ProjectCreatedEventData{
			ProjectID:   project.ID,
			TenantID:    project.TenantID,
			Name:        project.Name,
			Prefix:      project.Prefix,
			Description: project.Description,
			TeamID:      project.TeamID,
			UserID:      project.UserID,
			Active:      project.Active,
			Public:      project.Public,
			ColumnOrder: project.ColumnOrder,
			UpdatedAt:   project.UpdatedAt.String(),
			CreatedAt:   project.CreatedAt.String(),
		},
		Type: msg.TypeProjectCreated,
		Metadata: msg.Metadata{
			TenantID: values.TenantID,
			UserID:   values.UserID,
			TraceID:  values.TraceID,
		},
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	ph.js.Publish(msg.SubjectProjectCreated, bytes)

	return web.Respond(r.Context(), w, project, http.StatusCreated)
}

// Update handles project update requests.
func (ph *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) error {
	var update model.UpdateProject

	values, ok := web.FromContext(r.Context())
	if !ok {
		return web.CtxErr()
	}

	pid := chi.URLParam(r, "pid")

	if err := web.Decode(r, &update); err != nil {
		return fmt.Errorf("error decoding project update: %w", err)
	}

	up, err := ph.projectService.Update(r.Context(), pid, update, time.Now())
	if err != nil {
		switch err {
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("error updating project %q: %w", pid, err)
		}
	}

	e := msg.ProjectUpdatedEvent{
		Type: msg.TypeProjectUpdated,
		Data: msg.ProjectUpdatedEventData{
			Name:        &up.Name,
			Description: &up.Description,
			Active:      &up.Active,
			Public:      &up.Public,
			TeamID:      &up.TeamID,
			ProjectID:   up.ID,
			ColumnOrder: up.ColumnOrder,
			UpdatedAt:   up.UpdatedAt.String(),
		},
		Metadata: msg.Metadata{
			TenantID: values.TenantID,
			UserID:   values.UserID,
			TraceID:  values.TraceID,
		},
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	ph.js.Publish(msg.SubjectProjectUpdated, bytes)

	return web.Respond(r.Context(), w, up, http.StatusOK)
}

// Delete handles project delete requests.
func (ph *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	var err error
	pid := chi.URLParam(r, "pid")

	values, ok := web.FromContext(r.Context())
	if !ok {
		return web.CtxErr()
	}

	if _, err = ph.projectService.Retrieve(r.Context(), pid); err != nil {
		_, err = ph.projectService.RetrieveShared(r.Context(), pid)
		if err == nil {
			return web.NewRequestError(err, http.StatusUnauthorized)
		}
	}
	if err = ph.taskService.DeleteAll(r.Context(), pid); err != nil {
		return err
	}
	if err = ph.columnService.DeleteAll(r.Context(), pid); err != nil {
		return err
	}
	if err = ph.projectService.Delete(r.Context(), pid); err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("error deleting project %q: %w", pid, err)
		}
	}

	e := msg.ProjectDeletedEvent{
		Type: msg.TypeProjectDeleted,
		Metadata: msg.Metadata{
			TenantID: values.TenantID,
			TraceID:  values.TraceID,
			UserID:   values.UserID,
		},
		Data: msg.ProjectDeletedEventData{
			ProjectID: pid,
		},
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	ph.js.Publish(msg.SubjectProjectDeleted, bytes)

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}
