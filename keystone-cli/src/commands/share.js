const chalk = require('chalk')
const { cli } = require('cli-ux')
const { newShare } = require('@keystone.sh/core/lib/commands/share')
const { ROLES } = require('@keystone.sh/core/lib/constants')

const { CommandSignedIn } = require('../lib/commands')

class ShareCommand extends CommandSignedIn {
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
    const { args } = this.parse(ShareCommand)

    await this.newShare(args.env)
  }
}

ShareCommand.description = `Share your files with a non-blockstack user

Generate a token. 
The token should be set in the process environment of any user. 
This user will be able to run only ${chalk.yellow(
  '$ ks pull'
)} in order to pull locally files from the selected env.
`

ShareCommand.examples = [chalk.yellow(`$ ks share ${chalk.italic('env_name')}`)]

ShareCommand.args = [
  {
    name: 'env',
    required: true, // make the arg required with `required: true`
    description: `Environment you want the user to be created on.`, // help description
    hidden: false,
  },
]

module.exports = ShareCommand
