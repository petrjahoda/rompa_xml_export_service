package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strconv"
	"time"
)

type Fail struct {
	OID     int    `gorm:"column:OID"`
	Name    string `gorm:"column:Name"`
	Barcode string `gorm:"column:Barcode"`
}

func (Fail) TableName() string {
	return "fail"
}

type Workplace struct {
	OID              int    `gorm:"column:OID"`
	Name             string `gorm:"column:Name"`
	WorkplaceGroupID int    `gorm:"column:WorkplaceGroupID"`
	DeviceID         int    `gorm:"column:DeviceID"`
	Code             string `gorm:"column:Code"`
}

func (Workplace) TableName() string {
	return "workplace"
}

func (workplace Workplace) GetConstantFromDatabase() int {
	numberOfMillisecondsInMinute := 60000

	xmlTimerConstant := SwConfig{}
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError(workplace.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return 60 * numberOfMillisecondsInMinute
	}
	argument := "XmlExportInMinutes"
	db.Where("`Key` = ?", argument).Find(&xmlTimerConstant)
	returnValue, err := strconv.Atoi(xmlTimerConstant.Value)
	if err != nil {
		return 60 * numberOfMillisecondsInMinute
	}
	return returnValue * numberOfMillisecondsInMinute
}

type TerminalInputOrder struct {
	OID      int       `gorm:"column:OID"`
	DTS      time.Time `gorm:"column:DTS"`
	DTE      time.Time `gorm:"column:DTE"`
	OrderID  int       `gorm:"column:OrderID"`
	UserID   int       `gorm:"column:UserID"`
	DeviceID int       `gorm:"column:DeviceID"`
	Cavity   int       `gorm:"column:Cavity"`
}

func (TerminalInputOrder) TableName() string {
	return "terminal_input_order"
}

type TerminalInputOrderTerminalInputFail struct {
	TerminalInputOrderID int `gorm:"column:TerminalInputOrderID"`
	TerminalInputFailID  int `gorm:"column:TerminalInputFailID"`
}

func (TerminalInputOrderTerminalInputFail) TableName() string {
	return "terminal_input_order_terminal_input_fail"
}

type TerminalInputPackage struct {
	OID                  int       `gorm:"column:OID"`
	DT                   time.Time `gorm:"column:DT"`
	PackageID            int       `gorm:"column:PackageID"`
	TerminalInputOrderID int       `gorm:"column:TerminalInputOrderID"`
	Count                int       `gorm:"column:Count"`
}

func (TerminalInputPackage) TableName() string {
	return "terminal_input_package"
}

type WorkplacePort struct {
	OID          int    `gorm:"column:OID"`
	DevicePortID int    `gorm:"column:DevicePortID"`
	WorkplaceID  int    `gorm:"column:WorkplaceID"`
	Type         string `gorm:"column:Type"`
}

func (WorkplacePort) TableName() string {
	return "workplace_port"
}

type User struct {
	OID   int    `gorm:"column:OID"`
	Login string `gorm:"column:Login"`
}

func (User) TableName() string {
	return "user"
}

type Order struct {
	OID       int    `gorm:"column:OID"`
	Name      string `gorm:"column:Name"`
	ProductID int    `gorm:"column:ProductID"`
}

func (Order) TableName() string {
	return "order"
}

type Product struct {
	OID     int    `gorm:"column:OID"`
	Name    string `gorm:"column:Name"`
	Barcode string `gorm:"column:Barcode"`
}

func (Product) TableName() string {
	return "product"
}

type Package struct {
	OID           int    `gorm:"column:OID"`
	Barcode       int    `gorm:"column:Barcode"`
	PackageTypeID int    `gorm:"column:PackageTypeID"`
	OrderID       string `gorm:"column:OrderID"`
}

func (Package) TableName() string {
	return "package"
}

type TerminalInputFail struct {
	OID      int       `gorm:"column:OID"`
	DT       time.Time `gorm:"column:DT"`
	FailID   int       `gorm:"column:FailID"`
	UserID   int       `gorm:"column:UserID"`
	DeviceID int       `gorm:"column:DeviceID"`
	Count    int       `gorm:"column:Count"`
}

func (TerminalInputFail) TableName() string {
	return "terminal_input_fail"
}

type RompaWorkplaceData struct {
	OID                       int       `gorm:"column:OID"`
	WorkplaceID               string    `gorm:"column:WorkplaceID"`
	LatestTerminalFailID      int       `gorm:"column:LatestTerminalFailID"`
	LatestMachineFailDateTime time.Time `gorm:"column:LatestMachineFailDateTime"`
	LatestPackageID           int       `gorm:"column:LatestPackageID"`
}

func (RompaWorkplaceData) TableName() string {
	return "rompa_workplace_data"
}

type SwConfig struct {
	OID     int    `gorm:"column:OID"`
	SoftId  string `gorm:"column:SoftId"`
	Key     string `gorm:"column:Key"`
	Value   string `gorm:"column:Value"`
	Version string `gorm:"column:Version"`
	Note    string `gorm:"column:Note"`
}

func (SwConfig) TableName() string {
	return "sw_config"
}

func CheckDatabase() bool {
	var connectionString string
	var dialect string
	if DatabaseType == "postgres" {
		connectionString = "host=" + DatabaseIpAddress + " sslmode=disable port=" + DatabasePort + " user=" + DatabaseLogin + " dbname=" + DatabaseName + " password=" + DatabasePassword
		dialect = "postgres"
	} else if DatabaseType == "mysql" {
		connectionString = DatabaseLogin + ":" + DatabasePassword + "@tcp(" + DatabaseIpAddress + ":" + DatabasePort + ")/" + DatabaseName + "?charset=utf8&parseTime=True&loc=Local"
		dialect = "mysql"
	}
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogWarning("MAIN", "Database zapsi2 does not exist")
		return false
	}
	LogDebug("MAIN", "Database zapsi2 exists")
	return true
}

func CheckDatabaseType() (string, string) {
	var connectionString string
	var dialect string
	if DatabaseType == "postgres" {
		connectionString = "host=" + DatabaseIpAddress + " sslmode=disable port=" + DatabasePort + " user=" + DatabaseLogin + " dbname=" + DatabaseName + " password=" + DatabasePassword
		dialect = "postgres"
	} else if DatabaseType == "mysql" {
		connectionString = DatabaseLogin + ":" + DatabasePassword + "@tcp(" + DatabaseIpAddress + ":" + DatabasePort + ")/" + DatabaseName + "?charset=utf8&parseTime=True&loc=Local"
		dialect = "mysql"
	}
	return connectionString, dialect
}
