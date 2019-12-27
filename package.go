package main

import (
	"github.com/jinzhu/gorm"
)

func GetProductNameForPackage(productId int) string {
	var product Product
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return ""
	}
	defer db.Close()
	db.Where("OID = ?", productId).Find(&product)
	return product.Barcode
}

func GetProductIdForPackage(orderID int) int {
	var order Order
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0
	}
	defer db.Close()
	db.Where("OID = ?", orderID).Find(&order)
	return order.ProductID
}

func GetUserLoginForPackage(userId int) string {
	var user User
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return ""
	}
	defer db.Close()
	db.Where("OID = ?", userId).Find(&user)
	return user.Login
}

func GetOrderIdAndUserIdForPackage(inputPackage TerminalInputPackage) (orderId int, userId int, deviceId int) {
	var terminalInputOrder TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0, 0, 0
	}
	defer db.Close()
	db.Where("OID = ?", inputPackage.TerminalInputOrderID).Find(&terminalInputOrder)
	return terminalInputOrder.OrderID, terminalInputOrder.UserID, terminalInputOrder.DeviceID
}

func GetBarcodeForPackage(inputPackage TerminalInputPackage) int {
	var newPackage Package
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0
	}
	defer db.Close()
	db.Where("OID = ?", inputPackage.PackageID).Find(&newPackage)
	return newPackage.Barcode
}

func GetNewPackages(latestPackageId int) []TerminalInputPackage {
	var newPackages []TerminalInputPackage
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return nil
	}
	defer db.Close()
	db.Where("OID > ?", latestPackageId).Find(&newPackages)
	return newPackages
}

func (workplace Workplace) GetLatestPackageId() int {
	latestPackageId := RompaWorkplaceData{}
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError(workplace.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0
	}
	defer db.Close()
	db.Where("WorkplaceID = ?", workplace.OID).Find(&latestPackageId)
	return latestPackageId.LatestPackageID
}
