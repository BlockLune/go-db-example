package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	// _ "github.com/lib/pq"
)

var db *sql.DB

func initDb() {
	// get config from env variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	// construct the connection string
	// mysql
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	// postgres
	// connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// connect
	var err error
	db, err = sql.Open("mysql", connStr) // mysql
	// db, err = sql.Open("postgres", connStr) // postgres
	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to database!")
}

type User struct {
	ID       int64
	Email    string
	Password string
}

func createUsersTable() {
	query := `
    CREATE TABLE users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        email VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL
    )
    `

	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func addUser(email, password string) error {
	query := `
    INSERT INTO users (email, password)
    VALUES (?, ?)
    `

	_, err := db.Exec(query, email, password)
	if err != nil {
		return err
	}

	return nil
}

func addUserWithPrepareAndGetId(email, password string) (int64, error) {
	query := `
	INSERT INTO users (email, password)
	VALUES (?, ?)
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(email, password)
	if err != nil {
		return 0, err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func getUsers() ([]User, error) {
	query := `
    SELECT id, email, password
    FROM users
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func getUserByEmail(email string) (*User, error) {
	query := `
    SELECT id, email, password
    FROM users
    WHERE email = ?
    `

	row := db.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func updateUserPassword(email, newPassword string) error {
	query := `
    UPDATE users
    SET password = ?
    WHERE email = ?
    `

	_, err := db.Exec(query, newPassword, email)
	if err != nil {
		return err
	}

	return nil
}

func deleteUserByEmail(email string) error {
	query := `
    DELETE FROM users
    WHERE email = ?
    `

	_, err := db.Exec(query, email)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	initDb()
	createUsersTable()

	err := addUser("i@blocklune.cc", "password123")
	if err != nil {
		panic(err)
	}

	userId, err := addUserWithPrepareAndGetId("you@blocklune.cc", "password456")
	if err != nil {
		panic(err)
	}
	fmt.Println("New user ID:", userId)

	users, err := getUsers()
	if err != nil {
		panic(err)
	}
	fmt.Println("All users:")
	for _, user := range users {
		fmt.Printf("%d %s %s\n", user.ID, user.Email, user.Password)
	}

	user, err := getUserByEmail("i@blocklune.cc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("User with email i@blocklune.cc: %d %s %s\n", user.ID, user.Email, user.Password)

	err = updateUserPassword("i@blocklune.cc", "password789")
	if err != nil {
		panic(err)
	}

	user, err = getUserByEmail("i@blocklune.cc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("User with email i@blocklune.cc: %d %s %s\n", user.ID, user.Email, user.Password)

	err = deleteUserByEmail("you@blocklune.cc")
	if err != nil {
		panic(err)
	}

	users, err = getUsers()
	if err != nil {
		panic(err)
	}
	fmt.Println("All users:")
	for _, user := range users {
		fmt.Printf("%d %s %s\n", user.ID, user.Email, user.Password)
	}
}
