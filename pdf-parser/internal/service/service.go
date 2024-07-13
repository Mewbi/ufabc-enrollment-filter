package service

import (
	"bytes"
	"context"
	"time"

	"pdf-parser/internal/repository"
	"pdf-parser/types"

	"github.com/google/uuid"
)

type HttpClient interface {
	DownloadPDF(url string) (*bytes.Buffer, error)
}

type Parser interface {
	ParsePDF(content *bytes.Buffer) ([][]string, error)
	ConvertToCSV(rows [][]string) (*bytes.Buffer, error)
	SaveToCSVFile(content *bytes.Buffer, filename string) error
}

type Service struct {
	repository repository.Repository
	httpClient HttpClient
	parser     Parser
}

func New(repo repository.Repository, client HttpClient, parser Parser) *Service {
	return &Service{
		repository: repo,
		httpClient: client,
		parser:     parser,
	}
}

func (s *Service) CheckEnrollmentExist(ctx context.Context, url string) (*types.Enrollment, error) {
	return s.repository.GetEnrollmentByURL(ctx, url)
}

func (s *Service) GetEnrollment(ctx context.Context, id uuid.UUID) (*types.Enrollment, error) {
	return s.repository.GetEnrollmentByID(ctx, id.String())
}

func (s *Service) ListEnrollments(ctx context.Context) ([]types.Enrollment, error) {
	return s.repository.ListEnrollments(ctx)
}

func (s *Service) CreateEnrollment(ctx context.Context, infoPdf *types.InfoPDF) (*types.Enrollment, error) {
	content, err := s.httpClient.DownloadPDF(infoPdf.URL.String())
	if err != nil {
		return nil, err
	}

	rows, err := s.parser.ParsePDF(content)
	if err != nil {
		return nil, err
	}

	csvContent, err := s.parser.ConvertToCSV(rows)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	enrollment := &types.Enrollment{
		ID:        id,
		URL:       infoPdf.URL.String(),
		Name:      infoPdf.Name,
		Content:   csvContent.Bytes(),
		CreatedAt: time.Now(),
	}

	if err := s.repository.SaveEnrollment(ctx, enrollment); err != nil {
		return nil, err
	}

	return enrollment, nil
}
