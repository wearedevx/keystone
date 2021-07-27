module github.com/wearedevx/keystone/cli

go 1.16

require (
	github.com/bxcodec/faker/v3 v3.6.0
	github.com/cossacklabs/themis/gothemis v0.13.1
	github.com/eiannone/keyboard v0.0.0-20200508000154-caf4b762e807
	github.com/google/go-github/v32 v32.1.0
	github.com/jamesruan/sodium v0.0.0-20181216154042-9620b83ffeae
	github.com/jedib0t/go-pretty/v6 v6.2.1
	github.com/logrusorgru/aurora/v3 v3.0.0
	github.com/manifoldco/promptui v0.8.0
	github.com/rogpeppe/go-internal v1.8.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.0
	github.com/udhos/equalfile v0.3.0
	github.com/wearedevx/keystone/api v0.0.0-20210727120514-53b74e65257d
	github.com/xanzy/go-gitlab v0.49.0
	golang.org/x/oauth2 v0.0.0-20210427180440-81ed05c6b58c
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/wearedevx/keystone/cmd => ./cmd

replace github.com/wearedevx/keystone/api => ../api
