package audit

import (
	"github.com/gaozhengxin/bridgeaudit/mongodb"
)

var (
	dbAPI mongodb.AccountingAPI
)

func StartAudit() {
	dbAPI = mongodb.NewAccountingAPI()
}
