package handlers

import (
	"context"
	"fmt"
	"github.com/devpies/devpie-client-events/go/events"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mockPub "github.com/devpies/devpie-client-core/users/api/publishers/mocks"
	mockQuery "github.com/devpies/devpie-client-core/users/domain/mocks"
	"github.com/devpies/devpie-client-core/users/domain/projects"
	"github.com/devpies/devpie-client-core/users/domain/teams"
	mockAuth "github.com/devpies/devpie-client-core/users/platform/auth0/mocks"
	th "github.com/devpies/devpie-client-core/users/platform/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTeamMocks() *Team {
	mockRepo := th.Repo()
	mockNats := &events.Client{}
	mockAuth0 := &mockAuth.Auther{}
	mockTeamQueries := &mockQuery.TeamQuerier{}
	mockProjectQueries := &mockQuery.ProjectQuerier{}
	mockMembershipQueries := &mockQuery.MembershipQuerier{}
	mockUserQueries := &mockQuery.UserQuerier{}
	mockInviteQueries := &mockQuery.InviteQuerier{}
	mockPublishers := &mockPub.Publisher{}

	tq := TeamQueries{mockTeamQueries, mockProjectQueries, mockMembershipQueries, mockUserQueries, mockInviteQueries}

	return  &Team{
			repo:  mockRepo,
			nats: mockNats,
			auth0: mockAuth0,
			query: tq,
			publish: mockPublishers,
		}
}

func newTeam() teams.NewTeam {
	return teams.NewTeam{
		Name:      "TestTeam",
		ProjectID: "8695a94f-7e0a-4198-8c0a-d3e12727a5ba",
	}
}

func team() teams.Team {
	return teams.Team{
		ID: "39541c75-ca3e-4e2b-9728-54327772d001",
		Name:      "TestTeam",
		UserID:    "a4b54ec1-57f9-4c39-ab53-d936dbb6c177",
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

func teamJson(nt teams.NewTeam) string {
	return fmt.Sprintf(`{ "name": "%s", "projectId": "%s" }`,
		nt.Name, nt.ProjectID)
}

func TestTeams_Create_200(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nt := newTeam()
	tm := team()
	nm := newMembership(tm)
	m := membership(nm)

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", context.Background(), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Create", context.Background(), fake.repo, nt, uid, mock.AnythingOfType("time.Time")).Return(tm, nil)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Create", context.Background(), fake.repo, nm, mock.AnythingOfType("time.Time")).Return(m, nil)

	up := projects.UpdateProjectCopy{
		TeamID: &tm.ID,
	}

	fake.query.project.(*mockQuery.ProjectQuerier).On("Update", context.Background(), fake.repo, nt.ProjectID, up).Return(nil)
	fake.publish.(*mockPub.Publisher).On("MembershipCreatedForProject", fake.nats, m, nt.ProjectID, uid).Return(nil)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.Create(w, r)
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", strings.NewReader(teamJson(nt)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusCreated, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.query.membership.(*mockQuery.MembershipQuerier).AssertExpectations(t)
		fake.publish.(*mockPub.Publisher).AssertExpectations(t)
	})
}
