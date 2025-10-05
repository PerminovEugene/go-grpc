package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// Create new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Add watch directories
	root := "./internal"
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip hidden, vendor, tmp, bin directories
			if strings.HasPrefix(info.Name(), ".") || info.Name() == "vendor" ||
				info.Name() == "tmp" || info.Name() == "bin" {
				return filepath.SkipDir
			}
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Run tests initially
	runTests()

	// Watch for changes
	debounce := time.NewTimer(0)
	<-debounce.C // drain the timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Only react to Go file writes
			if strings.HasSuffix(event.Name, ".go") && (event.Op&fsnotify.Write == fsnotify.Write) {
				debounce.Reset(100 * time.Millisecond)
			}

		case <-debounce.C:
			fmt.Println("\n🔄 Files changed, running tests...")
			runTests()

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}

func runTests() {
	cmd := exec.Command("go", "test", "./...", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	fmt.Println(strings.Repeat("-", 60))
	if err != nil {
		fmt.Printf("❌ Tests failed in %v\n", duration)
	} else {
		fmt.Printf("✅ All tests passed in %v\n", duration)
	}
	fmt.Println("👀 Watching for changes... (Ctrl+C to stop)")
	fmt.Println(strings.Repeat("-", 60))
}
