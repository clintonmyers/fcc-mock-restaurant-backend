package helpers

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"testing"
)

func CreateAndMigrateTempDB(dir string) (*os.File, *gorm.DB, error) {
	file, err := ioutil.TempFile(dir, "test_database")
	if err != nil {
		return nil, nil, err
	}

	db, err := gorm.Open(sqlite.Open(file.Name()), &gorm.Config{
		PrepareStmt:          true,
		FullSaveAssociations: true,
	})
	if err != nil {
		defer file.Close()
		return nil, nil, err
	}

	// Migrate all models
	// By migrating them for all we can catch modeling errors introduced by changes
	//MigrateDB(&app.Configuration{DB: db})

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.UserRole{})
	db.AutoMigrate(&models.Testimonial{})
	db.AutoMigrate(&models.TestimonialImage{})
	db.AutoMigrate(&models.Restaurant{})
	db.AutoMigrate(&models.UserAddress{})
	db.AutoMigrate(&models.RestaurantAddress{})
	db.AutoMigrate(&models.Company{})
	db.AutoMigrate(&models.MenuImage{})
	db.AutoMigrate(&models.Menu{})
	db.AutoMigrate(&models.Role{})
	db.AutoMigrate(&models.MenuItem{})

	return file, db, err
}
func TestMainRepository_GetUserBySubscriberId(t *testing.T) {

	tempDir := t.TempDir()
	file, db, err := CreateAndMigrateTempDB(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var mainRepo UserRepository
	mainRepo = &MainRepository{DB: db}

	// -------------------------------- //
	// -------------------------------- //
	// -------------------------------- //

	// Create a user
	user := models.User{
		Username:  "Username",
		FirstName: "first",
		LastName:  "last",
		SubId:     "subscriber1",
		UserRole: []models.UserRole{
			{
				Model:        gorm.Model{},
				Role:         "admin",
				RestaurantID: 1,
			},
		},
		Addresses: []models.UserAddress{},
	}

	count, err := mainRepo.SaveUser(&user)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("expected to affect one row with insert")
	}

	// -------------------------------- //
	// -------------------------------- //
	// -------------------------------- //

	var ans models.User
	err = mainRepo.GetUserBySubId(&ans, "subscriber1")
	if err != nil {
		t.Fatal(err)
	}

	users := make([]models.User, 0)
	mainRepo.GetAllUsers(&users)

	for _, a := range users {
		fmt.Println(a)
	}

	if ans.ID != 1 {
		t.Fatal("expected ID to be returned to equal 1")
	}
	for _, role := range ans.UserRole {
		if role.RestaurantID != 1 {
			t.Fatal("invalid restaurant ID found")
		}
		if role.Role != "admin" {
			t.Fatal("invalid user role found")
		}
	}
}

func TestMainRepository_GetUserById(t *testing.T) {

	tempDir := t.TempDir()
	file, db, err := CreateAndMigrateTempDB(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var mainRepo UserRepository
	mainRepo = &MainRepository{DB: db}

	// -------------------------------- //
	// -------------------------------- //
	// -------------------------------- //

	// Create a user
	user := models.User{
		Username:  "Username",
		FirstName: "first",
		LastName:  "last",
		SubId:     "abc",
		UserRole: []models.UserRole{
			{
				Role: "role1",
			},
		},
		Addresses: []models.UserAddress{},
	}

	count, err := mainRepo.SaveUser(&user)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("expected to affect one row with insert")
	}

	// -------------------------------- //
	// -------------------------------- //
	// -------------------------------- //

	var shouldBeIdOne models.User
	err = mainRepo.GetUserById(&shouldBeIdOne, 1)
	if err != nil {
		t.Fatal(err)
	}

	if shouldBeIdOne.ID != 1 {
		t.Fatal("expected ID to be returned to equal 1")
	}
}

