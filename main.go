package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

func crackProgram(name string, binpath string, filter string, misc1 string, misc2 string) {
	// Check if bin folder exists
	_, err := os.Stat(binpath)
	if err != nil {
		log.Fatalln("Error finding Jetbrains binaries.")
	}

	// Update vmoptions
	_, err = os.Stat(binpath + "/" + "jetbrains_client64.vmoptions")
	if err != nil {
		log.Fatalln("Error finding Jetbrains client vmoptions file.")
	}

	// Open the jetbrains file
	file, err := os.Open(binpath + "/" + "jetbrains_client64.vmoptions")
	if err != nil {
		log.Fatalln("Error opening Jetbrains client vmoptions file.")
	}
	defer file.Close()

	// Read options data
	contents, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln("Error reading Jetbrains client vmoptions file.")
	}

	// Convert to string
	contentsStr := string(contents)

	// Check if we already cracked this file
	if strings.Contains(contentsStr, filter) {
		log.Println("Jetbrains client vmoptions file already cracked. Skipping.")
	} else {
		// Inject cracked agent
		contentsStr = contentsStr + "\n" + filter + "\n" + misc1 + "\n" + misc2

		// Flush back to file
		_, err = file.Write([]byte(contentsStr))
		if err != nil {
			log.Fatalln("Error writing Jetbrains client vmoptions file.")
		} else {
			log.Println("Jetbrains client vmoptions file successfully cracked.")
		}
	}

	// Begin cracking the program-specific file
	_, err = os.Stat(binpath + "/" + name + "64.vmoptions")
	if err != nil {
		log.Fatalln("Error finding program vmoptions file.")
	}

	// Open program vmoptions file
	file, err = os.Open(binpath + "/" + name + "64.vmoptions")
	if err != nil {
		log.Fatalln("Error opening program vmoptions file.")
	}

	defer file.Close()

	// Read all data from program vmoptions
	contents, err = io.ReadAll(file)
	if err != nil {
		log.Fatalln("Error reading program vmoptions file.")
	}

	// Check if we cracked already and crack if not.
	contentsStr = string(contents)
	if strings.Contains(contentsStr, filter) {
		log.Println("Program vmoptions file already cracked. Skipping.")
	} else {
		contentsStr = contentsStr + "\n" + filter + "\n" + misc1 + "\n" + misc2
		_, err = file.Write([]byte(contentsStr))
		if err != nil {
			log.Fatalln("Error writing program vmoptions file.")
		} else {
			log.Println("Program vmoptions file cracked.")
		}
	}
}

func main() {
	pathToCrack := flag.String("agent", "", "Path to a cracked netfilter agent")
	pathToProgram := flag.String("program", "", "Path to Jetbrains program. Defaults to all.")

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

	if programPath == "" {
		// Search in default parrot directory
		_, err := os.Stat("/home/amnesia/.local/share/JetBrains/Toolbox/apps/")
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
	} else {
		if strings.Contains(programPath, "bin") {
			fileinfo, err := os.Stat(programPath + "..")
			if err != nil {
				log.Fatalln("Error finding Jetbrains program folder.")
			}
			crackProgram(fileinfo.Name(), programPath, netfilterJava, miscJava1, miscJava2)
		} else {
			// Assume base root
			fileinfo, err := os.Stat(programPath)
			if err != nil {
				log.Fatalln("Error finding Jetbrains program folder.")
			}
			crackProgram(fileinfo.Name(), programPath+"/bin", netfilterJava, miscJava1, miscJava2)
		}
	}
}
