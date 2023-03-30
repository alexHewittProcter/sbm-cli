/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Searches Node.js projects",
	Long: `This command searches Node.js projects for a specific package provided by the user.`,
	Run: func(cmd *cobra.Command, args []string) {	
		partialSearch, err := cmd.Flags().GetBool("partial")
		if err != nil {
			fmt.Println("Error:", err)
		}

		var packageToSearch string
		if len(args) >= 1 && args[0] != "" {
			packageToSearch = args[0]
		} else {
			fmt.Println("Please provide a package to search")
			return
		}
	
		packageJSON,err := getPackageJsons(".")
		if err != nil {
			fmt.Println("Error:", err)
		}

		for i := 0; i < len(packageJSON); i++ {
			jsonPath := packageJSON[i]
			deps, name, err := getDependencies(jsonPath)
			if err != nil {
				fmt.Println("Error:", err)
			}

			for key,val := range deps {
				// fmt.Println(key)
				if(key == packageToSearch || partialSearch && strings.Contains(key,packageToSearch)) {
					printFoundPackage(name,strings.Replace(jsonPath,"../","",-1), key,val)
				}
			}
		}
		
	},
}

func printFoundPackage(name string, path string, foundPackage string, packageVersion string) {
	v := fmt.Sprintf("%s - %s:%s - %s", name, foundPackage, packageVersion, path)
	fmt.Println(v)
}

func getPackageJsons(path string) ([]string, error) {
	var packages []string
	err:= filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "node_modules" {
			return filepath.SkipDir
		}

		if info.Name() == "package.json" {
			packages = append(packages, currentPath)
		}

		return nil
	})
	

	return packages,err
}

type NodeJson struct {
    Dependencies   map[string]string `json:"dependencies"`
    DevDependencies   map[string]string `json:"devDependencies"`
	Name string `json:"name"`
}


func getDependencies(path string) (map[string]string, string, error) {
	file,err := os.ReadFile(path)
	if err != nil {
		return nil,"",err
	}

	var data NodeJson 

	err = json.Unmarshal([]byte(file),&data)
	if err != nil {
		return nil,"",err
	}

	dependencies := data.Dependencies
	devDependencies := data.DevDependencies

	if dependencies == nil{
		if devDependencies == nil {
			 return map[string]string{},"",nil
		} else {
			return devDependencies,"",nil
		}
	}

	for k, v := range devDependencies {
		dependencies[k] = v
	}

	return dependencies,data.Name, nil
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().BoolP("partial", "p", false, "Help message for toggle")
}
