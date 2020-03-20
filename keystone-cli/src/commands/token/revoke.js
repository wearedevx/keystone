const chalk = require('chalk')
const { cli } = require('cli-ux')
const { newShare } = require('@keystone.sh/core/lib/commands/share')

const { CommandSignedIn } = require('../../lib/commands')

class TokenCommand extends CommandSignedIn {
  async newShare(env) {
    await this.withUserSession(async userSession => {
      cli.action.start('Creating a new read only user')
      const project = await this.getProjectName()

      await newShare(userSession, {
        project,
        env,
      })
      cli.action.stop('done')

      this.log(
        'The token has ben revoked from the environment. It cannot be used to fetch files anymore.'
      )
    })
  }

  async run() {
    const { args } = this.parse(TokenCommand)

    await this.newShare(args.env)
  }
}

TokenCommand.description = `revoke access to your files with a non-blockstack user `

TokenCommand.examples = [
  chalk.blue(`$ ks token revoke ${chalk.italic('ENV_NAME')}`),
]

TokenCommand.args = [
  {
    name: 'env',
    required: true, // make the arg required with `required: true`
    description: `environment you want the token to be revoked on`, // help description
    hidden: false,
  },
]

module.exports = TokenCommand
