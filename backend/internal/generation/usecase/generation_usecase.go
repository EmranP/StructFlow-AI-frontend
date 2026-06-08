package usecase

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/generation/dto"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/generation/service"
	genservice "github.com/EmranP/Design-Struct-Project-AI/backend/internal/generation/service"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/generation/status"
	zipservice "github.com/EmranP/Design-Struct-Project-AI/backend/internal/generation/zip"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/repository"
	customerrors "github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/errors"
)

type generationUseCase struct {
	projectRepo repository.ProjectRepository
	generator   *genservice.Generator

	generationRepo     repository.GenerationRepository
	generationTempRepo repository.GeneratedTemplateRepository
	zipService         *zipservice.Service
}

func New(
	projectRepo repository.ProjectRepository,
	generationRepo repository.GenerationRepository,
	generationTempRepo repository.GeneratedTemplateRepository,
	generator *service.Generator,
	zipService *zipservice.Service,
) GenerationUseCase {
	return &generationUseCase{
		projectRepo:        projectRepo,
		generationRepo:     generationRepo,
		generationTempRepo: generationTempRepo,
		generator:          generator,
		zipService:         zipService,
	}
}

func (u *generationUseCase) Add(
	ctx context.Context,
	projectID uuid.UUID,
) (*domain.Generation, error) {

	project, err := u.projectRepo.
		GetByID(
			ctx,
			projectID,
		)

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil,
			customerrors.ErrProjectNotFound
	}

	generation := &domain.Generation{
		ID: uuid.New(),

		ProjectID: projectID,

		Status: status.Pending,
	}

	errNewGen := u.generationRepo.Create(ctx, generation)
	if errNewGen != nil {

		return nil, errNewGen
	}

	go u.generator.Process(
		ctx,
		generation.ID,
		project.ID,
	)

	return generation, nil
}

func (u *generationUseCase) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Generation, error) {

	generation, err := u.generationRepo.
		GetByID(
			ctx,
			id,
		)

	if err != nil {
		return nil, err
	}

	if generation == nil {
		return nil,
			customerrors.ErrGenerationNotFound
	}

	return generation, nil
}

func (u *generationUseCase) FindByProjectID(
	ctx context.Context,
	projectID uuid.UUID,
	limit,
	page int,
) ([]dto.GenerationAllResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	total, err := u.generationRepo.GetCount(ctx, projectID)
	if err != nil {

		return nil, 0, err
	}

	if total == 0 {
		return []dto.GenerationAllResponse{}, 0, nil
	}

	gens, err := u.generationRepo.GetAllByProjectID(
		ctx,
		projectID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	gensDto := make([]dto.GenerationAllResponse, len(gens))
	for i, g := range gens {
		var errMsg string
		if g.ErrorMessage != nil {
			errMsg = *g.ErrorMessage
		}

		gensDto[i] = dto.GenerationAllResponse{
			ID:           g.ID.String(),
			ProjectID:    g.ProjectID.String(),
			Status:       g.Status,
			ErrorMessage: errMsg,
			CreatedAt:    g.CreatedAt,
			UpdatedAt:    g.UpdatedAt,
		}
	}

	return gensDto, total, nil
}

func (u *generationUseCase) FindTemplates(
	ctx context.Context,
	genId uuid.UUID,
) ([]dto.GenerationTempResponse, error) {
	genTemps, err := u.generationTempRepo.GetByGenerationID(ctx, genId)
	if err != nil {
		return nil, err
	}

	genTempsDto := make([]dto.GenerationTempResponse, len(genTemps))
	for i, gTemp := range genTemps {
		genTempsDto[i] = dto.GenerationTempResponse{
			ID:        gTemp.ID.String(),
			GenID:     gTemp.GenerationID.String(),
			Type:      gTemp.Type,
			Content:   json.RawMessage(gTemp.Content),
			CreatedAt: gTemp.CreatedAt,
		}
	}

	return genTempsDto, nil
}

func (u *generationUseCase) Download(
	ctx context.Context,
	generationID uuid.UUID,
) ([]byte, error) {
	templates, err := u.generationTempRepo.
		GetByGenerationID(
			ctx,
			generationID,
		)
	if err != nil {
		return nil, err
	}

	if len(templates) == 0 {
		return nil,
			customerrors.ErrGenerationNotFound
	}

	return u.zipService.Generate(
		templates[0].Content,
	)
}
