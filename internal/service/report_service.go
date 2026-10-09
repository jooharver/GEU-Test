package service

import "github.com/jooharver/geu-test/internal/repository"

type ReportService interface {
	GetWasteSummary() ([]map[string]interface{}, error)
	GetPaymentSummary() ([]map[string]interface{}, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo}
}

func (s *reportService) GetWasteSummary() ([]map[string]interface{}, error) {
	return s.repo.GetWasteSummary()
}

func (s *reportService) GetPaymentSummary() ([]map[string]interface{}, error) {
	return s.repo.GetPaymentSummary()
}