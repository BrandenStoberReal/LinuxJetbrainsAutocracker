package main

import (
	"flag"
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
		log.Println("Jetbrains .vmoptions not found! Is this a real Jetbrains product? Skipping.")
	} else {
		// Open the jetbrains file
		jetbrainsFileBuffer, err := os.ReadFile(binpath + "/" + "jetbrains_client64.vmoptions")
		if err != nil {
			log.Fatalf("Error opening Jetbrains client .vmoptions file: %s", err.Error())
		}

		// Convert to string
		jetbrainsContentStr := string(jetbrainsFileBuffer)

		// Check if we already cracked this file
		if strings.Contains(jetbrainsContentStr, filter) {
			log.Println("Jetbrains client .vmoptions file already cracked. Skipping.")
		} else {
			// Inject cracked agent
			jetbrainsContentStr = jetbrainsContentStr + "\n" + filter + "\n" + misc1 + "\n" + misc2

			// Flush back to jetbrainsFile
			err = os.WriteFile(binpath+"/"+"jetbrains_client64.vmoptions", []byte(jetbrainsContentStr), 0644)
			if err != nil {
				log.Fatalf("Error writing Jetbrains client .vmoptions file: %s", err.Error())
			} else {
				log.Println("Jetbrains client .vmoptions file successfully cracked.")
			}
		}
	}

	// Begin cracking the program-specific file
	_, err = os.Stat(binpath + "/" + name + "64.vmoptions")
	if err != nil {
		log.Fatalf("Error finding program-specific .vmoptions file: %s", err.Error())
	} else {
		// Open program vmoptions file
		programFileBuffer, err := os.ReadFile(binpath + "/" + name + "64.vmoptions")
		if err != nil {
			log.Fatalln("Error opening program .vmoptions jetbrainsFile.")
		}

		// Check if we cracked already and crack if not.
		programContentsStr := string(programFileBuffer)

		if strings.Contains(programContentsStr, filter) {
			log.Println("Program-specific .vmoptions file already cracked. Skipping.")
		} else {
			programContentsStr = programContentsStr + "\n" + filter + "\n" + misc1 + "\n" + misc2
			err = os.WriteFile(binpath+"/"+name+"64.vmoptions", []byte(programContentsStr), 0644)
			if err != nil {
				log.Fatalf("Error writing program-specific .vmoptions file: %s", err.Error())
			} else {
				log.Println("Program .vmoptions file cracked.")
			}
		}
	}
}

func main() {
	pathToCrack := flag.String("agent", "", "Path to the JAR file of a cracked netfilter agent.")
	pathToProgram := flag.String("installDir", "", "Path to your Jetbrains program container directory. Defaults to the Linux toolbox install location. All jetbrains products in this directory will be cracked.")

	flag.Parse()

	// Resolve pointers
	netfilterPath := *pathToCrack
	programPath := *pathToProgram

	// Required argument check
	if netfilterPath == "" {
		log.Fatal("No netfilter agent provided! Please specify the path to a cracked agent JAR file.")
	}

	// Check for netfilter
	netfilterInfo, err := os.Stat(netfilterPath)
	if err != nil {
		log.Fatalf("Error finding netfilter path: %s", err.Error())
	}

	// Handle inference of jar file
	if netfilterInfo.IsDir() {
		oldnetPath := netfilterPath
		netfilterPath, err = filepath.Abs(netfilterPath + "/ja-netfilter.jar")
		if err != nil {
			log.Fatalf("Error resolving absolute path for netfilter: %s", err.Error())
		}

		// Ensure netfilter exists
		netfilterInfo, err = os.Stat(netfilterPath)
		if err != nil {
			log.Fatalf("Could not infer the netfilter jar from the provided folder. Does it exist?")
		} else {
			log.Printf("Found netfilter jar in %s", oldnetPath)
		}
	}

	// Required strings for crack to function
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
		} else {
			log.Println("Jetbrains Toolbox apps directory found! Beginning cracking...")
		}
	}

	// Search in default Linux toolbox dir
	_, err = os.Stat(programPath)
	if err != nil {
		// Default folder does not exist
		log.Fatalln("Jetbrains Toolbox install directory not found! Please set manually from the CLI.")
	} else {
		// Default folder exists
		entries, err := os.ReadDir("/home/amnesia/.local/share/JetBrains/Toolbox/apps/")
		if err != nil {
			log.Fatalf("Error listing Jetbrains products: %s", err.Error())
		}

		log.Println("-----------------------------------------------------------------------------")
		// Enum through products
		for _, entry := range entries {
			if entry.IsDir() {
				log.Printf("Found program to crack: %s", entry.Name())
				binFolder := "/home/amnesia/.local/share/JetBrains/Toolbox/apps/" + entry.Name() + "/bin"

				crackProgram(entry.Name(), binFolder, netfilterJava, miscJava1, miscJava2)
				log.Println("-----------------------------------------------------------------------------")
			}
		}
	}
}
