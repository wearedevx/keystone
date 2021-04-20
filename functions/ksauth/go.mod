module github.com/wearedevx/keystone/functions/ksauth

go 1.15

require (
	filippo.io/age v1.0.0-rc.1
	filippo.io/edwards25519 v1.0.0-beta.3 // indirect
	github.com/GoogleCloudPlatform/cloudsql-proxy v1.18.0
	github.com/GoogleCloudPlatform/functions-framework-go v1.1.0
	github.com/cossacklabs/themis/gothemis v0.13.1 // indirect
	github.com/eiannone/keyboard v0.0.0-20200508000154-caf4b762e807
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-openapi/strfmt v0.19.5 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-github v17.0.0+incompatible // indirect
	github.com/google/go-github/v32 v32.1.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/jedib0t/go-pretty/v6 v6.1.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/logrusorgru/aurora/v3 v3.0.0
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/pelletier/go-toml v1.9.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/wearedevx/keystone v0.0.0-20210412140218-f907250b7cf7
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1 // indirect
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602
	golang.org/x/sys v0.0.0-20210403161142-5e06dd20ab57 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/h2non/gock.v1 v1.0.16 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/postgres v1.0.0
	gorm.io/gorm v1.21.7
)

replace github.com/wearedevx/keystone => ../../
