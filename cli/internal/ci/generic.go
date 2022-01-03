package ci

import (
	"io"
	"log"
	"os"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

const GenericCI CiServiceType = "generic-ci"

type genericCiService struct {
	log                 *log.Logger
	err                 error
	name                string
	ctx                 *core.Context
	lastEnvironmentSent string
}

// GenericCi function returns an instance of genericCiService
func GenericCi(ctx *core.Context, name string) CiService {
	kf := keystonefile.KeystoneFile{}
	kf.Load(ctx.Wd)

	return &genericCiService{
		log:                 log.New(log.Writer(), "[GenericCi] ", 0),
		err:                 nil,
		name:                name,
		ctx:                 ctx,
		lastEnvironmentSent: "",
	}
}

// Name method returns the name of the ci service
func (g *genericCiService) Name() string { return g.name }

// Type method returns the type of ci service
func (g *genericCiService) Type() string { return string(GenericCI) }

// Usage method returns the string explaining how to use keystone with
// this type of service.
func (g *genericCiService) Usage() string {
	return ui.RenderTemplate(
		"generic-ci-usage",
		`They have been packaged in a 'keystone.tar.gz' archive.
To use them in your pipeline, upload that archive to your ci service,
and run the following script: 

  {{ "# extract and remove the archive" | bright_black }}
  tar -zxvf keystone.tar.gz; rm keystone.tar.gz;
  {{ "# source the dotenv file" | bright_black }}
  set -o allexport; source .keystone/cache/$ks_environment/.env; set +o allexport;
  {{ "# move the files in the cache" | bright_black }}
  if [ "$(ls -A .keystone/cache/$ks_environment/files/*)" ]; then
    cp -r .keystone/cache/$ks_environment/files/* ./;
  fi
    `,
		map[string]string{
			"Environment": g.lastEnvironmentSent,
		},
	)
}

// Setup method does nothing for generic service
func (g *genericCiService) Setup() CiService { return g }

// GetOptions method returns an empty map for generic service
func (g *genericCiService) GetOptions() map[string]string { return map[string]string{} }

// PushSecret method will create an archive for the current environment
// that will only contain secrets and files in mentionned ih the keysotne file
func (g *genericCiService) PushSecret(
	message models.MessagePayload,
	environment string,
) CiService {
	archive, err := getArchiveBuffer(g.ctx, environment)
	if err != nil {
		g.err = err
		return g
	}

	file, err := os.OpenFile("keystone.tar.gz", os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		g.err = err
		return g
	}
	defer utils.Close(file)

	_, err = io.Copy(file, archive)
	if err != nil {
		g.err = err
		return g
	}

	return g
}

// CleanSecrets method does nothing for the generic environment
func (g *genericCiService) CleanSecret(environment string) CiService {
	return g
}

// CheckSetup method does nothing for the generic environment
func (g *genericCiService) CheckSetup() CiService {
	return g
}

// Error method returns the last error encountered
func (g *genericCiService) Error() error {
	return g.err
}
