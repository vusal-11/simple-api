package repository

import (
	"database/sql"
	
	"fmt"
	"simple_api/internal/service"
)

type UserRepository struct {
	db *sql.DB
}


func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db:db}
}

func (r *UserRepository) Create(name,email,password string) (service.User,error){
	var user service.User
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRow(query, name, email, password).Scan(&user.ID)
	if err != nil {
		return user, fmt.Errorf("unable to create user: %w", err)
	}

	user.Name = name
	user.Email = email
	user.Password = password 
	return user, nil

}

func (r *UserRepository) GetAll()([]service.User,error) {
	rows , err := r.db.Query("SELECT id,name, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("unable to get users: %w",err)
	}
	defer rows.Close()

	var users []service.User

	for rows.Next() {
		var user service.User
		if err := rows.Scan(&user.ID,&user.Name,&user.Email,&user.Password);err != nil{
			return nil, fmt.Errorf("unable to scan user: %w",err)
		}
		users = append(users, user)
	}

	return users, nil

}

func (r *UserRepository) GetByID(id string) (service.User, error) {
	var user service.User
	err := r.db.QueryRow("SELECT id, name, email, password FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return user, fmt.Errorf("unable to get user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Update(id, name, password, email string) (service.User, error) {
	
	_, err := r.db.Exec("UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4", name, email, password, id)
	if err != nil {
		return service.User{}, fmt.Errorf("unable to update user: %w", err)
	}

	
	return r.GetByID(id)
}

func (r *UserRepository) Delete(id string) error {

	_, err := r.db.Exec("DELETE FROM users WHERE id = $1",id)
	if err != nil {
		return fmt.Errorf("unable to delete user: %w", err)
	}
	return nil
}

