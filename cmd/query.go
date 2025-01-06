// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/account"
	"github.com/dingxin-tech/mc-cli/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(queryCmd)
}

var queryCmd = &cobra.Command{
	Use:   "query [sql]",
	Short: "Execute query on MaxCompute",
	Run: func(cmd *cobra.Command, args []string) {
		odpsIns := getOdpsIns()
		if len(args) == 0 {
			fmt.Println("No query found")
			return
		}

		sql := strings.Join(args, " ")
		sql = ensureSemicolon(sql)

		hints := map[string]string{
			"odps.sql.select.output.format": "HumanReadable",
			"odps.sql.submit.mode":          "script",
		}
		ins, err := odpsIns.ExecSQl(sql, hints)
		if err != nil {
			fmt.Printf("Submit SQL Error: %v\n", err)
			return
		}

		// 启动 goroutine 获取 logView
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			logView, err := odpsIns.LogView().GenerateLogView(ins, 72)
			if err == nil {
				fmt.Println("Logview:\n", logView)
			} else {
				fmt.Printf("Generate LogView Error: %v\n", err)
			}
		}()

		err = ins.WaitForSuccess()
		if err != nil {
			fmt.Printf("Execute SQL Error: %v\n", err)
			return
		}
		result, err := ins.GetResult()
		if err != nil {
			fmt.Printf("Get Result Error: %v\n", err)
			return
		}
		fmt.Printf("\nResult:\n")
		for _, res := range result {
			fmt.Printf(res.Content())
		}
		wg.Wait()
	},
}

func ensureSemicolon(sql string) string {
	// 去除 SQL 字符串两端的空白字符
	trimmedSQL := strings.TrimSpace(sql)

	// 检查 trimmedSQL 是否以分号结尾
	if !strings.HasSuffix(trimmedSQL, ";") {
		// 如果没有，以分号结尾
		trimmedSQL += ";"
	}

	return trimmedSQL
}

func getOdpsIns() *odps.Odps {
	aliyunAccount := account.NewAliyunAccount(viper.GetString(common.AccessId), viper.GetString(common.AccessKey))
	odpsIns := odps.NewOdps(aliyunAccount, viper.GetString(common.Endpoint))
	odpsIns.SetDefaultProjectName(viper.GetString(common.ProjectName))
	return odpsIns
}
