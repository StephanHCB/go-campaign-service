package mysqldb

import (
	"context"
	"github.com/StephanHCB/go-campaign-service/internal/entity"
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/internal/repository/database/dbrepo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/zerolog/log"
	"time"
)

// use this for mocking

var fatalFunc = fatal
var warnFunc = warn
var infoFunc = info

type MysqlRepository struct {
	db *gorm.DB
}

func Create() dbrepo.Repository {
	return &MysqlRepository{}
}

func (r *MysqlRepository) Open() {
	db, err := gorm.Open("mysql", configuration.DatabaseMysqlConnectString())
	if err != nil {
		fatalFunc(err, "failed to open mysql connection")
	}

	// see https://making.pusher.com/production-ready-connection-pooling-in-go/
	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(50)
	db.DB().SetConnMaxLifetime(time.Minute * 10)

	r.db = db
}

func (r *MysqlRepository) Close() {
	err := r.db.Close()
	if err != nil {
		fatalFunc(err, "failed to close mysql connection")
	}
}

func (r *MysqlRepository) Migrate() {
	err := r.db.AutoMigrate(&entity.Campaign{}, &entity.Recipient{}).Error
	if err != nil {
		fatalFunc(err, "failed to migrate mysql db")
	}
}

func (r *MysqlRepository) AddCampaign(ctx context.Context, a *entity.Campaign) (uint, error) {
	err := r.db.Create(a).Error
	if err != nil {
		warnFunc(ctx, err, "mysql error during campaign insert")
	}
	return a.ID, err
}

func (r *MysqlRepository) UpdateCampaign(ctx context.Context, a *entity.Campaign) error {
	err := r.db.Save(a).Error
	if err != nil {
		warnFunc(ctx, err, "mysql error during campaign update")
	}
	return err
}

func (r *MysqlRepository) GetCampaignById(ctx context.Context, id uint) (*entity.Campaign, error) {
	var a entity.Campaign
	err := r.db.First(&a, id).Error
	if err != nil {
		infoFunc(ctx, err, "mysql error during campaign select - might be ok")
	}
	return &a, err
}

func fatal(err error, msg string) {
	log.Fatal().Err(err).Msg(msg + ": " + err.Error())
}

func warn(ctx context.Context, err error, msg string) {
	log.Ctx(ctx).Warn().Err(err).Msg(msg + ": " + err.Error())
}

func info(ctx context.Context, err error, msg string) {
	log.Ctx(ctx).Info().Err(err).Msg(msg + ": " + err.Error())
}
