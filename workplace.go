package main

import (
	"github.com/jinzhu/gorm"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (workplace Workplace) CheckRompaWorkplaceData() {
	var rompaWorkplaceData RompaWorkplaceData
	var terminalInputPackage TerminalInputPackage
	var terminalInputFail TerminalInputFail
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError(workplace.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Where("WorkplaceId = ?", workplace.OID).Find(&rompaWorkplaceData)
	db.Order("OID desc").Last(&terminalInputFail)
	db.Order("OID desc").Last(&terminalInputPackage)
	if rompaWorkplaceData.OID == 0 {
		LogInfo(workplace.Name, "No Record in RompaWorkplaceData, creating")
		LogInfo(workplace.Name, "Latest terminal input package: "+strconv.Itoa(terminalInputPackage.OID))
		LogInfo(workplace.Name, "Latest terminal input fail: "+strconv.Itoa(terminalInputFail.OID))
		rompaWorkplaceData.WorkplaceID = strconv.Itoa(workplace.OID)
		rompaWorkplaceData.LatestPackageID = terminalInputPackage.OID
		rompaWorkplaceData.LatestTerminalFailID = terminalInputFail.OID
		rompaWorkplaceData.LatestMachineFailDateTime = time.Now()
		db.Save(&rompaWorkplaceData)
	}
}

func (workplace Workplace) ProcessTerminalBoxes() {
	LogInfo(workplace.Name, "Processing packages")
	latestPackageId := workplace.GetLatestPackageId()
	LogInfo(workplace.Name, "Latest packageId: "+strconv.Itoa(latestPackageId))
	newPackages := GetNewPackages(latestPackageId)
	LogInfo(workplace.Name, "Number of new packages: "+strconv.Itoa(len(newPackages)))
	for _, newPackage := range newPackages {
		LogInfo(workplace.Name, "Processing package "+strconv.Itoa(newPackage.OID))
		barcode := GetBarcodeForPackage(newPackage)
		orderId, userId, deviceId := GetOrderIdAndUserIdForPackage(newPackage)
		userLogin := GetUserLoginForPackage(userId)
		productId := GetProductIdForPackage(orderId)
		productName := GetProductNameForPackage(productId)
		if deviceId == workplace.DeviceID {
			LogInfo(workplace.Name, "SAVING PACKAGE: "+newPackage.DT.String()+"-"+strconv.Itoa(barcode)+"-"+strconv.Itoa(newPackage.Count)+"-"+productName+"-"+userLogin+"-"+time.Now().String())
			data := "<?xml version=\"1.0\"?>" +
				"\n<ProcessBoxID>" +
				"\n\t<ApplicationArea>" +
				"\n\t\t<Sender>" +
				"\n\t\t\t<LogicalID>lid://infor.file.zapsi</LogicalID>" +
				"\n\t\t\t<ComponentID>erp</ComponentID>" +
				"\n\t\t\t<ConfirmationCode>OnError</ConfirmationCode>" +
				"\n\t\t</Sender>" +
				"\n\t\t<CreationDateTime>" + time.Now().Format("2006-01-02T15:04:05Z") + "</CreationDateTime>" +
				"\n\t\t<BODID>" + strconv.Itoa(barcode) + "</BODID>" +
				"\n\t</ApplicationArea>" +
				"\n\t<DataArea>" +
				"\n\t\t<Process>" +
				"\n\t\t\t<TenantID>infor</TenantID>" +
				"\n\t\t\t<AccountingEntityID>1000</AccountingEntityID>" +
				"\n\t\t\t<LocationID/>" +
				"\n\t\t\t<ActionCriteria>" +
				"\n\t\t\t\t<ActionExpression actionCode=\"Change\"/>" +
				"\n\t\t\t</ActionCriteria>" +
				"\n\t\t</Process>" +
				"\n\t\t<BoxID>" +
				"\n\t\t\t<BoxID>" +
				"\n\t\t\t\t<ID variationID=\"320\">" + strconv.Itoa(barcode) + "</ID>" +
				"\n\t\t\t</BoxID>" +
				"\n\t\t\t<ActualQuantity>" + strconv.Itoa(newPackage.Count) + "</ActualQuantity>" +
				"\n\t\t\t<Item>" +
				"\n\t\t\t\t<ID>" + productName + "</ID>" +
				"\n\t\t\t</Item>" +
				"\n\t\t\t<Operator>" + userLogin + "</Operator>" +
				"\n\t\t\t<ProductionDate>" + newPackage.DT.Format("2006-01-02T15:04:05Z") + "</ProductionDate>" +
				"\n\t\t\t<Status>completed</Status>" +
				"\n\t\t</BoxID>" +
				"\n\t</DataArea>" +
				"\n</ProcessBoxID>"

			logDirectory := filepath.Join("/home/data")
			logFileName := strconv.Itoa(barcode) + ".xml"
			logFullPath := strings.Join([]string{logDirectory, logFileName}, "/")
			f, err := os.OpenFile(logFullPath, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				LogError(workplace.Name, "Box saving: cannot open file: "+err.Error())
				return
			}
			if _, err := f.WriteString(data); err != nil {
				LogError(workplace.Name, "Cannot write data to file: "+err.Error())
			}
			f.Close()
		}
		workplace.UpdateRompaWorkplaceDataLatestPackageId(newPackage.OID)
	}
}

func (workplace Workplace) ProcessTerminalFails() {
	latestTerminalFailId := workplace.GetLatestTerminalFailId()
	terminalFails := workplace.GetTerminalFailsFrom(latestTerminalFailId)
	LogInfo(workplace.Name, "Number of terminal fails: "+strconv.Itoa(len(terminalFails)))
	if len(terminalFails) > 0 {
		for _, terminalFail := range terminalFails {
			rejectionReasonCode := GetFailBarcodeForFail(terminalFail)
			terminalInputOrderId := GetTerminalInputOrderIdForFail(terminalFail)
			orderId := GetOrderFor(terminalInputOrderId)
			productionSchedule, productId := GetOrderNameForFail(orderId)
			productName := GetProductBarcodeForFail(productId)
			userLogin := GetUserLoginForFail(terminalFail)
			LogInfo(workplace.Name, "SAVING TERMINAL FAIL: "+rejectionReasonCode)
			data := "<?xml version=\"1.0\"?>" +
				"\n<ProcessBoxID>" +
				"\n\t<ApplicationArea>" +
				"\n\t\t<Sender>" +
				"\n\t\t\t<LogicalID>lid://infor.file.zapsi</LogicalID>" +
				"\n\t\t\t<ComponentID>erp</ComponentID>" +
				"\n\t\t\t<ConfirmationCode>OnError</ConfirmationCode>" +
				"\n\t\t</Sender>" +
				"\n\t\t<CreationDateTime>" + time.Now().Format("2006-01-02T15:04:05Z") + "</CreationDateTime>" +
				"\n\t\t<BODID>" + strconv.Itoa(terminalFail.FailID) + "</BODID>" +
				"\n\t</ApplicationArea>" +
				"\n\t<DataArea>" +
				"\n\t\t<Process>" +
				"\n\t\t\t<TenantID>infor</TenantID>" +
				"\n\t\t\t<AccountingEntityID>1000</AccountingEntityID>" +
				"\n\t\t\t<LocationID/>" +
				"\n\t\t\t<ActionCriteria>" +
				"\n\t\t\t\t<ActionExpression actionCode=\"Add\"/>" +
				"\n\t\t\t</ActionCriteria>" +
				"\n\t\t</Process>" +
				"\n\t\t<BoxID>" +
				"\n\t\t\t<BoxID>" +
				"\n\t\t\t\t<ID></ID>" +
				"\n\t\t\t</BoxID>" +
				"\n\t\t\t<ProductionSchedule>" + productionSchedule + "</ProductionSchedule>" +
				"\n\t\t\t<Item>" +
				"\n\t\t\t\t<ID>" + productName + "</ID>" +
				"\n\t\t\t</Item>" +
				"\n\t\t\t<ActualQuantity>" + strconv.Itoa(terminalFail.Count) + "</ActualQuantity>" +
				"\n\t\t\t<Operator>" + userLogin + "</Operator>" +
				"\n\t\t\t<ProductionDate>" + terminalFail.DT.Format("2006-01-02T15:04:05Z") + "</ProductionDate>" +
				"\n\t\t\t<Status>rejected</Status>" +
				"\n\t\t\t<RejectReason>" + rejectionReasonCode + "</RejectReason>" +
				"\n\t\t</BoxID>" +
				"\n\t</DataArea>" +
				"\n</ProcessBoxID>"
			logDirectory := filepath.Join("/home/data")
			logFileName := strconv.Itoa(terminalFail.OID) + ".xml"
			logFullPath := strings.Join([]string{logDirectory, logFileName}, "/")
			f, err := os.OpenFile(logFullPath, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				LogError(workplace.Name, "Box saving: create file: "+err.Error())
				return
			}
			if _, err := f.WriteString(data); err != nil {
				LogError(workplace.Name, "Cannot write data to file: "+err.Error())
			}
			latestNumber := terminalFail.OID + terminalFail.Count - 1
			workplace.UpdateRompaWorkplaceDataLatestTerminalFailId(latestNumber)
			f.Close()

		}
	}
}

func (workplace Workplace) ProcessMachineFails(start time.Time) {
	failDevicePortId := workplace.GetFailDevicePortId()
	machineFailId := "101"
	rejectionReasonCode := "R0402"
	latestMachineFailDateTime := workplace.GetLatestMachineFailDateTime()
	latestOrders := GetOrdersClosedInLastTenSeconds(start, workplace)
	machineFailsCount := workplace.GetNumberOfMachineFails(latestMachineFailDateTime, start, failDevicePortId)
	openOrders := GetActualOpenOrders(workplace)
	LogInfo(workplace.Name, "Number of open orders: "+strconv.Itoa(len(openOrders)))
	LogInfo(workplace.Name, "Number of closed orders in last 10 seconds: "+strconv.Itoa(len(latestOrders))+", ready to transfer : "+strconv.FormatBool(len(latestOrders) > 0))
	LogInfo(workplace.Name, "Time from latest machine fail transfer: "+start.Sub(latestMachineFailDateTime).String()+", ready to transfer : "+strconv.FormatBool(start.Sub(latestMachineFailDateTime) > time.Minute*60))
	LogInfo(workplace.Name, "Number of machine fails: "+strconv.Itoa(machineFailsCount)+", ready to transfer : "+strconv.FormatBool(machineFailsCount > 0))
	if len(latestOrders) > 0 {
		for _, openOrder := range latestOrders {
			if machineFailsCount > 0 {
				productionSchedule, productId := GetOrderNameAndProductIdFor(openOrder.OrderID)
				productName := GetProductBarcodeFor(productId)
				userLogin := GetUserLogin(openOrder.UserID)
				println("Order name : " + productionSchedule)
				LogInfo(workplace.Name, "SAVING MACHINE FAIL: "+rejectionReasonCode)
				data := "<?xml version=\"1.0\"?>" +
					"\n<ProcessBoxID>" +
					"\n\t<ApplicationArea>" +
					"\n\t\t<Sender>" +
					"\n\t\t\t<LogicalID>lid://infor.file.zapsi</LogicalID>" +
					"\n\t\t\t<ComponentID>erp</ComponentID>" +
					"\n\t\t\t<ConfirmationCode>OnError</ConfirmationCode>" +
					"\n\t\t</Sender>" +
					"\n\t\t<CreationDateTime>" + start.Format("2006-01-02T15:04:05Z") + "</CreationDateTime>" +
					"\n\t\t<BODID>" + machineFailId + "</BODID>" +
					"\n\t</ApplicationArea>" +
					"\n\t<DataArea>" +
					"\n\t\t<Process>" +
					"\n\t\t\t<TenantID>infor</TenantID>" +
					"\n\t\t\t<AccountingEntityID>1000</AccountingEntityID>" +
					"\n\t\t\t<LocationID/>" +
					"\n\t\t\t<ActionCriteria>" +
					"\n\t\t\t\t<ActionExpression actionCode=\"Add\"/>" +
					"\n\t\t\t</ActionCriteria>" +
					"\n\t\t</Process>" +
					"\n\t\t<BoxID>" +
					"\n\t\t\t<BoxID>" +
					"\n\t\t\t\t<ID></ID>" +
					"\n\t\t\t</BoxID>" +
					"\n\t\t\t<ProductionSchedule>" + productionSchedule + "</ProductionSchedule>" +
					"\n\t\t\t<Item>" +
					"\n\t\t\t\t<ID>" + productName + "</ID>" +
					"\n\t\t\t</Item>" +
					"\n\t\t\t<ActualQuantity>" + strconv.Itoa(machineFailsCount*openOrder.Cavity) + "</ActualQuantity>" +
					"\n\t\t\t<Operator>" + userLogin + "</Operator>" +
					"\n\t\t\t<ProductionDate>" + start.Format("2006-01-02T15:04:05Z") + "</ProductionDate>" +
					"\n\t\t\t<Status>rejected</Status>" +
					"\n\t\t\t<RejectReason>" + rejectionReasonCode + "</RejectReason>" +
					"\n\t\t</BoxID>" +
					"\n\t</DataArea>" +
					"\n</ProcessBoxID>"
				logDirectory := filepath.Join("/home/data")
				logFileName := strconv.Itoa(workplace.OID) + strconv.Itoa(start.Year()) + strconv.Itoa(start.YearDay()) + strconv.Itoa(start.Hour()) + strconv.Itoa(start.Minute()) + "-" + productName + ".xml"
				logFullPath := strings.Join([]string{logDirectory, logFileName}, "/")
				f, err := os.OpenFile(logFullPath, os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					LogError(workplace.Name, "Machine fails saving: create file: "+err.Error())
					return
				}
				if _, err := f.WriteString(data); err != nil {
					LogError(workplace.Name, "Cannot write data to file: "+err.Error())
				}
				f.Close()
			}
			workplace.UpdateRompaWorkplaceDataMachineFailDateTime(start)

		}
	}
	latestMachineFailDateTime = workplace.GetLatestMachineFailDateTime()
	waitTime := workplace.GetConstantFromDatabase()
	LogInfo(workplace.Name, "Duration now: "+strconv.Itoa(int(start.Sub(latestMachineFailDateTime).Milliseconds())))
	LogInfo(workplace.Name, "Duration set: "+strconv.Itoa(waitTime))
	LogInfo(workplace.Name, "Duration is over the limit: "+strconv.FormatBool(int(start.Sub(latestMachineFailDateTime).Milliseconds()) > waitTime))
	if (int(start.Sub(latestMachineFailDateTime).Milliseconds()) > waitTime) && len(openOrders) > 0 {
		LogInfo(workplace.Name, "Processing machine fails, number of machine fails: "+strconv.Itoa(machineFailsCount))
		for _, openOrder := range openOrders {
			if machineFailsCount > 0 {
				productionSchedule, productId := GetOrderNameAndProductIdFor(openOrder.OrderID)
				productBarcode := GetProductBarcodeFor(productId)
				userLogin := GetUserLogin(openOrder.UserID)
				LogInfo(workplace.Name, "SAVING MACHINE FAIL: "+rejectionReasonCode)
				data := "<?xml version=\"1.0\"?>" +
					"\n<ProcessBoxID>" +
					"\n\t<ApplicationArea>" +
					"\n\t\t<Sender>" +
					"\n\t\t\t<LogicalID>lid://infor.file.zapsi</LogicalID>" +
					"\n\t\t\t<ComponentID>erp</ComponentID>" +
					"\n\t\t\t<ConfirmationCode>OnError</ConfirmationCode>" +
					"\n\t\t</Sender>" +
					"\n\t\t<CreationDateTime>" + time.Now().Format("2006-01-02T03:04:05Z") + "</CreationDateTime>" +
					"\n\t\t<BODID>" + machineFailId + "</BODID>" +
					"\n\t</ApplicationArea>" +
					"\n\t<DataArea>" +
					"\n\t\t<Process>" +
					"\n\t\t\t<TenantID>infor</TenantID>" +
					"\n\t\t\t<AccountingEntityID>1000</AccountingEntityID>" +
					"\n\t\t\t<LocationID/>" +
					"\n\t\t\t<ActionCriteria>" +
					"\n\t\t\t\t<ActionExpression actionCode=\"Add\"/>" +
					"\n\t\t\t</ActionCriteria>" +
					"\n\t\t</Process>" +
					"\n\t\t<BoxID>" +
					"\n\t\t\t<BoxID>" +
					"\n\t\t\t\t<ID></ID>" +
					"\n\t\t\t</BoxID>" +
					"\n\t\t\t<ProductionSchedule>" + productionSchedule + "</ProductionSchedule>" +
					"\n\t\t\t<Item>" +
					"\n\t\t\t\t<ID>" + productBarcode + "</ID>" +
					"\n\t\t\t</Item>" +
					"\n\t\t\t<ActualQuantity>" + strconv.Itoa(machineFailsCount*openOrder.Cavity) + "</ActualQuantity>" +
					"\n\t\t\t<Operator>" + userLogin + "</Operator>" +
					"\n\t\t\t<ProductionDate>" + start.Format("2006-01-02T03:04:05Z") + "</ProductionDate>" +
					"\n\t\t\t<Status>rejected</Status>" +
					"\n\t\t\t<RejectReason>" + rejectionReasonCode + "</RejectReason>" +
					"\n\t\t</BoxID>" +
					"\n\t</DataArea>" +
					"\n</ProcessBoxID>"
				logDirectory := filepath.Join("/home/data")
				logFileName := strconv.Itoa(workplace.OID) + strconv.Itoa(start.Year()) + strconv.Itoa(start.YearDay()) + strconv.Itoa(start.Hour()) + strconv.Itoa(start.Minute()) + "-" + productBarcode + ".xml"
				logFullPath := strings.Join([]string{logDirectory, logFileName}, "/")
				f, err := os.OpenFile(logFullPath, os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					LogError(workplace.Name, "Box saving: create file: "+err.Error())
					return
				}
				if _, err := f.WriteString(data); err != nil {
					LogError(workplace.Name, "Cannot write data to file: "+err.Error())
				}
				f.Close()
			}
			workplace.UpdateRompaWorkplaceDataMachineFailDateTime(start)
		}
	}
}

func (workplace Workplace) UpdateRompaWorkplaceDataMachineFailDateTime(start time.Time) {
	var rompaWorkplaceData RompaWorkplaceData
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError(workplace.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Model(&rompaWorkplaceData).Where("WorkplaceId = ?", workplace.OID).Update("LatestMachineFailDateTime", start)
}

func (workplace Workplace) UpdateRompaWorkplaceDataLatestPackageId(latestPackageId int) {
	var rompaWorkplaceData RompaWorkplaceData
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError(workplace.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Model(&rompaWorkplaceData).Where("WorkplaceId = ?", workplace.OID).Update("LatestPackageID", latestPackageId)
}

func (workplace Workplace) UpdateRompaWorkplaceDataLatestTerminalFailId(latestTerminalFailId int) {
	var rompaWorkplaceData RompaWorkplaceData
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError(workplace.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
	}
	if db != nil {
		defer db.Close()
		db.Model(&rompaWorkplaceData).Where("WorkplaceId = ?", workplace.OID).Update("LatestTerminalFailID", latestTerminalFailId)
	}
}

func (workplace Workplace) Sleep(start time.Time) {
	if time.Since(start) < (downloadInSeconds * time.Second) {
		sleepTime := downloadInSeconds*time.Second - time.Since(start)
		LogInfo(workplace.Name, "Sleeping for "+sleepTime.String())
		time.Sleep(sleepTime)
	}
}

func (workplace Workplace) GetLatestTerminalFailId() int {
	workplaceData := RompaWorkplaceData{}
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError(workplace.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0
	}
	defer db.Close()
	db.Where("WorkplaceID = ?", workplace.OID).Find(&workplaceData)
	return workplaceData.LatestTerminalFailID
}

func (workplace Workplace) GetLatestMachineFailDateTime() time.Time {
	workplaceData := RompaWorkplaceData{}
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError(workplace.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.Now().Location())
	}
	defer db.Close()
	db.Where("WorkplaceID = ?", workplace.OID).Find(&workplaceData)
	return workplaceData.LatestMachineFailDateTime
}

func (workplace Workplace) CheckOrderChangeInLastTenSeconds(start time.Time) bool {
	terminalInputOrder := TerminalInputOrder{}
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError(workplace.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return false
	}
	defer db.Close()
	db.Where("DeviceID = ?", workplace.DeviceID).Order("DTS desc").Last(&terminalInputOrder)
	LogInfo(workplace.Name, "Latest open order id: "+strconv.Itoa(terminalInputOrder.OID))
	LogInfo(workplace.Name, terminalInputOrder.DTE.String())
	LogInfo(workplace.Name, start.String())
	LogInfo(workplace.Name, start.Sub(terminalInputOrder.DTE).String())
	LogInfo(workplace.Name, strconv.FormatBool(terminalInputOrder.DTE.Before(start)))
	if start.Sub(terminalInputOrder.DTE) < 10*time.Second && terminalInputOrder.DTE.Before(start) {
		LogInfo(workplace.Name, "Order closed inside interval of last 10 seconds")
		return true
	}
	return false
}

func (workplace Workplace) GetTerminalFailsFrom(terminalInputFailId int) []TerminalInputFail {
	var terminalFails []TerminalInputFail
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError(workplace.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return terminalFails
	}
	defer db.Close()
	db.Raw("SELECT count(*) as Count,OID, FailID, DT, UserID, DeviceID FROM `terminal_input_fail`where OID>? and DeviceID =? group by FailID,DT order by OID asc", terminalInputFailId, workplace.DeviceID).Find(&terminalFails)
	return terminalFails

}

func (workplace Workplace) GetNumberOfMachineFails(latestMachineFailDateTime time.Time, endTime time.Time, failDevicePortId int) int {
	var count int
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError(workplace.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0
	}
	defer db.Close()
	db.Table("device_input_digital").Where("DevicePortID = ?", failDevicePortId).Where("DT > ?", latestMachineFailDateTime).Where("DT < ?", endTime).Where("Data =?", 1).Count(&count)
	return count
}

func (workplace Workplace) GetFailDevicePortId() int {
	var workplacePort WorkplacePort
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError(workplace.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0
	}
	defer db.Close()
	db.Where("WorkplaceID = ?", workplace.OID).Where("Type = ?", "fail").Find(&workplacePort)
	return workplacePort.DevicePortID
}
