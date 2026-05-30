package repo

import (
	"fmt"
	"time"

	"github.com/Mgeorg1/go_randomaliens/gen/pb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func CreateDb(addr string, port int, user, password, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		addr,
		port,
		user,
		password,
		dbName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (r *Repo) Close() error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func NewRepo(addr string, port int, user, password, dbName string) (*Repo, error) {
	db, err := CreateDb(addr, port, user, password, dbName)
	if err != nil {
		return nil, err
	}
	repo := &Repo{db: db}
	if err := repo.AutoMigrate(); err != nil {
		return nil, err
	}
	return repo, nil
}

type FrequencyEvent struct {
	SessionID string `gorm:"primaryKey, type:uuid"`
	Frequency float64
	Timestamp time.Time `gorm:"type:timestamptz"`
}

func (event *FrequencyEvent) FromPbEvent(pbEvent *pb.FrequencyEvent) {
	event.SessionID = pbEvent.SessionId
	event.Frequency = pbEvent.Frequency
	event.Timestamp = pbEvent.Timestamp.AsTime()
}

func (r *Repo) SaveFrequencyEvent(event *FrequencyEvent) error {
	return r.db.Create(event).Error
}

func (r *Repo) AutoMigrate() error {
	return r.db.AutoMigrate(&FrequencyEvent{})
}
