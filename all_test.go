package go_common

import (
	"gitlab-prod.ptit365.com/go-common/logs"
	"gitlab-prod.ptit365.com/go-common/utils"
	"gitlab-prod.ptit365.com/go-common/workqueue"
	"testing"
	"time"
)

func TestAll(t *testing.T){
	logs.CreateDefaultLogger("tester")
	utils.ToDatetimeStringWithoutDash(time.Now())
	workqueue.New()
}