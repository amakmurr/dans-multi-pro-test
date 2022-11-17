package internal

import (
	"database/sql"
	"errors"
	"github.com/amakmurr/dans-multi-pro-test/internal/core"
	"github.com/amakmurr/dans-multi-pro-test/pkg/jwt"
	"net/http"
	"strconv"

	"github.com/amakmurr/dans-multi-pro-test/pkg/dans"
	v1 "github.com/amakmurr/dans-multi-pro-test/pkg/openapi"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type Endpoint struct {
	userRepository core.UserRepository
	jwtClient      jwt.Client
	dansClient     dans.ClientInterface
}

func NewEndpoint(userRepository core.UserRepository, jwtClient jwt.Client, dansClient dans.ClientInterface) *Endpoint {
	return &Endpoint{
		userRepository: userRepository,
		jwtClient:      jwtClient,
		dansClient:     dansClient,
	}
}

func (e *Endpoint) login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := e.userRepository.GetByUsername(r.Context(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handleError(w, r, ErrUnauthorized)
			return
		}
		handleError(w, r, err)
		return
	}
	if !user.CheckPassword(password) {
		handleError(w, r, ErrUnauthorized)
		return
	}
	accessToken, err := e.jwtClient.GenerateToken(jwt.UserClaims{
		UserID:   user.ID,
		Username: user.Username,
	})
	if err != nil {
		handleError(w, r, err)
		return
	}

	resp := &v1.LoginResponse{
		Data: v1.AccessToken{
			AccessToken: accessToken,
		},
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

func (e *Endpoint) getJobList(w http.ResponseWriter, r *http.Request) {
	params := &dans.GetJobListParams{}
	p := r.URL.Query()
	if p.Has("description") {
		params.Description = lo.ToPtr(p.Get("description"))
	}
	if p.Has("location") {
		params.Location = lo.ToPtr(p.Get("location"))
	}
	if p.Has("full_time") {
		params.FullTime = lo.ToPtr(true)
	}
	page := 1
	if p.Has("page") {
		var err error
		page, err = strconv.Atoi(p.Get("page"))
		if err != nil {
			handleError(w, r, ErrValidation)
			return
		}
	}
	params.Page = &page

	jobs, err := e.dansClient.GetJobList(r.Context(), params)
	if err != nil {
		handleError(w, r, err)
		return
	}

	resp := &v1.JobsResponse{}
	resp.Data = constructJobsToResp(jobs)
	resp.Meta.Page = page
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

func (e *Endpoint) getJobDetail(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		handleError(w, r, ErrValidation)
		return
	}

	job, err := e.dansClient.GetJobDetail(r.Context(), id)
	if err != nil {
		handleError(w, r, err)
		return
	}
	if job.Title == "" {
		handleError(w, r, ErrDataNotFound)
		return
	}

	resp := &v1.JobDetailResponse{}
	resp.Data = constructJobToResp(*job)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

func constructJobsToResp(jobs []*dans.Job) v1.Jobs {
	var j v1.Jobs
	for _, job := range jobs {
		if job != nil {
			j = append(j, constructJobToResp(*job))
		}
	}
	return j
}

func constructJobToResp(job dans.Job) v1.Job {
	return v1.Job{
		Company:     job.Company,
		CompanyLogo: job.CompanyLogo,
		CompanyUrl:  job.CompanyUrl,
		CreatedAt:   job.CreatedAt,
		Description: job.Description,
		HowToApply:  job.HowToApply,
		Id:          job.Id,
		Location:    job.Location,
		Title:       job.Title,
		Type:        job.Type,
		Url:         job.Url,
	}
}
