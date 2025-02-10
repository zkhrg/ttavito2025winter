package repository

import (
	"database/sql"
	"fmt"

	"gitverse-internship-zg/services/user-service/domain/entities"
	"gitverse-internship-zg/services/user-service/domain/interfaces"

	sq "github.com/Masterminds/squirrel"
)

type EntityRepo struct {
	DB      *sql.DB
	Builder sq.StatementBuilderType
}

func NewEntityRepo(db *sql.DB) interfaces.EntityRepository {
	return &EntityRepo{
		DB:      db,
		Builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *EntityRepo) GetByID(id string) (*entities.Entity, error) {
	query, args, err := r.Builder.
		Select("id", "first_name").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.DB.QueryRow(query, args...)
	entity := &entities.Entity{}
	if err := row.Scan(&entity.ID, &entity.FirstName); err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *EntityRepo) EditByID(newUser entities.EntityRequest) (*entities.Entity, error) {
	query, args, err := r.Builder.
		Update("users").
		Set("first_name", newUser.FirstName).
		Set("last_name", newUser.SecondName).
		Set("birth_date", newUser.BirthDate).
		Set("gender", newUser.Gender).
		Where(sq.Eq{"id": newUser.ID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	// Выполняем запрос на обновление
	result, err := r.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// Проверяем, сколько строк было обновлено
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	// Если строка не была обновлена, возвращаем ошибку
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no rows were updated")
	}

	query, args, err = r.Builder.
		Select("id", "first_name").
		From("users").
		Where(sq.Eq{"id": newUser.ID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.DB.QueryRow(query, args...)
	entity := &entities.Entity{}
	if err := row.Scan(&entity.ID, &entity.BirthDate); err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *EntityRepo) CreateUser(userData entities.EntityRequest) (*entities.Entity, error) {
	// Строим запрос на вставку нового пользователя
	query, args, err := r.Builder.
		Insert("users").
		Columns("first_name", "last_name", "birth_date", "gender").
		Values(userData.FirstName, userData.SecondName, userData.BirthDate, userData.Gender).
		Suffix("RETURNING id, first_name, last_name, birth_date, gender").
		ToSql()
	if err != nil {
		return nil, err
	}

	// Выполняем запрос на вставку
	row := r.DB.QueryRow(query, args...)

	// Создаем структуру для возврата вставленных данных
	entity := &entities.Entity{}
	if err := row.Scan(&entity.ID, &entity.FirstName, &entity.SecondName, &entity.BirthDate, &entity.Gender); err != nil {
		return nil, err
	}

	// Возвращаем вставленный объект
	return entity, nil
}
