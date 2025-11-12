package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/cheggaaa/pb/v3"
)

/*
Version: 0.0.1
Author: Mang Zhang, Shenzhen China
Release Date: 2025-11-11
Project Name: GoPercyUpgrade
Description: A library is for application upgrade, only one line code is enough for
			Mac, Linux, Windows three platform application upgrade.
Copy Rights: MIT License
Email: m13692277450@outlook.com
Mobile: +86-13692277450
HomePage: www.pavogroup.top , github.com/13692277450
*/

/*
Version.json file structure:
{
  "versionwindows": "1.0.1",
  "versionlinux": "1.0.1",
  "versionmac": "1.0.1",
  "downloadUrlwindows": "http://example.com/windows",
  "downloadUrllinux": "http://example.com/linux",
  "downloadUrlmac": "http://example.com/mac",
  "noteswindows": "Bug fixes and improvements",
  "noteslinux": "Bug fixes and improvements",
  "notesmac": "Bug fixes and improvements",
  "pub_date": "2025-11-09T00:00:00Z",
  "shownotesmessages": true,
  "removeoldfiles": true
}
*/

type UpgradeConfig struct {
	VersionWindows    string `json:"versionwindows"`
	VersionLinux      string `json:"versionlinux"`
	VersionMac        string `json:"versionmac"`
	DownloadUrlWinows string `json:"downloadUrlwindows"`
	DownloadUrlLinux  string `json:"downloadUrllinux"`
	DownloadUrlMac    string `json:"downloadUrlmac"`
	// Sha256Windows        string `json:"sha256windows"`
	// Sha256Linux          string `json:"sha256linux"`
	// Sha256Mac            string `json:"sha256mac"`
	NotesWindows         string `json:"noteswindows"`
	NotesLinux           string `json:"noteslinux"`
	NotesMac             string `json:"notesmac"`
	PubDate              string `json:"pub_date"`
	VeriftyAfterDownload bool   `json:"veriftyafterdownload"`
	// ShowProgressBar      bool   `json:"showprogressbar"`
	// BackgroundMode       bool   `json:"backgroundmode"`
	ShowNotesMessages bool `json:"shownotesmessages"`

	RemoveOldFiles bool `json:"removeoldfiles"`
}

var (
	cyan         = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	green        = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	downloadUrl  string
	updateConfig UpgradeConfig
)

type Version struct {
	MajorVersion int
	MinorVersion int
	PatchVersion int
}

// Helper function to validate version format (x.x.x)

