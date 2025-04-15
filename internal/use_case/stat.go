package usecase

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"

	"backend/internal/infrastructure/repository"
	"backend/internal/models"
)

type StatService struct {
	statRepo *repository.StatRepository
}

func NewStatService(repo *repository.StatRepository) *StatService {
	return &StatService{statRepo: repo}
}

func (s *StatService) GetStatList(id uuid.UUID) (models.Stat, error) {
	return s.statRepo.GetStat(id)
}
func (s *StatService) ExportStatByPartner(partnerID uuid.UUID) ([]models.ActivationX, error) {
	return s.statRepo.ExportPromoByPartner(partnerID)
}
func (s *StatService) GetStatByPromo(partnerID uuid.UUID, promoID uuid.UUID) (models.StatPromo, error) {
	return s.statRepo.GetStatByPromo(partnerID, promoID)
}
func (s *StatService) ExportPartnerStat(partnerID uuid.UUID, requesterID uuid.UUID) (string, error) {
	if partnerID != requesterID {
		return "", errors.New("forbidden: access denied")
	}

	data, err := s.statRepo.ExportPromoByPartner(partnerID)
	if err != nil {
		return "", err
	}

	if len(data) == 0 {
		return "", errors.New("no data found for this partner")
	}

	filePath := fmt.Sprintf("partner_stats_%s.xlsx", partnerID.String())
	err = generateExcel(data, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func generateExcel(data []models.ActivationX, filePath string) error {
	f := excelize.NewFile()
	sheetName := "Statistics"
	f.SetSheetName("Sheet1", sheetName)
	headers := []string{"Client №.", "Promo Name", "Promo ID", "Buy", "Current", "Overall", "Created At", "Updated At"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}
	totalOverall := 0
	for i, record := range data {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), record.PromoName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), record.PromoID.String())
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), record.Buy)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), record.Current)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), record.Overall)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), record.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), record.UpdatedAt.Format("2006-01-02 15:04:05"))
		totalOverall += record.Overall
	}
	finalRow := len(data) + 2
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", finalRow), "TOTAL:")
	f.SetCellValue(sheetName, fmt.Sprintf("F%d", finalRow), totalOverall)
	statRow := finalRow + 1
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", statRow), "TOTAL Clients:")
	f.SetCellValue(sheetName, fmt.Sprintf("F%d", statRow), len(data))

	if err := f.SaveAs(filePath); err != nil {
		return err
	}

	return nil
}
