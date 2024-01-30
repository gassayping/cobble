package new

import (
	"bufio"
	"cobble/project"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

func Run(args []string) {
	answer := make(chan string, 1)
	go getAnswers(answer, args)

	project := project.Project{}
	project.Name = <-answer

	os.MkdirAll(project.Name+"/src/"+project.Name+"_BP/scripts", os.ModePerm)
	os.Chdir(project.Name)

	if exec.Command("git", "init").Run() != nil {
		panic("Failed to initialize git")
	}

	if exec.Command("npm", "init", "-y").Run() != nil {
		panic("Failed to initialize npm")
	}

	if exec.Command("npm", "install", "typescript").Run() != nil {
		panic("Could not install typescript")
	}
	writeTSConfig(project.Name)

	project.Description = <-answer
	project.APIVersion = <-answer

	if exec.Command("npm", "install", "@minecraft/server@"+project.APIVersion).Run() != nil {
		panic("Could not install @minecraft/server@" + project.APIVersion)
	}

	project.IsStable = <-answer == "true"
	project.UsesUI = <-answer == "true"

	if project.UsesUI {
		if project.IsStable {
			if exec.Command("npm", "install", "@minecraft/server-ui").Run() != nil {
				panic("Could not install @minecraft/server-ui")
			}
		} else {
			if exec.Command("npm", "install", "@minecraft/server-ui@1.2.0-beta.1.20.50-stable").Run() != nil {
				panic("Could not install @minecraft/server-ui@1.2.0-beta.1.20.50-stable")
			}
		}
	}
	writeBPManifest(&project)

}

func getAnswers(answer chan string, args []string) {
	reader := bufio.NewReader(os.Stdin)
	var proj project.Project

	if len(args) == 0 {
		fmt.Print("Project Name: ")
		proj.Name, _ = reader.ReadString('\n')
		proj.Name = strings.TrimSpace(proj.Name)
	} else {
		proj.Name = args[0]
	}
	answer <- proj.Name

	fmt.Print("Project Description: ")
	proj.Description, _ = reader.ReadString('\n')
	proj.Description = strings.TrimSpace(proj.Description)

	versions := getAvailableVersions()
	fmt.Println("Available Stable Versions:\n" +
		"\t" + strings.Join(versions.Stables, "\n\t") + "\n\t" + versions.Latest)
	fmt.Print("Version: ")
	proj.APIVersion, _ = reader.ReadString('\n')
	proj.APIVersion = strings.TrimSpace(proj.APIVersion)
	proj.IsStable = len(proj.APIVersion) == 5 || strings.Contains(proj.APIVersion, "-rc")

	answer <- proj.Description
	answer <- proj.APIVersion

	fmt.Print("Use @minecraft/server-ui? (y/n): ")
	var input string
	fmt.Scanln(&input)
	proj.UsesUI = input == "" || input[0] == 'y'

	answer <- fmt.Sprint(proj.IsStable)
	answer <- fmt.Sprint(proj.UsesUI)
}
func writeBPManifest(project *project.Project) {
	os.WriteFile("src/"+project.Name+"_BP/manifest.json", []byte(
		fmt.Sprintf(`{
	"format_version": 2,
	"header": {
		"name": "%v",
		"description": "%v",
		"uuid": "%v",
		"version": [
			0,
			0,
			1
		],
		"min_engine_version": [
			1,
			20,
			50
		]
	},
	"modules": [
		{
			"type": "data",
			"uuid": "%v",
			"version": [
				1,
				0,
				0
			]
		},
		{
			"type": "script",
			"language": "javascript",
			"uuid": "%v",
			"version": [
				0,
				0,
				1
			],
			"entry": "scripts/main.js"
		}
	],
	"dependencies": [
		{
			"module_name": "@minecraft/server",
			"version": "%v"
		},
		{
			"module_name": "@minecraft/server-ui",
			"version": "1.2.0-beta"
		}
	]
}`,
			project.Name, project.Description, uuid.New(), uuid.New(), uuid.New(), project.APIVersion),
	), 0666)
}

func writeTSConfig(projName string) {
	os.WriteFile("tsconfig.json", []byte(
		fmt.Sprintf(`{
	"compileOnSave": true,
	"compilerOptions": {
		"lib": [
			"ESNext",
			"DOM"
		],
		"target": "ESNext",
		"moduleResolution": "node",
		"outDir": "src/%v_BP/scripts",
		"removeComments": true
	}
}`, projName),
	), 0666)
}