func GoPercyUpgradeConfig(CurrentVersion, NewVersion string) {

	// Get update version info from server.GoPercyUpgrade
	versionBody, err := http.Get(NewVersion)
	if err != nil {
		fmt.Printf("ERROR: Get update version info failed: %v\n", err)
		return
	}
	defer versionBody.Body.Close()
	versionData, err := io.ReadAll(versionBody.Body)
	if err != nil {
		fmt.Printf("ERROR: Read update version info failed: %v\n", err)
		return
	}
	err = json.Unmarshal(versionData, &updateConfig)
	if err != nil {
		fmt.Printf("ERROR: JSON unmarshal failed: %v\n", err)
		return
	}

	sysType := runtime.GOOS
	switch sysType {
	case "windows":

		versionDiffer := VersionCompareResult(CurrentVersion, updateConfig.VersionWindows)
		switch versionDiffer {
		case "lower":
			return
		case "equal":
			return
		case "newer":
			//need upgrade
			downloadUrl = updateConfig.DownloadUrlWinows
			if downloadUrl == "" { //If no Windows download URL, just return
				return
			}
			if updateConfig.ShowNotesMessages {
				fmt.Println(green.Render("Upgrade version: ", updateConfig.VersionWindows))
				fmt.Println(cyan.Render("Upgrade notes: ", updateConfig.NotesWindows))
				fmt.Println(green.Render("Publish data: ", updateConfig.PubDate))
			}
			upgradeWindows(downloadUrl)
		}
	case "linux":

		versionDiffer := VersionCompareResult(CurrentVersion, updateConfig.VersionLinux)
		switch versionDiffer {
		case "lower":
			return
		case "equal":
			return
		case "newer":
			downloadUrl = updateConfig.DownloadUrlLinux
			if downloadUrl == "" { //If no Linux download URL, just return
				return
			}
			if updateConfig.ShowNotesMessages {
				fmt.Println(green.Render("Upgrade version: ", updateConfig.VersionLinux))
				fmt.Println(cyan.Render("Upgrade notes: ", updateConfig.NotesLinux))
				fmt.Println(green.Render("Publish data: ", updateConfig.PubDate))
			}
			upgradeLinux(downloadUrl)
		}
	case "darwin":

		versionDiffer := VersionCompareResult(CurrentVersion, updateConfig.VersionMac)
		switch versionDiffer {
		case "lower":
			return
		case "equal":
			return
		case "newer":
			downloadUrl = updateConfig.DownloadUrlMac
			if downloadUrl == "" { //If no Mac download URL, just return
				return
			}
			if updateConfig.ShowNotesMessages {
				fmt.Println(green.Render("Upgrade version: ", updateConfig.VersionMac))
				fmt.Println(cyan.Render("Upgrade notes: ", updateConfig.NotesMac))
				fmt.Println(green.Render("Publish data: ", updateConfig.PubDate))
			}
			upgradeMac(downloadUrl)
		}
	default:
		fmt.Println("Your OS is not supported.")
	}
}
func upgradeMac(downloadUrl string) { //Mac version upgrade function
	filePathMac := filepath.Base(os.Args[0])
	go func() {
		fmt.Println(cyan.Render("Starting download upgrade from: ", downloadUrl+"\n"))
		for i := 1; i < 15; i++ {
			fmt.Print(".")
			time.Sleep(500 * time.Millisecond)
		}
	}()

	tempFile := filePathMac + ".tmp" //move current application to .tmp file
	err := downloadFile(downloadUrl, tempFile)
	if err != nil {
		fmt.Println("Download new file error.")
		return
	}

	// Rename current executable to .old and rename the new one to current executable
	oldFile := filePathMac + ".old"
	if updateConfig.RemoveOldFiles {
		os.Remove(oldFile) // Remove old backup if exists
	}
	if err := os.Rename(filePathMac, oldFile); err != nil {
		// It's okay if the current file doesn't exist (first install)
		fmt.Printf("Move to old file was failure: %v\n", err)
	}

	// Create update shell script for macOS (same behavior as Linux)
	shContent := "#!/bin/sh\n" +
		"sleep 2\n" +
		"mv \"" + tempFile + "\" \"" + filePathMac + "\"\n" +
		"chmod +x \"" + filePathMac + "\"\n" +
		"rm -- \"$0\"\n"

	shFile := "update_mac.sh"
	if err := os.WriteFile(shFile, []byte(shContent), 0755); err != nil {
		fmt.Println("Writing update script was failure.")
		return
	}

	// Run the shell script and wait a moment for it to take over
	fmt.Println("Executing update script...")
	cmd := exec.Command("sh", shFile)
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error executing update script: %v\n", err)
		return
	}
	time.Sleep(3 * time.Second)

	// Verify the new file exists
	if _, err := os.Stat(filePathMac); os.IsNotExist(err) {
		fmt.Printf("Error: New file %s was not created\n", filePathMac)
		return
	}

	fmt.Println("Update completed successfully. Old version saved as " + oldFile + ", the old version will be removed automatically when application launch next time.")
	os.Exit(0) // Exit the program after successful update
}

