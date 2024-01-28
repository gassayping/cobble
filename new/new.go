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
	reader := bufio.NewReader(os.Stdin)
	var proj project.Project
	if len(args) == 0 {
		fmt.Print("Project Name: ")
		proj.Name, _ = reader.ReadString('\n')
	} else {
		proj.Name = args[0]
	}

	fmt.Print("Project Description: ")
	proj.Description, _ = reader.ReadString('\n')

	versions := getAvailableVersions()
	fmt.Println("Available Stable Versions:")
	fmt.Println("\t" + strings.Join(versions.Stables, "\n\t") + "\n\t" + versions.Latest)
	fmt.Print("Version: ")
	proj.APIVersion, _ = reader.ReadString('\n')
	proj.IsStable = len(proj.APIVersion) == 5 || strings.Contains(proj.APIVersion, "-rc")

	fmt.Print("Use @minecraft/server-ui? (y/n): ")
	var input string
	fmt.Scanln(&input)
	proj.UsesUI = input == "" || input[0] == 'y'

	initialize(&proj)

}

func initialize(project *project.Project) {
	os.MkdirAll(project.Name+"/src/"+project.Name+"_BP/scripts", os.ModePerm)
	writeBPManifest(project)
	os.Chdir(project.Name)
	if exec.Command("npm", "init", "-y").Run() != nil {
		fmt.Println("Could not run npm init")
	}
	if exec.Command("npm", "install", "@minecraft/server@"+project.APIVersion).Run() != nil {
		fmt.Println("Could not install @minecraft/server@" + project.APIVersion)
		return
	}
	if project.UsesUI {
		if project.IsStable {
			if exec.Command("npm", "install", "@minecraft/server-ui").Run() != nil {
				fmt.Println("Could not install @minecraft/server-ui")
				return
			}
		} else {
			if exec.Command("npm", "install", "@minecraft/server-ui@1.2.0-beta").Run() != nil {
				fmt.Println("Could not install @minecraft/server-ui@1.2.0-beta")
				return
			}
		}
	}

	if exec.Command("npm", "install", "typescript").Run() != nil {
		fmt.Println("Could not install typescript")
		return
	}
	writeTSConfig(project.Name)

}

func writeBPManifest(project *project.Project) {
	os.WriteFile(project.Name+"/src/"+project.Name+"_BP/manifest.json", []byte(
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
