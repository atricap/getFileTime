package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows"
)

const (
	timeLayout string = "2006-01-02 15:04:05.999999999"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: getFileTime <file>\n")
		return
	}

	filePath := os.Args[1]

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: file '%s' does not exist\n", filePath)
		return
	}

	// Get handle to the file
	handle, err := windows.CreateFile(
		windows.StringToUTF16Ptr(filePath),
		0,
		0,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return
	}
	defer windows.CloseHandle(handle)

	// Get current file times
	var creationTimeFt, lastAccessTimeFt, lastWriteTimeFt windows.Filetime
	err = windows.GetFileTime(handle, &creationTimeFt, &lastAccessTimeFt, &lastWriteTimeFt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting file times: %v\n", err)
		return
	}

	// Convert FILETIME to time.Time
	creationTime := time.Unix(0, creationTimeFt.Nanoseconds())
	lastAccessTime := time.Unix(0, lastAccessTimeFt.Nanoseconds())
	lastWriteTime := time.Unix(0, lastWriteTimeFt.Nanoseconds())

	fmt.Printf("Creation: %s\nAccess: %s\nModification: %s\n",
		creationTime.Format(timeLayout),
		lastAccessTime.Format(timeLayout),
		lastWriteTime.Format(timeLayout),
	)
}