// DownloadUpgrade
func upgradeWindows(downloadUrl string) {
	filePathWindows := filepath.Base(os.Args[0])
	go func() {
		fmt.Println(cyan.Render("Starting download upgrade from: ", downloadUrl+"\n"))
		for i := 1; i < 15; i++ {
			fmt.Print(".")
			time.Sleep(500 * time.Millisecond)
		}
	}()
	tempFile := filePathWindows + ".tmp" // save version as .tmp file
	err := downloadFile(downloadUrl, tempFile)
	if err != nil {
		fmt.Println("Download new file error.")
		return
	}
	// Rename current executable to .old and rename the new one to current executable
	oldFile := filePathWindows + ".old"
	if updateConfig.RemoveOldFiles {
		os.Remove(oldFile) // Remove old backup if exists
	}
	if err := os.Rename(filePathWindows, oldFile); err != nil {
		fmt.Println("Move to old file was failure.")
		return
	}
	// Create update batch script
	batchContent := `@echo off
timeout /t 2 /nobreak >nul
move /Y "` + tempFile + `" "` + filePathWindows + `"
del "%~f0"
`
	batchFile := "update.bat"
	if err := os.WriteFile(batchFile, []byte(batchContent), 0755); err != nil {
		fmt.Println("Running batch file was failure.")
		return
	}
	// Run the batch file and wait for completion
	fmt.Println("Executing update script...")
	cmd := exec.Command("cmd.exe", "/C", batchFile)
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error executing update script: %v\n", err)
		return
	}
	time.Sleep(3 * time.Second)
	// Verify the new file exists
	if _, err := os.Stat(filePathWindows); os.IsNotExist(err) {
		fmt.Printf("Error: New file %s was not created\n", filePathWindows)
		return
	}
	fmt.Println("Update completed successfully. Old version saved as " + oldFile + ", the old version will be removed automatically when application launch next time.")
	os.Exit(0) // Exit the program after successful update
}
func upgradeLinux(downloadUrl string) {
	filePathLinux := filepath.Base(os.Args[0])

	go func() {
		fmt.Println(cyan.Render("Starting download upgrade from: ", downloadUrl+"\n"))
		for i := 1; i < 15; i++ {
			fmt.Print(".")
			time.Sleep(500 * time.Millisecond)
		}
	}()
	tempFile := filePathLinux + ".tmp"
	err := downloadFile(downloadUrl, tempFile)
	if err != nil {
		fmt.Println("Download new file error.")
		return
	}
	// Rename current executable to .old and rename the new one to current executable
	oldFile := filePathLinux + ".old"
	if updateConfig.RemoveOldFiles {
		os.Remove(oldFile) // Remove old backup if exists
	}
	if err := os.Rename(filePathLinux, oldFile); err != nil {
		// It's okay if the current file doesn't exist (first install)
		fmt.Printf("Move to old file was failure: %v\n", err)
	}

	// Create update shell script
	shContent := "#!/bin/sh\n" +
		"sleep 2\n" +
		"mv \"" + tempFile + "\" \"" + filePathLinux + "\"\n" +
		"chmod +x \"" + filePathLinux + "\"\n" +
		"rm -- \"$0\"\n"

	shFile := "update.sh"
	if err := os.WriteFile(shFile, []byte(shContent), 0755); err != nil {
		fmt.Println("Writing update script was failure.")
		return
	}

	// Run the shell script and wait a moment for it to take over
	fmt.Println("Executing update script...")
	cmd := exec.Command("sh", shFile)
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error executing update script: %v\n", err)
		return
	}
	time.Sleep(3 * time.Second)

	// Verify the new file exists
	if _, err := os.Stat(filePathLinux); os.IsNotExist(err) {
		fmt.Printf("Error: New file %s was not created\n", filePathLinux)
		return
	}

	fmt.Println("Update completed successfully. Old version saved as " + oldFile + ", the old version will be removed automatically when application launch next time.")
	os.Exit(0) // Exit the program after successful update
}

func downloadFile(url, filepath string) error {
	// HTTP GET
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Download new file was failure.")
		return err
	}
	defer resp.Body.Close()
	// Get file size
	size := resp.ContentLength
	// creat progress bar
	bar := pb.Full.Start64(size)
	defer bar.Finish()
	// create temp file first
	//tempFile := filePathWindows + ".tmp"
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Create tmp file was failure.")
		return err
	}
	defer file.Close()
	// create writer with progress bar
	writer := bar.NewProxyWriter(file)
	bar.SetRefreshRate(time.Second)
	// write file
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		fmt.Println("Create new file was failure.")
		return err
	}
	time.Sleep(2 * time.Second)
	return nil
}

func SortVersion(versionStr string) (*Version, error) {
	partsVersion := strings.Split(versionStr, ".")
	//fmt.Printf("leng is : %d", len(partsVersion))
	if len(partsVersion) != 3 {
		return nil, fmt.Errorf("invalid version format: %s (split len: %d), verion number must be 3 parts: x.x.x", versionStr, len(partsVersion))
	}

	majorVersion, err := strconv.Atoi(partsVersion[0])
	if err != nil {
		return nil, err
	}

	minorVersion, err := strconv.Atoi(partsVersion[1])
	if err != nil {
		return nil, err
	}

	patchVersion, err := strconv.Atoi(partsVersion[2])
	if err != nil {
		return nil, err
	}

	return &Version{MajorVersion: majorVersion, MinorVersion: minorVersion, PatchVersion: patchVersion}, nil
}

func CompareVersions(ver1, ver2 *Version) int {
	if ver1.MajorVersion != ver2.MajorVersion {
		return ver1.MajorVersion - ver2.MajorVersion
	}
	if ver1.MinorVersion != ver2.MinorVersion {
		return ver1.MinorVersion - ver2.MinorVersion
	}
	return ver1.PatchVersion - ver2.PatchVersion
}

func VersionCompareResult(currentVersion, newVersion string) string {
	if currentVersion == "" {
		fmt.Println("ERROR: Empty current version.")
		return ""
	}
	if newVersion == "" {
		fmt.Println("ERROR: Empty new version.")
		return ""
	}
	verStr1, err := SortVersion(currentVersion)
	if err != nil {
		fmt.Printf("ERROR: Parsing current version %q: %v\n", currentVersion, err)
		return ""
	}
	VerStr2, err := SortVersion(newVersion)
	if err != nil {
		fmt.Println("Error parsing new version:", err)
		return ""
	}
	result := CompareVersions(verStr1, VerStr2)
	if result > 0 {
		return "lower" //current version is newer than new version
	} else if result < 0 {
		return "newer" //need upgrade
	} else {
		return "equal"
		//equal versions
	}
}
