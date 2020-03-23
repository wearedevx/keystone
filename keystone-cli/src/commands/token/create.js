const chalk = require('chalk')
const { cli } = require('cli-ux')
const { newShare } = require('@keystone.sh/core/lib/commands/share')
const { ROLES } = require('@keystone.sh/core/lib/constants')

const { CommandSignedIn } = require('../../lib/commands')

class TokenCommand extends CommandSignedIn {
  async newShare(env) {
    await this.withUserSession(async userSession => {
      cli.action.start('Creating a new read only user')
      const project = await this.getProjectName()

      const { privateKey, membersDescriptor } = await newShare(userSession, {
        project,
        env,
      })
      cli.action.stop('done')

      const data = JSON.stringify({
        project,
        env,
        members: membersDescriptor.content[ROLES.ADMINS],
        privateKey,
      })

      const buff = new Buffer(data)
      const token = buff.toString('base64')
      this.log(`\n▻ ${chalk.yellow(token)}\n`)
      this.log(
        `${`This token must be stored in .env file of your project if you want to use it.\n${chalk.grey(
          `KEYSTONE_SHARED=${chalk.italic('TOKEN')}`
        )}`}`
      )
    })
  }

  async run() {
    const { args } = this.parse(TokenCommand)

    await this.newShare(args.env)
  }
}

TokenCommand.description = `give access to your files with a non-blockstack user

generate a token
the token should be set in the system environment
it allow the user to only run ${chalk.blue(
  '$ ks pull'
)} in order to pull locally files from the selected env
`

TokenCommand.examples = [
  chalk.blue(`$ ks token create ${chalk.italic('ENV_NAME')}`),
]

TokenCommand.args = [
  {
    name: 'env',
    required: true, // make the arg required with `required: true`
    description: `environment you want the token to be created on`, // help description
    hidden: false,
  },
]

module.exports = TokenCommand
