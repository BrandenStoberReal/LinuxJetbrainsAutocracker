package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func crackProgram(name string, binpath string, filter string, misc1 string, misc2 string) {
	// Check if bin folder exists
	_, err := os.Stat(binpath)
	if err != nil {
		log.Fatalln("No bin folder found. Is this a Jetbrains program? Skipping.")
	}

	// Update vmoptions
	_, err = os.Stat(binpath + "/" + "jetbrains_client64.vmoptions")
	if err != nil {
		log.Fatalln("Jetbrains .vmoptions not found. Is this a real Jetbrains product? Skipping.")
	}

	// Open the jetbrains jetbrainsFile
	jetbrainsFile, err := os.Open(binpath + "/" + "jetbrains_client64.vmoptions")
	if err != nil {
		log.Fatalf("Error opening Jetbrains client .vmoptions file: %s", err.Error())
	}
	defer jetbrainsFile.Close()

	// Read options data
	contents, err := io.ReadAll(jetbrainsFile)
	if err != nil {
		log.Fatalf("Error reading Jetbrains client .vmoptions file: %s", err.Error())
	}

	// Convert to string
	contentsStr := string(contents)

	// Check if we already cracked this file
	if strings.Contains(contentsStr, filter) {
		log.Println("Jetbrains client .vmoptions file already cracked. Skipping.")
	} else {
		// Inject cracked agent
		contentsStr = contentsStr + "\n" + filter + "\n" + misc1 + "\n" + misc2

		// Flush back to jetbrainsFile
		_, err = jetbrainsFile.Write([]byte(contentsStr))
		if err != nil {
			log.Fatalf("Error writing Jetbrains client .vmoptions file: %s", err.Error())
		} else {
			log.Println("Jetbrains client .vmoptions file successfully cracked.")
		}
	}

	// Begin cracking the program-specific file
	_, err = os.Stat(binpath + "/" + name + "64.vmoptions")
	if err != nil {
		log.Fatalf("Error finding program-specific .vmoptions file: %s", err.Error())
	}

	// Open program vmoptions file
	programfile, err := os.Open(binpath + "/" + name + "64.vmoptions")
	if err != nil {
		log.Fatalln("Error opening program .vmoptions jetbrainsFile.")
	}

	defer programfile.Close()

	// Read all data from program vmoptions
	contents, err = io.ReadAll(programfile)
	if err != nil {
		log.Fatalf("Error reading program-specific .vmoptions file: %s", err.Error())
	}

	// Check if we cracked already and crack if not.
	contentsStr = string(contents)
	if strings.Contains(contentsStr, filter) {
		log.Println("Program-specific .vmoptions file already cracked. Skipping.")
	} else {
		contentsStr = contentsStr + "\n" + filter + "\n" + misc1 + "\n" + misc2
		_, err = programfile.Write([]byte(contentsStr))
		if err != nil {
			log.Fatalf("Error writing program-specific .vmoptions file: %s", err.Error())
		} else {
			log.Println("Program .vmoptions file cracked.")
		}
	}
}

func main() {
	pathToCrack := flag.String("agent", "", "Path to a cracked netfilter agent")
	pathToProgram := flag.String("program", "", "Path to your Jetbrains program container directory. Defaults to the Linux toolbox install location. All jetbrains products in this directory will be cracked.")

	flag.Parse()

	netfilterPath := *pathToCrack
	programPath := *pathToProgram

	// Required argument
	if netfilterPath == "" {
		log.Fatal("No netfilter agent found! Please specify the path to a cracked agent JAR file.")
	}

	// Required arguments for crack to function
	netfilterJava := "-javaagent:" + netfilterPath + "=jetbrains"
	miscJava1 := "--add-opens=java.base/jdk.internal.org.objectweb.asm=ALL-UNNAMED"
	miscJava2 := "--add-opens=java.base/jdk.internal.org.objectweb.asm.tree=ALL-UNNAMED"

	// Determine toolbox path
	if programPath == "" {
		// Fetch home directory
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Your home dir is: %s!", dirname)
		programPath, err = filepath.Abs(dirname + "/.local/share/JetBrains/Toolbox/apps/")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Search in default parrot directory
	_, err := os.Stat(programPath)
	if err != nil {
		// Default folder does not exist
		log.Println("Jetbrains Toolbox install directory not found! Please set manually from the CLI.")
		os.Exit(1)
	} else {
		// Default folder exists
		entries, err := os.ReadDir("/home/amnesia/.local/share/JetBrains/Toolbox/apps/")
		if err != nil {
			log.Fatalln("Error listing Jetbrains products.")
		}

		// Enum through products
		for _, entry := range entries {
			log.Printf("Found program to crack: %s", entry.Name())
			binFolder := "/home/amnesia/.local/share/JetBrains/Toolbox/apps/" + entry.Name() + "/bin"

			crackProgram(entry.Name(), binFolder, netfilterJava, miscJava1, miscJava2)
		}
	}
}
