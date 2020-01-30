const { flags } = require('@oclif/command')
const { cli } = require('cli-ux')
const chalk = require('chalk')
const { ROLES } = require('@keystone.sh/core/lib/constants')
const { add } = require('@keystone.sh/core/lib/commands/add')
const { assertUserIsAdmin } = require('@keystone.sh/core/lib/member')
const { CommandSignedIn } = require('../lib/commands')

class AddCommand extends CommandSignedIn {
  async add(blockstackId, email, project) {
    const userSession = await this.getUserSession()

    await assertUserIsAdmin(userSession, { project })

    const invitee = {
      blockstackId,
      email,
      role: ROLES.READERS,
    }
    let success
    try {
      cli.action.start(`Adding ${blockstackId} to the poject`)
      const { added, memberAdded } = await add(userSession, {
        invitee,
        project,
      })

      if (added) {
        success = true
        this.log(
          `▻ ${chalk.yellow(
            memberAdded
          )} added to ${project} ${chalk.green.bold('✓')}`
        )
      } else {
        success = false

        this.log(
          `▻ Failed to add ${chalk.yellow(
            memberAdded
          )} to ${project} ${chalk.red.bold('✗')}`,
          memberAdded.error
        )
      }
    } catch (error) {
      success = false

      console.error(error)
      this.log(`\n${chalk.red(error.code)} : ${error.message}`)
    }
    cli.action.stop(success ? 'done' : 'failed')
  }

  async run() {
    try {
      const { argv, args } = this.parse(AddCommand)
      // at least 1 arguments required, a blockstack id
      if (argv.length >= 1) {
        const project = await this.getProjectName()
        await this.add(args.blockstackId, args.email, project)
      } else {
        // await this.prompt()
        this.log(
          `We need at  least a blockstack id and an email associated to an invitation`
        )
      }
    } catch (error) {
      console.error(error)
      this.log(error.message)
    }
  }
}

AddCommand.args = [
  {
    name: 'blockstackId',
    required: true, // make the arg required with `required: true`
    description: 'Blockstack_id to add', // help description
    hidden: false,
  },
  {
    name: 'email',
    required: true, // make the arg required with `required: true`
    description: 'email associated to an invitation', // help description
    hidden: false,
  },
]

AddCommand.description = `Add a member to a project.

Adding a member give them access to the project.
The member should have accepted your invitation for this to work

You  can add the member to an environment with : ${chalk.yellow(
  'ks env config'
)}
`

AddCommand.examples = [
  `${chalk.yellow(
    '$ ks add example.id.blockstack example@mail.com'
  )} ${chalk.gray.italic('#add a user to a project')}`,
]

module.exports = AddCommand
