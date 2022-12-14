/*
Copyright © 2022 Isan Rivkin isanrivkin@gmail.com

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
	"fmt"

	"github.com/isan-rivkin/surf/lib/awsu"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	tableNamePattern       string
	ddbQuery               string
	ddbIncludeGlobalTables *bool
	ddbListTables          *bool
	ddbFilterTables        *bool
	ddbFilterKeys          *bool
	ddbFilterData          *bool
	ddbFilterAllOpts       *bool
	ddbStopOnFirstMatch    *bool
)

// ddbCmd represents the ddb command
var ddbCmd = &cobra.Command{
	Use:   "ddb",
	Short: "DynamoDB search tool",
	Long: `
	TBD...

	surf ddb -q val -t table-* [--filter-tables, --filter-keys, --filter-data, --filter-all]

	surf ddb --list-tables
`,
	Run: func(cmd *cobra.Command, args []string) {
		// MARSHAL ATTRIBUTES UTILITY https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbattribute/
		fmt.Println("ddb called")
		auth, err := awsu.NewSessionInput(awsProfile, awsRegion)

		if err != nil {
			log.WithError(err).Fatalf("failed creating session in AWS")
		}

		client, err := awsu.NewDDB(auth)

		if err != nil {
			log.WithError(err).Fatalf("failed creating ddb session")
		}
		ddb := awsu.NewDDBClient(client)
		if *ddbListTables {
			listDDBTables(ddb, true, *ddbIncludeGlobalTables)
			return
		} else {
			err = ddb.ScanTable(tableNamePattern)
			if err != nil {
				log.WithError(err).Fatalf("failed scanning table")
			}
		}
	},
}

func listDDBTables(ddb awsu.DDBApi, withNonGlobal, withGlobal bool) {
	tables, err := ddb.ListCombinedTables(withNonGlobal, withGlobal)
	if err != nil {
		log.WithError(err).Fatalf("failed listing tables")
	}
	for _, t := range tables {
		fmt.Println(t.TableName())
	}
}

func init() {
	rootCmd.AddCommand(ddbCmd)

	ddbCmd.PersistentFlags().StringVarP(&awsProfile, "profile", "p", getDefaultProfileEnvVar(), "~/.aws/credentials chosen account")
	ddbCmd.PersistentFlags().StringVarP(&awsRegion, "region", "r", "", "~/.aws/config default region if empty")
	ddbCmd.PersistentFlags().StringVarP(&ddbQuery, "query", "q", "", "filter query regex supported")
	ddbCmd.PersistentFlags().StringVarP(&tableNamePattern, "table", "t", "", "regex table pattern name to match")

	ddbListTables = ddbCmd.Flags().Bool("list-tables", false, "list all available tables")
	ddbIncludeGlobalTables = ddbCmd.Flags().Bool("include-global-tables", true, "if true will include global tables during search")
	ddbStopOnFirstMatch = ddbCmd.Flags().Bool("stop-first-match", false, "if true stop stop searching on first match found")
}
