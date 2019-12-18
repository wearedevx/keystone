const chalk = require('chalk')
const fs = require('fs')
const { newShare, pullShared } = require('@keystone/core/lib/commands/share')
const { flags } = require('@oclif/command')
const { SHARE_FILENAME } = require('@keystone/core/lib/constants')

const { CommandSignedIn } = require('../lib/commands')

class ShareCommand extends CommandSignedIn {
  async newShare(action, env) {
    await this.withUserSession(async userSession => {
      const { username } = userSession.loadUserData()
      const project = await this.getProjectName()
      if (action === 'new') {
        const addedShare = await newShare(userSession, { project, env })

        fs.writeFile(
          SHARE_FILENAME,
          JSON.stringify({
            project,
            env,
            member: username,
            privateKey: addedShare.privateKey,
          }),
          err => console.log(err)
        )

        this.log(
          `Private key to decrypt shared user files :\n▻ ${chalk.yellow(
            addedShare.privateKey
          )}`
        )
      }
    })
  }

  async pull(pathToConfig) {
    await this.withUserSession(async userSession => {
      const { project, env, member, privateKey } = JSON.parse(
        fs.readFileSync(pathToConfig)
      )
      userSession.sharedPrivateKey = privateKey
      const absoluteProjectPath = await this.getConfigFolderPath()

      await pullShared(userSession, {
        project,
        env,
        origin: member,
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
