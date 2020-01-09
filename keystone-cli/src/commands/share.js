const chalk = require('chalk')
const fs = require('fs')
const { newShare, pullShared } = require('@keystone.sh/core/lib/commands/share')
const { flags } = require('@oclif/command')
const { SHARE_FILENAME, ROLES } = require('@keystone.sh/core/lib/constants')

const { CommandSignedIn } = require('../lib/commands')

class ShareCommand extends CommandSignedIn {
  async newShare(action, env) {
    await this.withUserSession(async userSession => {
      const project = await this.getProjectName()
      if (action === 'new') {
        const { privateKey, membersDescriptor } = await newShare(userSession, {
          project,
          env,
        })

        console.log(userSession.getFile)

        const data = JSON.stringify({
          project,
          env,
          members: membersDescriptor.content[ROLES.ADMINS],
          privateKey,
          userSession,
        })
        console.log(data)

        const buff = new Buffer(data)
        const token = buff.toString('base64')
        this.log(`\n▻ ${chalk.yellow(token)}\n`)
        this.log(
          `${`This token must be stored in .env file of your project if you want to use it.\n${chalk.grey(
            `KEYSTONE_SHARED=${chalk.italic('TOKEN')}`
          )}`}`
        )
      }
    })
  }

  async pull(pathToConfig) {
    await this.withUserSession(async userSession => {
      const { project, env, members, privateKey } = JSON.parse(
        fs.readFileSync(pathToConfig)
      )
      userSession.sharedPrivateKey = privateKey
      const absoluteProjectPath = await this.getConfigFolderPath()

      await pullShared(userSession, {
        project,
        env,
        origins: members,
        privateKey,
        absoluteProjectPath,
      })
    })
  }

  async run() {
    const { args, flags } = this.parse(ShareCommand)

    if (args.action === 'new') {
      if (!flags.env)
        throw new Error(
          'You need to give the name of the envivronment you want to create the user in !'
        )
      this.newShare(args.action, flags.env)
    } else if (args.action === 'pull') {
      if (!flags.link)
        throw new Error('You need to give the path to the link file ! ')
      this.pull(flags.link)
    } else {
      this.log(`The action ${chalk.bold(args.action)} is not a valid one`)
    }
  }
}

ShareCommand.description = `Share your file file with a non-blockstack user
`

ShareCommand.examples = [chalk.yellow('$ ks share')]

ShareCommand.args = [
  {
    name: 'action',
    required: true, // make the arg required with `required: true`
    description: `new || pull. Create a new shared user or pull files based on ${SHARE_FILENAME} file.`, // help description
    hidden: false,
  },
]

ShareCommand.flags = {
  env: flags.string({
    char: 'e',
    multiple: false,
    description: `Env you want to create the user in.`,
  }),
  link: flags.string({
    char: 'l',
    multiple: false,
    description: `Path to your link file.`,
  }),
}

module.exports = ShareCommand
