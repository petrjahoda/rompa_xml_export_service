package main

import "github.com/jinzhu/gorm"

func GetUserLoginForFail(fail TerminalInputFail) string {
	var user User
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return ""
	}
	db.Where("OID = ?", fail.UserID).Find(&user)
	return user.Login
}

func GetProductBarcodeForFail(productId int) string {
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

func GetOrderNameForFail(orderId int) (orderName string, productId int) {
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

func GetOrderFor(terminalInputOrderId int) int {
	var terminalInputOrder TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0
	}
	db.Where("OID = ?", terminalInputOrderId).Find(&terminalInputOrder)
	return terminalInputOrder.OrderID
}

func GetTerminalInputOrderIdForFail(terminalInputFail TerminalInputFail) int {
	var connection TerminalInputOrderTerminalInputFail
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0
	}
	db.Where("TerminalInputFailID = ?", terminalInputFail.OID).Find(&connection)
	return connection.TerminalInputOrderID
}

func GetFailBarcodeForFail(terminalInputFail TerminalInputFail) string {
	var fail Fail
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return ""
	}
	db.Where("OID = ?", terminalInputFail.FailID).Find(&fail)
	return fail.Barcode
}
