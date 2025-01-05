package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type ProcessedData struct {
	depotID        string
	manifestNumber string
	decryptionKey  string
}

func main() {
if len(os.Args) < 2 {
        fmt.Println("Please drag and drop files onto the program")
        fmt.Println("Supported files: .manifest and config.vdf\n")
        fmt.Println("Press Enter to exit...")
        var input string
        fmt.Scanln(&input)
    }
    
    manifestFiles := make(map[string]string)
    configData := make(map[string]ProcessedData)
    
    for _, filePath := range os.Args[1:] {
        fileName := filepath.Base(filePath)
        
        switch {
        case strings.HasSuffix(fileName, ".manifest"):
            processManifestFile(filePath, manifestFiles)
        case fileName == "config.vdf":
            processConfigFile(filePath, configData)
        default:
            fmt.Printf("Unsupported file: %s\n", fileName)
        }
    }
    
    generateOutputFile(manifestFiles, configData)
}

func processManifestFile(filePath string, manifestFiles map[string]string) {
	fileName := filepath.Base(filePath)
	parts := strings.Split(fileName, "_")
	if len(parts) < 2 {
		log.Printf("Invalid manifest filename format: %s\n", fileName)
		return
	}

	depotID := parts[0]
	manifestNumber := strings.TrimSuffix(parts[1], ".manifest")
	
	manifestFiles[depotID] = manifestNumber
	fmt.Printf("Processed Manifest: Depot ID = %s, Manifest Number = %s\n", 
		depotID, manifestNumber)
}

func processConfigFile(filePath string, configData map[string]ProcessedData) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading config file: %v\n", err)
		return
	}

	depotRegex := regexp.MustCompile(`"(\d+)"\s*{\s*"DecryptionKey"\s*"([^"]+)"`)
	matches := depotRegex.FindAllSubmatch(content, -1)

	if len(matches) == 0 {
		log.Println("Could not find any depot IDs and decryption keys")
		return
	}

	for _, match := range matches {
		if len(match) >= 3 {
			depotID := string(match[1])
			decryptionKey := string(match[2])
			
			configData[depotID] = ProcessedData{
				depotID:        depotID,
				decryptionKey:  decryptionKey,
			}

			fmt.Printf("Processed Config: Depot ID = %s, Decryption Key = %s\n", 
				depotID, decryptionKey)
		}
	}
}

func generateOutputFile(manifestFiles map[string]string, configData map[string]ProcessedData) {
	var outputEntries []string

	for depotID, manifestNumber := range manifestFiles {
		if configEntry, exists := configData[depotID]; exists {
			outputLine := fmt.Sprintf(
				"addappid(%s, 1, \"%s\")\n"+
				"setManifestid(%s, \"%s\", 0)", 
				depotID, configEntry.decryptionKey, 
				depotID, manifestNumber,
			)
			outputEntries = append(outputEntries, outputLine)
		}
	}

	if len(outputEntries) == 0 {
		fmt.Println("No matching depot IDs found between manifest and config files")
		return
	}

	sort.Slice(outputEntries, func(i, j int) bool {
		idI, _ := strconv.Atoi(regexp.MustCompile(`\((\d+),`).FindStringSubmatch(outputEntries[i])[1])
		idJ, _ := strconv.Atoi(regexp.MustCompile(`\((\d+),`).FindStringSubmatch(outputEntries[j])[1])
		return idI < idJ
	})

	var appID string
	fmt.Print("What's the game's APP ID? ")
	fmt.Scanln(&appID)

	outputContent := "-- manifest & lua provided by: https://www.piracybound.com/discord\n" +
					"-- via manilua\n" +
					fmt.Sprintf("addappid(%s)\n", appID) +
					strings.Join(outputEntries, "\n")

	outputFilename := fmt.Sprintf("%s.lua", appID)
	err := ioutil.WriteFile(outputFilename, []byte(outputContent), 0644)
	if err != nil {
		log.Printf("Error writing output file: %v\n", err)
		return
	}

	fmt.Printf("Output file generated: %s\n", outputFilename)
	fmt.Println(outputContent)
}