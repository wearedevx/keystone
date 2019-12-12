const inquirer = require('inquirer')
const { flags } = require('@oclif/command')
const chalk = require('chalk')
const { CommandSignedIn } = require('../lib/commands')
const {
  deleteInvites,
  checkInvitations,
} = require('@keystone/core/lib/invitation')

const { invite } = require('@keystone/core/lib/commands/invite')

const askEmail = async defaultEmail => {
  const answer = await inquirer.prompt([
    {
      name: 'email',
      type: 'input',
      message: `What's your email address?`,
      default: defaultEmail,
    },
  ])
  return answer.email
}

class InviteCommand extends CommandSignedIn {
  async check() {
    await this.withUserSession(async userSession => {
      try {
        const projectsStatuses = await checkInvitations(userSession)
        projectsStatuses.forEach(statuses => {
          if (statuses.invite === 'fulfilled') {
            this.log(
              `▻ ${chalk.yellow(statuses.name)} ${chalk.green.bold('✓')}`
            )
          } else {
            this.log(`▻ ${statuses.name} ${chalk.red.bold('✗')}`)
          }
        })
      } catch (err) {
        this.log(err)
      }
    })
  }

  async invite(emails, project, role, removal) {
    await this.withUserSession(async userSession => {
      if (removal) {
        const deletedInvites = await deleteInvites(userSession, {
          project,
          emails,
        })
        deletedInvites.forEach(i => {
          this.log(
            `▻ invitation for ${chalk.yellow(
              i.email
            )} has been deleted ${chalk.green.bold('✓')}`
          )
        })
        return true
      }

      const { email } = userSession.loadUserData()

      const from = await askEmail(email)

      try {
        const invitations = await invite(userSession, {
          from,
          project,
          emails,
          role,
        })
        invitations.forEach(invitation => {
          if (invitation.sent) {
            this.log(
              `▻ invitation as ${role} sent to ${chalk.yellow(
                invitation.email
              )} ${chalk.green.bold('✓')}`
            )
          } else {
            this.log(
              `▻ invitation to ${chalk.yellow(
                invitation.email
              )} failed ${chalk.red.bold('✗')}`,
              invitation.error
            )
          }
        })
      } catch (error) {
        console.error(error)
        this.log(`${error.message}`)
      }
    })
  }

  async run() {
    try {
      const { argv, flags } = this.parse(InviteCommand)

      if (flags.check) {
        await this.check()
        return
      }

      if (flags.accept) {
        await this.accept()
        return
      }
      // at least 1 arguments required, an email
      if (argv.length >= 1) {
        const project = await this.getProjectName(flags)
        await this.invite(argv, project, flags.role, flags.removal)
      } else {
        // await this.prompt()
        this.log(`We need at least one email to send an invitation`)
      }
    } catch (error) {
      this.log(error.message)
    }
  }
}

InviteCommand.args = [
  {
    name: 'emails',
    required: false, // make the arg required with `required: true`
    description: 'Emails for invitations to be sent', // help description
    hidden: false,
  },
]

InviteCommand.flags = {
  ...CommandSignedIn.flags,
  role: flags.string({
    char: 'r',
    multiple: false,
    options: ['reader', 'contributor', 'admin'],
    default: 'reader',
    description: `Assigns a role`,
  }),
  removal: flags.boolean({
    multiple: false,
    default: false,
    description: `Deletes an invitation`,
  }),
  check: flags.boolean({
    multiple: false,
    default: false,
    description: `Check your pending invitations`,
  }),
}

InviteCommand.description = `Invites one or more people by email to a project.

By default, people you invite are readers. 
You can change the role set by using the role flag. You have 3 choices:
- reader: can only read files from the project.
- contributor: can read, write and add new files to the project
- admin: all the above plus ask people to join the project
`

InviteCommand.examples = [
  `${chalk.yellow('$ ks invite friend@example.com')} ${chalk.gray.italic(
    '#Send an invitation to friend@example.com as a reader on the project'
  )}`,
  `${chalk.yellow(
    '$ ks invite friend@example.com friend2@example.com --role=admin'
  )} ${chalk.gray.italic('#Invite as admin on the project')}`,
  `${chalk.yellow(
    '$ ks invite friend@example.com friend2@example.com --removal'
  )} ${chalk.gray.italic('#Removes the invitations for friend and friend2')}`,
]

module.exports = InviteCommand
