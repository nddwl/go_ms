package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type FollowCount struct {
	ID          int64 `gorm:"primary_key"`
	UserID      int64
	FollowCount int
	FansCount   int
	CreateTime  time.Time
	UpdateTime  time.Time
}

type FollowCountModel struct {
	db *gorm.DB
}

func NewFollowCountModel(db *gorm.DB) *FollowCountModel {
	return &FollowCountModel{db}
}

func (m *FollowCountModel) Insert(ctx context.Context, data *FollowCount) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *FollowCountModel) FindOne(ctx context.Context, id int64) (*FollowCount, error) {
	var (
		f   FollowCount
		err error
	)
	err = m.db.WithContext(ctx).Where("id = ?", id).First(&f).Error
	return &f, err
}

func (m *FollowCountModel) Update(ctx context.Context, data *FollowCount) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *FollowCountModel) IncrFollowCount(ctx context.Context, userId int64) error {
	return m.db.WithContext(ctx).Exec("INSERT INTO follow_count (user_id, follow_count) VALUES (?, 1) ON DUPLICATE KEY UPDATE follow_count = follow_count + 1", userId).Error
}

func (m *FollowCountModel) DecrFollowCount(ctx context.Context, userId int64) error {
	return m.db.WithContext(ctx).Exec("INSERT INTO follow_count (user_id, follow_count) VALUES (?, 1) ON DUPLICATE KEY UPDATE follow_count = follow_count - 1", userId).Error
}

func (m *FollowCountModel) IncrFansCount(ctx context.Context, userId int64) error {
	return m.db.WithContext(ctx).Exec("INSERT INTO follow_count (user_id, fans_count) VALUES (?, 1) ON DUPLICATE KEY UPDATE fans_count = fans_count + 1", userId).Error
}

func (m *FollowCountModel) DecrFansCount(ctx context.Context, userId int64) error {
	return m.db.WithContext(ctx).Exec("INSERT INTO follow_count (user_id, fans_count) VALUES (?, 1) ON DUPLICATE KEY UPDATE fans_count = fans_count - 1", userId).Error
}

func (m *FollowCountModel) FindByUserIds(ctx context.Context, userIds []int64) ([]*FollowCount, error) {
	var (
		f   []*FollowCount
		err error
	)
	err = m.db.WithContext(ctx).Where("user_id in (?)", userIds).Find(&f).Error
	return f, err
}
