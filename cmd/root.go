/*
Copyright © 2020 Ross Lonstein <ross@develop.lonsteins.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
	log "github.com/sirupsen/logrus"
)

// Verbose toggle debugging noise
var Verbose bool

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "yt",
	Short: "yt, a YAML Tool",
	Long: `yt, YAML Tool, a quick hack for extracting data using JSONPath`,
	Args: cobra.RangeArgs(1,2),
	Run: yt,
}

// the actual work
func yt(cmd *cobra.Command, args []string) {
	if Verbose {
		log.SetLevel(log.DebugLevel)
	}
	var jp string
	var data []byte
	var err error
	if len(args) == 1 {
		jp = args[0]
		log.Debug("Will read from stdin...")
		data, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal("Read from stdin failed: ", err)
		}
	} else {
		jp = args[1]
		fn := args[0]
		log.Debug("Will read from '", fn, "'...")
		data, err = ioutil.ReadFile(fn)
		if err != nil {
			log.Fatal("Read failed: ", err)
		}
	}
	log.Debug("Read ", len(data), " bytes")
	var doc yaml.Node
	log.Debug("Unmarshaling...")
	if err := yaml.Unmarshal(data, &doc); err != nil {
		log.Fatal("Unmarshaling failed: ", err)
	}
	log.Debug("Constructing JSON path...")
	path, err := yamlpath.NewPath(jp)
	if err != nil {
		log.Fatal(err)
	}
	results, err := path.Find(&doc)
	if err != nil {
		log.Fatal("JSONPath error:", err)
	}
	log.Debug("Found ", len(results), " matches")
	for _, data := range results {
		out, err := encode(data)
		if err != nil {
			log.Error("Encoding error: ", err)
		} else {
			fmt.Print(out)
		}
	}
}

// Execute invoked by main to run the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Println(`Usage:
    yt [-v] [<path/to/yaml/doc>] <JSONPATH expression>

Options:
  -v       verbose, enable debug messages
  -h       help

Examples:
  $ yt foo.yml '$.bar'
  $ yt '$.bar' < foo.yml | yt '*.baz'

Additional reference:
  https://goessner.net/articles/JsonPath/
  https://golang.org/pkg/regexp/
`)
		return nil
	})
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

func encode(a *yaml.Node) (string, error) {
	var buf bytes.Buffer
	e := yaml.NewEncoder(&buf)

	defer func() {
		// can't do anything about it, but possible
		err := e.Close()
		if err != nil {
			log.Info("Error closing encoder: ", err)
		}
	}()
	
	e.SetIndent(2)

	if err := e.Encode(a); err != nil {
		return "", err
	}

	return buf.String(), nil
}
