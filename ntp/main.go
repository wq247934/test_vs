package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
	"time"
)

func main() {
	syncChronyHandler()
}

func syncChronyHandler() {
	cmd := exec.Command("chronyc", "makestep")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("SyncChronyFailed", stderr.String())
		return
	}
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeout := time.After(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			lastSyncTime, err := fetchChronyLastSyncTime()
			if err != nil {
				fmt.Println(err)
				return
			}
			if !lastSyncTime.Equal(time.Unix(0, 0)) {
				fmt.Println("last_sync_time:", lastSyncTime.Add(time.Hour*8).Format("2006-01-02 15:04:05"))
				return
			}
		case <-timeout:
			fmt.Println("Timeout", "Failed to fetch time within 5 seconds")
			return
		}
	}
}

func fetchChronyLastSyncTime() (time.Time, error) {
	cmd := exec.Command("chronyc", "tracking")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
		return time.Time{}, err
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Ref time (UTC)") {
			parts := strings.Split(line, ":")
			timeString := strings.Join(parts[1:], ":")
			t, err := time.Parse(" Mon Jan _2 15:04:05 2006", timeString)
			if err != nil {
				return time.Time{}, err
			}
			return t, nil
		}
	}
	return time.Time{}, errors.New("LastSyncTimeNotFound")
}
