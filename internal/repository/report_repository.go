package repository

import "gorm.io/gorm"

type ReportRepository interface {
	GetWasteSummary() ([]map[string]interface{}, error)
	GetPaymentSummary() ([]map[string]interface{}, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db}
}

// Rekap jumlah pickup berdasarkan tipe dan status
func (r *reportRepository) GetWasteSummary() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Table("waste_pickups").
		Select("type, status, count(*) as total_pickup").
		Group("type, status").
		Find(&results).Error
	return results, err
}

// Rekap total pendapatan berdasarkan status pembayaran
func (r *reportRepository) GetPaymentSummary() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Table("payments").
		Select("status, count(*) as total_transaksi, sum(amount) as total_pendapatan").
		Group("status").
		Find(&results).Error
	return results, err
}