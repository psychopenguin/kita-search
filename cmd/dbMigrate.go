// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/jinzhu/gorm"
	"github.com/psychopenguin/kita-search/pkg/kita"
	"github.com/spf13/cobra"
)

// dbMigrateCmd represents the dbMigrate command
var dbMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		migrate()
	},
}

func init() {
	dbCmd.AddCommand(dbMigrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dbMigrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dbMigrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func migrate() {
	db, err := gorm.Open("sqlite3", "kita.db")
	if err != nil {
		panic("Failed to open db")
	}
	defer db.Close()
	db.AutoMigrate(&kita.Kita{}, &kita.District{})
	db.Model(&kita.Kita{}).AddForeignKey("district_id", "districts(id)", "RESTRICT", "RESTRICT")
}
