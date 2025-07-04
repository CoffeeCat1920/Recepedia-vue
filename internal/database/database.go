package database

import (
	"big/internal/modals"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {

	// User Functions
	AddUser(user *modals.User) error
	GetUserByName(name string) (*modals.User, error)
	GetUserByUUid(uuid string) (*modals.User, error)
	DeleteUserByUUid(uuid string) error
	NumberOfUsers() int
	GetAllUsers() ([]modals.User, error)

	// Session Functions
	AddSession(session *modals.Session) error
	GetSession(sessionId string) (*modals.Session, error)
	DeleteSession(sessionId string) error
	DeleteSessionByUser(ownerId string) error

	// Admin Sessions
	AddAdminSession(session *modals.AdminSession) error
	GetAdminSession(sessionId string) (*modals.AdminSession, error)

	// Recipe
	AddRecipe(recipe *modals.Recipe) error
	GetRecipe(UUID string) (*modals.Recipe, error)
	MostViewedRecipes() ([]modals.Recipe, error)
	IncreaseRecipeViews(recipe *modals.Recipe) error
	SearchRecipe(name string) ([]modals.Recipe, error)
	GetRecipesByUser(uuid string) ([]modals.Recipe, error)
	DeleteRecipe(uuid string) error
	EditRecipeName(uuid string, name string) error
	DeleteRecipeByUser(userUUid string) error
	NumberOfRecipes() int
	GetAllRecipes() ([]modals.Recipe, error)

	// Likes
	AddLike(like *modals.Like) error
	DeleteLikeFromUserRecipeId(userid string, recipeid string) error
	IsLiked(userid string, recipeid string) error

	Health() map[string]string

	Close() error
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	schema     = os.Getenv("BLUEPRINT_DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func NewTest(db *sql.DB) Service {
	return &service{db: db}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

func (s *service) doesExists(value, attribute, table string) bool {
	q := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1)", table, attribute)

	var exists bool
	err := s.db.QueryRow(q, value).Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}
