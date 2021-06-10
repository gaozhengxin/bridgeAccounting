package accounting

import (
	"github.com/gaozhengxin/bridgeaudit/mongodb"
)

var (
	dbAPI mongodb.AccountingAPI
)

func StartAccounting() {
	dbAPI = mongodb.NewAccountingAPI()
}
