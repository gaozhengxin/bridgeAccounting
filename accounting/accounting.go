package accounting

import (
	"github.com/gaozhengxin/bridgeAccounting/mongodb"
)

var (
	dbAPI mongodb.AccountingAPI
)

func StartAccounting() {
	dbAPI = mongodb.NewAccountingAPI()
}
