package config

import (
	"fmt"
	"testing"
)

// testAccount["account_type"] = "github"

// var test_account = make(map[account_type:"github" email:"abigael.laldji@protonmail.com"
// ext_id:56883564 fullname:"Michel" private_key:"REC2-!wFGY٩)(Mx_Ni" public_key:"UEC2-*F$hB?0;@֞S$?`R\̴Bqy" user_id:"00fb7666-de43-4559-b4e4-39b172117dd8" username:"LAbigael"])

func TestAddAccount(t *testing.T) {

	testAccount := map[string]string{
		"account_type": "github",
		"email":        "abigael.laldji@protonmail.com",
		"ext_id":       "56883564",
		"fullname":     "Michel",
		"private_key":  "REC2-!wFGY٩)(Mx_Ni",
		"public_key":   "UEC2-*F$hB?0;@֞S$?`R\\̴Bqy",
		"user_id":      "00fb7666-de43-4559-b4e4-39b172117dd8",
		"username":     "LAbigael",
	}

	accountIndex := AddAccount(testAccount)
	fmt.Println(accountIndex)
}
