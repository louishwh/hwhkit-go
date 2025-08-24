package database

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含常用字段
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Repository 通用仓储接口
type Repository[T any] interface {
	Create(entity *T) error
	GetByID(id uint) (*T, error)
	Update(entity *T) error
	Delete(id uint) error
	List(offset, limit int) ([]*T, error)
	Count() (int64, error)
	FindByCondition(condition map[string]interface{}) ([]*T, error)
	FindOneByCondition(condition map[string]interface{}) (*T, error)
}

// BaseRepository 基础仓储实现
type BaseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		db: db,
	}
}

// Create 创建实体
func (r *BaseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// GetByID 根据ID获取实体
func (r *BaseRepository[T]) GetByID(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update 更新实体
func (r *BaseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete 软删除实体
func (r *BaseRepository[T]) Delete(id uint) error {
	var entity T
	return r.db.Delete(&entity, id).Error
}

// HardDelete 硬删除实体
func (r *BaseRepository[T]) HardDelete(id uint) error {
	var entity T
	return r.db.Unscoped().Delete(&entity, id).Error
}

// List 分页获取实体列表
func (r *BaseRepository[T]) List(offset, limit int) ([]*T, error) {
	var entities []*T
	err := r.db.Offset(offset).Limit(limit).Find(&entities).Error
	return entities, err
}

// Count 获取实体总数
func (r *BaseRepository[T]) Count() (int64, error) {
	var count int64
	var entity T
	err := r.db.Model(&entity).Count(&count).Error
	return count, err
}

// FindByCondition 根据条件查询实体列表
func (r *BaseRepository[T]) FindByCondition(condition map[string]interface{}) ([]*T, error) {
	var entities []*T
	query := r.db
	for key, value := range condition {
		query = query.Where(key, value)
	}
	err := query.Find(&entities).Error
	return entities, err
}

// FindOneByCondition 根据条件查询单个实体
func (r *BaseRepository[T]) FindOneByCondition(condition map[string]interface{}) (*T, error) {
	var entity T
	query := r.db
	for key, value := range condition {
		query = query.Where(key, value)
	}
	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// GetDB 获取数据库实例（用于复杂查询）
func (r *BaseRepository[T]) GetDB() *gorm.DB {
	return r.db
}

// Paginate 分页查询
func (r *BaseRepository[T]) Paginate(page, pageSize int, condition map[string]interface{}) (PaginationResult[T], error) {
	var entities []*T
	var total int64
	
	query := r.db.Model(new(T))
	
	// 应用条件
	for key, value := range condition {
		query = query.Where(key, value)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return PaginationResult[T]{}, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return PaginationResult[T]{}, err
	}
	
	return PaginationResult[T]{
		Data:       entities,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + int64(pageSize) - 1) / int64(pageSize),
	}, nil
}

// PaginationResult 分页结果
type PaginationResult[T any] struct {
	Data       []*T  `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int64 `json:"total_pages"`
}

// BatchCreate 批量创建
func (r *BaseRepository[T]) BatchCreate(entities []*T, batchSize int) error {
	return r.db.CreateInBatches(entities, batchSize).Error
}

// BatchUpdate 批量更新
func (r *BaseRepository[T]) BatchUpdate(updates map[string]interface{}, condition map[string]interface{}) error {
	query := r.db.Model(new(T))
	for key, value := range condition {
		query = query.Where(key, value)
	}
	return query.Updates(updates).Error
}

// Exists 检查实体是否存在
func (r *BaseRepository[T]) Exists(condition map[string]interface{}) (bool, error) {
	var count int64
	query := r.db.Model(new(T))
	for key, value := range condition {
		query = query.Where(key, value)
	}
	err := query.Count(&count).Error
	return count > 0, err
}