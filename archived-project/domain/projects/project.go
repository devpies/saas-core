package projects

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/devpies/devpie-client-core/projects/domain/memberships"
	"github.com/devpies/devpie-client-core/projects/platform/database"
)

var (
	ErrNotFound      = errors.New("project not found")
	ErrInvalidID     = errors.New("id provided was not a valid UUID")
	ErrNotAuthorized = errors.New("user does not have correct membership")
)

type ProjectQuerier interface {
	RetrieveTeamID(ctx context.Context, repo database.Storer, pid string) (string, error)

	Retrieve(ctx context.Context, repo database.Storer, pid, uid string) (Project, error)
	RetrieveShared(ctx context.Context, repo database.Storer, pid, uid string) (Project, error)

	List(ctx context.Context, repo database.Storer, np NewProject, uid string, now time.Time) (Project, error)

	Create(ctx context.Context, repo database.Storer, np NewProject, uid string, now time.Time) (Project, error)
	Delete(ctx context.Context, repo database.Storer, pid, uid string) error
}

func RetrieveTeamID(ctx context.Context, repo database.Storer, pid string) (string, error) {
	var p Project

	if _, err := uuid.Parse(pid); err != nil {
		return "", ErrInvalidID
	}

	stmt := repo.Select(
		"project_id",
		"name",
		"prefix",
		"description",
		"team_id",
		"user_id",
		"active",
		"public",
		"column_order",
		"updated_at",
		"created_at",
	).From("projects").Where(sq.Eq{"project_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return "", errors.Wrapf(err, "building query: %v", args)
	}

	row := repo.QueryRowxContext(ctx, q, pid)
	err = row.Scan(&p.ID, &p.Name, &p.Prefix, &p.Description, &p.TeamID, &p.UserID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrNotFound
		}
		return "", err
	}

	return p.TeamID, nil
}

func Retrieve(ctx context.Context, repo database.Storer, pid, uid string) (Project, error) {
	var p Project
	if _, err := uuid.Parse(pid); err != nil {
		return p, ErrInvalidID
	}

	stmt := repo.Select(
		"project_id",
		"name",
		"prefix",
		"description",
		"team_id",
		"user_id",
		"active",
		"public",
		"column_order",
		"updated_at",
		"created_at",
	).From("projects").Where(sq.Eq{"project_id": "?", "user_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return p, errors.Wrapf(err, "building query: %v", args)
	}

	row := repo.QueryRowxContext(ctx, q, pid, uid)
	err = row.Scan(&p.ID, &p.Name, &p.Prefix, &p.Description, &p.TeamID, &p.UserID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, ErrNotFound
		}
		return p, err
	}

	return p, nil
}

func RetrieveShared(ctx context.Context, repo database.Storer, pid, uid string) (Project, error) {
	var p Project

	tid, err := RetrieveTeamID(ctx, repo, pid)
	if err != nil {
		return p, err
	}

	if _, err = uuid.Parse(pid); err != nil {
		return p, ErrInvalidID
	}

	m, err := memberships.Retrieve(ctx, repo, uid, tid)
	if err != nil {
		return p, ErrNotAuthorized
	}

	stmt := repo.Select(
		"project_id",
		"name",
		"prefix",
		"description",
		"team_id",
		"user_id",
		"active",
		"public",
		"column_order",
		"updated_at",
		"created_at",
	).From("projects").Where(sq.Eq{"project_id": "?", "team_id": "?"})

	q, args, err := stmt.ToSql()

	if err != nil {
		return p, errors.Wrapf(err, "building query: %v", args)
	}

	row := repo.QueryRowxContext(ctx, q, pid, m.TeamID)
	err = row.Scan(&p.ID, &p.Name, &p.Prefix, &p.Description, &p.TeamID, &p.UserID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, ErrNotFound
		}
		return p, err
	}

	return p, nil
}

func List(ctx context.Context, repo database.Storer, uid string) ([]Project, error) {
	var p Project
	var ps = make([]Project, 0)

	q := `SELECT * FROM projects
		  WHERE team_id IN (SELECT team_id FROM memberships WHERE user_id = $1)
		  UNION 
		  SELECT * FROM projects 
		  WHERE user_id = $1
		  GROUP BY project_id`

	rows, err := repo.QueryxContext(ctx, q, uid)
	if err != nil {
		return nil, errors.Wrap(err, "selecting projects")
	}
	for rows.Next() {
		err = rows.Scan(&p.ID, &p.Name, &p.Prefix, &p.Description, &p.UserID, &p.TeamID, &p.Active, &p.Public, (*pq.StringArray)(&p.ColumnOrder), &p.UpdatedAt, &p.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scanning row into Struct")
		}
		ps = append(ps, p)
	}

	return ps, nil
}

func Create(ctx context.Context, repo database.Storer, np NewProject, uid string, now time.Time) (Project, error) {
	p := Project{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Prefix:      fmt.Sprintf("%s-", np.Name[:3]),
		Active:      true,
		UserID:      uid,
		TeamID:      np.TeamID,
		ColumnOrder: []string{"column-1", "column-2", "column-3", "column-4"},
		UpdatedAt:   now.UTC(),
		CreatedAt:   now.UTC(),
	}

	stmt := repo.Insert(
		"projects",
	).SetMap(map[string]interface{}{
		"project_id":   p.ID,
		"name":         p.Name,
		"prefix":       p.Prefix,
		"team_id":      p.TeamID,
		"description":  "",
		"user_id":      p.UserID,
		"column_order": pq.Array(p.ColumnOrder),
		"updated_at":   p.UpdatedAt,
		"created_at":   p.CreatedAt,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return p, errors.Wrapf(err, "inserting project: %v", p)
	}

	return p, nil
}

func Update(ctx context.Context, repo database.Storer, pid, uid string, update UpdateProject, now time.Time) (Project, error) {
	var p Project

	p, err := Retrieve(ctx, repo, pid, uid)
	if err != nil {
		p, err = RetrieveShared(ctx, repo, pid, uid)
		if err != nil {
			return p, err
		}
	}

	if update.Name != nil {
		p.Name = *update.Name
	}
	if update.Description != nil {
		p.Description = *update.Description
	}
	if update.Active != nil {
		p.Active = *update.Active
	}
	if update.Public != nil {
		p.Public = *update.Public
	}
	if update.ColumnOrder != nil {
		p.ColumnOrder = update.ColumnOrder
	}
	if update.TeamID != nil {
		p.TeamID = *update.TeamID
	}

	stmt := repo.Update(
		"projects",
	).SetMap(map[string]interface{}{
		"name":         p.Name,
		"description":  p.Description,
		"active":       p.Active,
		"public":       p.Public,
		"column_order": pq.Array(p.ColumnOrder),
		"team_id":      p.TeamID,
		"updated_at":   now.UTC(),
	}).Where(sq.Eq{"project_id": pid})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return p, errors.Wrap(err, "updating project")
	}

	return p, nil
}

func Delete(ctx context.Context, repo database.Storer, pid, uid string) error {
	if _, err := uuid.Parse(pid); err != nil {
		return ErrInvalidID
	}

	stmt := repo.Delete(
		"projects",
	).Where(sq.Eq{"project_id": pid, "user_id": uid})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting project %s", pid)
	}

	return nil
}
