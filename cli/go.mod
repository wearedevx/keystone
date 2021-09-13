module github.com/wearedevx/keystone/cli

go 1.16

require (
	cloud.google.com/go v0.89.0 // indirect
	github.com/briandowns/spinner v1.16.0
	github.com/bxcodec/faker/v3 v3.6.0
	github.com/cossacklabs/themis/gothemis v0.13.1
	github.com/danwakefield/fnmatch v0.0.0-20160403171240-cbb64ac3d964 // indirect
	github.com/denormal/go-gitignore v0.0.0-20180930084346-ae8ad1d07817 // indirect
	github.com/eiannone/keyboard v0.0.0-20200508000154-caf4b762e807
	github.com/google/go-github/v32 v32.1.0
	github.com/jamesruan/sodium v0.0.0-20181216154042-9620b83ffeae
	github.com/jedib0t/go-pretty/v6 v6.2.1
	github.com/lib/pq v1.10.2 // indirect
	github.com/logrusorgru/aurora/v3 v3.0.0
	github.com/manifoldco/promptui v0.8.0
	github.com/rogpeppe/go-internal v1.8.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.0
	github.com/udhos/equalfile v0.3.0
	github.com/wearedevx/keystone/api v0.0.0-20210910095952-c0e2ba65920c
	github.com/xanzy/go-gitlab v0.49.0
	go.uber.org/zap v1.18.1 // indirect
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	google.golang.org/api v0.52.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/wearedevx/keystone/cmd => ./cmd

replace github.com/wearedevx/keystone/api => ../api