func TestMainRepository_GetAllUsers(t *testing.T) {

	tempDir := t.TempDir()
	file, db, err := CreateAndMigrateTempDB(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var mainRepo UserRepository
	mainRepo = &MainRepository{DB: db}
	//
	//// Migrate the DB tables
	//db.AutoMigrate(&models.User{})
	//db.AutoMigrate(&models.UserRole{})

	// -------------------------------- //
	// -------------------------------- //
	// -------------------------------- //

	// Create a user
	user := models.User{
		Username:  "Username",
		FirstName: "first",
		LastName:  "last",
		SubId:     "subscriber1",
		UserRole: []models.UserRole{
			{
				Model:        gorm.Model{},
				Role:         "role1",
				RestaurantID: 1,
			},
		},
		Addresses: []models.UserAddress{},
	}

	count, err := mainRepo.SaveUser(&user)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("expected to affect one row with insert")
	}

	// -------------------------------- //
	// -------------------------------- //
	// -------------------------------- //

	users := make([]models.User, 0)

	//db.Find(&users)
	err = mainRepo.GetAllUsers(&users)
	//err = mainRepo.GetAllUsers(&users)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(users)

	if len(users) != 1 {
		t.Fatal("expected to only find one row inserted")
	}
}

func TestMainRepository_GetUserByUsername(t *testing.T) {

	tempDir := t.TempDir()
	file, db, err := CreateAndMigrateTempDB(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var mainRepo UserRepository
	mainRepo = &MainRepository{DB: db}
	//
	//// Migrate the DB tables
	//db.AutoMigrate(&models.User{})
	//db.AutoMigrate(&models.UserRole{})

	// -------------------------------- //
	// -------------------------------- //
	// -------------------------------- //

	// Create a user
	user := models.User{
		Username:  "Username",
		FirstName: "first",
		LastName:  "last",
		UserRole: []models.UserRole{
			{
				Role: "role1",
			},
		},
		Addresses: []models.UserAddress{},
	}

	count, err := mainRepo.SaveUser(&user)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("expected to affect one row with insert")
	}

	// -------------------------------- //
	// -------------------------------- //
	// -------------------------------- //

	var shouldBeIdOne models.User
	err = mainRepo.GetUserByUsername(&shouldBeIdOne, "Username")
	if err != nil {
		t.Fatal(err)
	}

	if shouldBeIdOne.ID != 1 {
		t.Fatal("expected ID to be returned to equal 1")
	}
}

func TestMainRepository_SaveUser(t *testing.T) {
	tempDir := t.TempDir()
	file, db, err := CreateAndMigrateTempDB(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var mainRepo UserRepository
	mainRepo = &MainRepository{DB: db}

	//// Migrate the DB tables
	//db.AutoMigrate(&models.User{})
	//db.AutoMigrate(&models.UserRole{})
	//
	// In order to test we just want to make sure that we can query a user from the db

	// Create a user
	user := models.User{
		Username:  "Username",
		FirstName: "first",
		LastName:  "last",
		UserRole: []models.UserRole{
			{
				Role: "role1",
			},
		},
		Addresses: []models.UserAddress{},
	}

	count, err := mainRepo.SaveUser(&user)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("expected to affect one row with insert")
	}

	users := make([]models.User, 0)

	//db.Find(&users)
	err = mainRepo.GetAllUsers(&users)
	//err = mainRepo.GetAllUsers(&users)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(users)

	if len(users) != 1 {
		t.Fatal("expected to only find one row inserted")
	}
	//var shouldBeIdOne models.User
	//err = mainRepo.GetUserById(&shouldBeIdOne, 1)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//if shouldBeIdOne.ID != 1 {
	//	t.Fatal("expected ID to be returned to equal 1")
	//}
	//
	//var getByUsername models.User
	//err = mainRepo.GetUserByUsername(&getByUsername, "Username")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//if getByUsername.ID != 1 {
	//	t.Fatal(err)
	//}

}
