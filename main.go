package main

import (
	"github.com/jinzhu/gorm"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var (
	activeWorkplaces  []Workplace
	runningWorkplaces []Workplace
	workplaceSync     sync.Mutex
)

const version = "2020.4.3.2"
const deleteLogsAfter = 240 * time.Hour
const downloadInSeconds = 10

func main() {
	LogDirectoryFileCheck("MAIN")
	LogInfo("MAIN", "Program version "+version+" started")
	CreateConfigIfNotExists()
	LoadSettingsFromConfigFile()
	LogDebug("MAIN", "Using ["+DatabaseType+"] on "+DatabaseIpAddress+":"+DatabasePort+" with database "+DatabaseName)
	databaseAvailable := false
	networkFolderAvailable := false
	for {
		start := time.Now()
		LogInfo("MAIN", "Program running")
		databaseAvailable = CheckDatabase()
		if !networkFolderAvailable {
			networkFolderAvailable = MapNetworkFolder()
		}
		if databaseAvailable && networkFolderAvailable {
			UpdateActiveWorkplaces("MAIN")
			DeleteOldLogFiles()
			LogInfo("MAIN", "Active workplaces: "+strconv.Itoa(len(activeWorkplaces))+", running workplaces: "+strconv.Itoa(len(runningWorkplaces)))
			for _, activeWorkplace := range activeWorkplaces {
				activeWorkplaceIsRunning := CheckWorkplace(activeWorkplace)
				if !activeWorkplaceIsRunning {
					go RunWorkplace(activeWorkplace)
				}
			}

		}
		if time.Since(start) < (downloadInSeconds * time.Second) {
			sleepTime := downloadInSeconds*time.Second - time.Since(start)
			LogInfo("MAIN", "Sleeping for "+sleepTime.String())
			time.Sleep(sleepTime)
		}
	}
}

func RunWorkplace(workplace Workplace) {
	LogInfo(workplace.Name, "Workplace started running")
	workplaceSync.Lock()
	runningWorkplaces = append(runningWorkplaces, workplace)
	workplaceSync.Unlock()
	workplace.CheckRompaWorkplaceData()
	workplaceIsActive := true
	for workplaceIsActive {
		start := time.Now()
		workplace.ProcessTerminalBoxes()
		workplace.ProcessTerminalFails()
		workplace.ProcessMachineFails(start)
		LogInfo(workplace.Name, "Processing takes "+time.Since(start).String())
		workplace.Sleep(start)
		workplaceIsActive = CheckActive(workplace)
	}
	RemoveWorkplaceFromRunningWorkplaces(workplace)
	LogInfo(workplace.Name, "Workplace not active, stopped running")
}

func RemoveWorkplaceFromRunningWorkplaces(workplace Workplace) {
	for idx, runningWorkplace := range runningWorkplaces {
		if workplace.Name == runningWorkplace.Name {
			workplaceSync.Lock()
			runningWorkplaces = append(runningWorkplaces[0:idx], runningWorkplaces[idx+1:]...)
			workplaceSync.Unlock()
		}
	}
}

func CheckActive(workplace Workplace) bool {
	for _, activeWorkplace := range activeWorkplaces {
		if activeWorkplace.Name == workplace.Name {
			LogInfo(workplace.Name, "Workplace still active")
			return true
		}
	}
	LogInfo(workplace.Name, "Workplace not active")
	return false
}

func CheckWorkplace(workplace Workplace) bool {
	for _, runningWorkplace := range runningWorkplaces {
		if runningWorkplace.Name == workplace.Name {
			return true
		}
	}
	return false
}

func UpdateActiveWorkplaces(reference string) {
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError(reference, "Problem opening "+DatabaseName+" database: "+err.Error())
		activeWorkplaces = nil
		return
	}
	db.Where("WorkplaceGroupId = ?", 2).Find(&activeWorkplaces)
}

func MapNetworkFolder() bool {
	LogInfo("MAIN", "Creating directory data")
	cmd := exec.Command("mkdir", "/home/data")
	err := cmd.Run()
	if err != nil {
		LogError("MAIN", "Problem creating directory /home/data, already exists: "+err.Error())
	}
	LogInfo("MAIN", "Directory /home/data created")
	LogInfo("MAIN", "Mapping network directory")
	cmd = exec.Command("mount", "-t", "cifs", "-v", "-o", "username=zapsi,password=Jahoda123,domain=ROMPACZ", "//10.60.1.9/Interface/Zapsi2LN", "/home/data")
	err = cmd.Run()
	if err != nil {
		LogError("MAIN", "Problem mapping network directory: "+err.Error())
		return false
	}
	LogInfo("MAIN", "Network directory mapped successfully")
	return true
}
