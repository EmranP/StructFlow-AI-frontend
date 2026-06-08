package usecase

import (
	"context"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/project/dto"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/repository"
	"github.com/google/uuid"
)

type projectUseCase struct {
	projectRepo repository.ProjectRepository
}

type countResult struct {
	totalCount int64
	err        error
}

func New(projectRepo repository.ProjectRepository) ProjectUseCase {
	return &projectUseCase{
		projectRepo: projectRepo,
	}
}

func (u *projectUseCase) Add(
	ctx context.Context,
	userID uuid.UUID,
	p *dto.ProjectRequest,
) (uuid.UUID, error) {

	project := &domain.Project{
		ID: uuid.New(),

		UserID: userID,

		Title: p.Title,

		ProjectType: p.ProjectType,

		Stack: p.Stack,

		Architecture: p.Architecture,

		Features: p.Features,

		AdditionalInfo: p.AdditionalInfo,
	}

	newProject, err := u.projectRepo.Create(
		ctx,
		project,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return newProject, nil
}

func (u *projectUseCase) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*dto.ProjectResponse, error) {
	project, err := u.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	projectDto := &dto.ProjectResponse{
		ID:             project.ID,
		Title:          project.Title,
		ProjectType:    project.ProjectType,
		Stack:          project.Stack,
		Architecture:   project.Architecture,
		Features:       project.Features,
		AdditionalInfo: project.AdditionalInfo,
		CreatedAt:      project.CreatedAt,
		UpdatedAt:      project.UpdatedAt,
	}

	return projectDto, nil
}

func (u *projectUseCase) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
	page,
	limit int,
) ([]dto.ProjectResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	total, err := u.projectRepo.GetCount(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []dto.ProjectResponse{}, 0, nil
	}

	projects, err := u.projectRepo.GetAllByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	projectsDto := make([]dto.ProjectResponse, len(projects))
	for i, p := range projects {
		projectsDto[i] = dto.ProjectResponse{
			ID:             p.ID,
			Title:          p.Title,
			ProjectType:    p.ProjectType,
			Stack:          p.Stack,
			Architecture:   p.Architecture,
			Features:       p.Features,
			AdditionalInfo: p.AdditionalInfo,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
		}
	}

	return projectsDto, total, nil
}

func (u *projectUseCase) Edit(
	ctx context.Context,
	id,
	userID uuid.UUID,
	data *dto.ProjectRequest,
) error {
	updateProject := &domain.Project{
		Title:          data.Title,
		ProjectType:    data.ProjectType,
		Stack:          data.Stack,
		Architecture:   data.Architecture,
		Features:       data.Features,
		AdditionalInfo: data.AdditionalInfo,
	}

	errUpdate := u.projectRepo.Update(ctx, id, userID, updateProject)
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (u *projectUseCase) Remove(
	ctx context.Context,
	id,
	userID uuid.UUID,
) error {
	errDel := u.projectRepo.Delete(ctx, id, userID)
	if errDel != nil {
		return errDel
	}

	return nil
}

func (u *projectUseCase) RemoveAll(
	ctx context.Context,
	userID uuid.UUID,
) error {
	errDel := u.projectRepo.DeleteAll(ctx, userID)
	if errDel != nil {
		return errDel
	}

	return nil
}
