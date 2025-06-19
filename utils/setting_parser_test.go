package utils_test

import (
	"strings"
	"testing"

	"github.com/dingxin-tech/mc-cli/utils"
)

func TestParse(t *testing.T) {
	sql := "set odps.instance.priority=1; set odps.sql.timezone=1; select 1"
	res := utils.Parse(sql)

	for key := range res.Settings {
		value := res.Settings[key]
		t.Log(key, value)
	}

	if len(res.Settings) != 2 {
		t.Error("settings length is not 2")
	}

	t.Log(res.RemainingQuery)
	if res.RemainingQuery != "  select 1;" {
		t.Error("remaining query is not select 1")
	}
}

func TestParse_Case1(t *testing.T) {
	sql := "SET odps.namespace.schema=true;\nSET LABEL 1 TO TABLE default.wrk_gh_events(`repo_id`, `repo_name`, `org_id`, `org_login`);"
	res := utils.Parse(sql)

	for key := range res.Settings {
		value := res.Settings[key]
		t.Log(key, value)
	}

	if len(res.Settings) != 1 {
		t.Error("settings length is not 1")
	}

	t.Log(res.RemainingQuery)
	if strings.Trim(res.RemainingQuery, "\n ") != "SET LABEL 1 TO TABLE default.wrk_gh_events(`repo_id`, `repo_name`, `org_id`, `org_login`);" {
		t.Error("remaining query is not right")
	}
}
