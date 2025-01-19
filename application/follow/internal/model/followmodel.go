package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Follow struct {
	ID             int64 `gorm:"primary_key"`
	UserID         int64
	FollowedUserID int64
	FollowStatus   int
	CreateTime     time.Time
	UpdateTime     time.Time
}

func (Follow) TableName() string {
	return "follow"
}

type FollowModel struct {
	db *gorm.DB
}

func NewFollowModel(db *gorm.DB) *FollowModel {
	return &FollowModel{db}
}

func (m *FollowModel) Insert(ctx context.Context, data *Follow) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *FollowModel) FindOne(ctx context.Context, id int64) (*Follow, error) {
	var (
		f   Follow
		err error
	)
	err = m.db.WithContext(ctx).Where("id = ?", id).First(&f).Error
	return &f, err
}

func (m *FollowModel) Update(ctx context.Context, data *Follow) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *FollowModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&Follow{}).Where("id = ?", id).Updates(values).Error
}

func (m *FollowModel) FindByUserIdAndFollowedUserId(ctx context.Context, userId int64, followedUserId int64) (*Follow, error) {
	var (
		f   Follow
		err error
	)
	err = m.db.WithContext(ctx).Where("user_id = ?", userId).Where("followed_user_id", followedUserId).First(&f).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &f, err
}

func (m *FollowModel) FindByUserId(ctx context.Context, userId int64, limit int) ([]*Follow, error) {
	var (
		f   []*Follow
		err error
	)
	err = m.db.WithContext(ctx).Where("user_id = ?", userId).Where("follow_status = ?", 1).Limit(limit).Order("id desc").Find(&f).Error
	return f, err
}

func (m *FollowModel) FindByFollowedUserIds(ctx context.Context, userId int64, followedUserIds []int64) ([]*Follow, error) {
	var (
		f   []*Follow
		err error
	)
	err = m.db.WithContext(ctx).Where("followed_user_id = ?", userId).Where("follow_user_id in (?)", followedUserIds).Find(&f).Error
	return f, err
}

func (m *FollowModel) FindByFollowedUserId(ctx context.Context, userId int64, limit int) ([]*Follow, error) {
	var (
		f   []*Follow
		err error
	)
	err = m.db.WithContext(ctx).Where("followed_user_id = ?", userId).Where("follow_status = ?", 1).Limit(limit).Order("id desc").Find(&f).Error
	return f, err

}
