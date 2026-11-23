package user

import (
  "context"
  "errors"
  "fmt"
  "log"
  "strings"

  "github.com/michael870917/gocourse_domain/domain"
  "gorm.io/gorm"
)

type Repository interface {
  Create(ctx context.Context, user *domain.User) (*domain.User, error)
  GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
  // GetByID retrieves a pointer user by their ID and an error
  GetByID(ctx context.Context, id string) (*domain.User, error)
  Delete(ctx context.Context, id string) error
  Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error
  Count(ctx context.Context, filters Filters) (int64, error)
}

type UserRepository struct {
  logger *log.Logger
  db     *gorm.DB
}

func NewRepository(logger *log.Logger, db *gorm.DB) Repository {
  return &UserRepository{db: db, logger: logger}
}

func (repository *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
  if error := repository.db.WithContext(ctx).Create(user).Error; error != nil {
    repository.logger.Println("Error creating user in repository: ", error)
    return nil, error
  }

  repository.logger.Println("Created user in repository: ", user)
  return user, nil
}

func (repository *UserRepository) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
  var users []domain.User

  tx := repository.db.WithContext(ctx).Model(&users)
  tx = applyFilters(tx, filters)
  tx = tx.Offset(offset).Limit(limit)
  result := tx.Order("created_at desc").Find(&users)
  if result.Error != nil {
    repository.logger.Println("Error retrieving users from repository: ", result.Error)
    return nil, result.Error
  }

  return users, nil
}

func (repository *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
  user := domain.User{ID: id}
  if error := repository.db.WithContext(ctx).First(&user).Error; error != nil {
    repository.logger.Println("Error retrieving user from repository: ", error)

    if errors.Is(error, gorm.ErrRecordNotFound) {
      return nil, ErrNotFound{UserId: id}
    }

    return nil, error
  }

  repository.logger.Println("Retrieved user from repository: ", user)
  return &user, nil
}

func (repository *UserRepository) Delete(ctx context.Context, id string) error {
  user := domain.User{ID: id}
  result := repository.db.WithContext(ctx).Delete(&user)
  if result.Error != nil {
    repository.logger.Println("Error deleting user from repository: ", result.Error)
    return result.Error
  }

  if result.RowsAffected == 0 {
    repository.logger.Println("No user found to delete with ID: ", id)
    return ErrNotFound{UserId: id}
  }

  repository.logger.Println("Deleted user from repository with ID: ", id)
  return nil
}

func (repository *UserRepository) Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error {
  validateUser, _ := repository.GetByID(ctx, id)
  if validateUser == nil {
    repository.logger.Println("User not found in repository with ID: ", id)
    return fmt.Errorf("user with ID: %s not found", id)
  }

  values := make(map[string]interface{})
  if firstName != nil {
    values["first_name"] = *firstName
  }
  if lastName != nil {
    values["last_name"] = *lastName
  }
  if email != nil {
    values["email"] = *email
  }
  if phone != nil {
    values["phone"] = *phone
  }

  if len(values) == 0 {
    repository.logger.Println("No fields to update for user with ID: ", id)
    return nil // No fields to update
  }

  result := repository.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(values)
  if result.Error != nil {
    repository.logger.Println("Error updating user in repository: ", result.Error)
    return result.Error
  }

  if result.RowsAffected == 0 {
    repository.logger.Println("No user found to update with ID: ", id)
    return ErrNotFound{UserId: id}
  }

  repository.logger.Println("Updated user in repository with ID: ", id)
  return nil
}

func (repository *UserRepository) Count(ctx context.Context, filters Filters) (int64, error) {
  var count int64
  tx := repository.db.WithContext(ctx).Model(&domain.User{})
  tx = applyFilters(tx, filters)
  if error := tx.Count(&count).Error; error != nil {
    repository.logger.Println("Error counting users in repository: ", error)
    return 0, error
  }

  return count, nil
}

func applyFilters(db *gorm.DB, filters Filters) *gorm.DB {
  if filters.FirstName != "" {
    filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
    db = db.Where("lower(first_name) like ?", filters.FirstName)
  }
  if filters.LastName != "" {
    filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
    db = db.Where("lower(last_name) like ?", filters.LastName)
  }

  return db
}
