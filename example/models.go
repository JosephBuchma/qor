package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/qor/qor/admin"
	"github.com/qor/qor/l10n"
	"github.com/qor/qor/media_library"
	"github.com/qor/qor/publish"
)

type CreditCard struct {
	gorm.Model
	Number string
	Issuer string
}

type Address struct {
	gorm.Model
	UserId   int64
	Address1 string
	Address2 string
}

type Role struct {
	gorm.Model
	Name string
}

type Language struct {
	gorm.Model
	Name string
}

type User struct {
	gorm.Model
	Name         string
	Gender       string
	Birthday     *time.Time
	Description  string `sql:"size:622550"`
	File         media_library.FileSystem
	RoleID       uint
	Languages    []Language `gorm:"many2many:user_languages;"`
	CreditCard   CreditCard
	CreditCardID uint
	Addresses    []Address
	publish.Status
}

func (User) ViewableLocales() []string {
	return []string{l10n.Global, "zh-CN", "JP", "EN", "DE"}
}

func (user User) EditableLocales() []string {
	if user.Name == "global_admin" {
		return []string{l10n.Global, "zh-CN", "EN"}
	} else {
		return []string{"zh-CN", "EN"}
	}
}

func (u User) DisplayName() string {
	return u.Name
}

type Product struct {
	gorm.Model
	Name         *string
	Description  *string
	CollectionID uint
	l10n.Locale
	publish.Status
}

type Collection struct {
	gorm.Model
	Name string
}

type Order struct {
	gorm.Model
	UserID     uint
	Status     string
	Amount     float64
	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Price     float64
	Quantity  uint
}

type DiscountRule struct {
	gorm.Model
	DiscountID string
	Type       string
	Value      string
}

type DiscountBenefit struct {
	gorm.Model
	DiscountID string
	Type       string
	Value      string
}

type DiscountCode struct {
	gorm.Model
	DiscountID     string
	Code           string
	AvailableTimes uint
}

type Discount struct {
	gorm.Model
	Name     string
	Benefits []DiscountBenefit
	Rules    []DiscountRule
	Codes    []DiscountCode
	Begins   *time.Time
	Expires  *time.Time
}

var DB gorm.DB
var Publish *publish.Publish

func init() {
	var err error
	// CREATE USER 'qor'@'localhost' IDENTIFIED BY 'qor';
	// CREATE DATABASE qor_example;
	// GRANT ALL PRIVILEGES ON qor_example.* TO 'qor'@'localhost';
	dbuser, dbpwd := "qor", "qor"
	if os.Getenv("WEB_ENV") == "online" {
		dbuser, dbpwd = os.Getenv("DB_USER"), os.Getenv("DB_PWD")
	}
	DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@/qor_example?charset=utf8&parseTime=True&loc=Local", dbuser, dbpwd))
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&User{}, &CreditCard{}, &Address{}, &Role{}, &Language{}, &Product{}, &Collection{}, &Order{}, &OrderItem{}, &admin.AssetManager{})

	Publish = publish.New(&DB)
	Publish.Support(&Product{}).AutoMigrate()

	l10n.RegisterCallbacks(&DB)

	var AdminRole Role
	DB.FirstOrCreate(&AdminRole, Role{Name: "admin"})
	DB.FirstOrCreate(&Role{}, Role{Name: "dev"})
	DB.FirstOrCreate(&Role{}, Role{Name: "customer_support"})

	DB.FirstOrCreate(&User{}, User{Name: "admin", RoleID: AdminRole.ID})
	DB.FirstOrCreate(&User{}, User{Name: "global_admin", RoleID: AdminRole.ID})

	DB.FirstOrCreate(&Language{}, Language{Name: "CN"})
	DB.FirstOrCreate(&Language{}, Language{Name: "JP"})
	DB.FirstOrCreate(&Language{}, Language{Name: "EN"})
	DB.FirstOrCreate(&Language{}, Language{Name: "DE"})

	DB.LogMode(true)
}
