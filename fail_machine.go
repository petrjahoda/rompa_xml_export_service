package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

func GetUserLogin(userId int) string {
	var user User
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return ""
	}
	db.Where("OID = ?", userId).Find(&user)
	return user.Login
}

func GetOrdersClosedInLastTenSeconds(start time.Time, workplace Workplace) []TerminalInputOrder {
	var recentlyClosedOrders []TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return nil
	}
	db.Where("DeviceID = ?", workplace.DeviceID).Where("DTE > ?", start.Add(-10*time.Second)).Find(&recentlyClosedOrders)
	return recentlyClosedOrders
}

func GetActualOpenOrders(workplace Workplace) []TerminalInputOrder {
	var openOrders []TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return nil
	}
	db.Where("DeviceID = ?", workplace.DeviceID).Where("DTE is null").Find(&openOrders)
	return openOrders
}

func GetOrderNameAndProductIdFor(orderId int) (orderName string, productId int) {
	var order Order
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return "", 0
	}
	db.Where("OID = ?", orderId).Find(&order)
	return order.Name, order.ProductID
}

func GetProductBarcodeFor(productId int) string {
	var product Product
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return ""
	}
	db.Where("OID = ?", productId).Find(&product)
	return product.Barcode
}
