const debug = require('debug')('keystone:command:project')
const { cli } = require('cli-ux')
const chalk = require('chalk')
const { assertUserIsAdminOrContributor } = require('@ks/core/lib/member')
const {
  getLatestProjectDescriptor,
  getLatestMembersDescriptor,
  getMembers,
} = require('@ks/core/lib/descriptor')

const { config } = require('@ks/core/lib/commands/project')
const { CommandSignedIn } = require('../lib/commands')

class ProjectCommand extends CommandSignedIn {
  async saveChanges(project, projectDescriptor) {
    await this.withUserSession(async userSession => {
      try {
        config(userSession, {
          project,
          descriptor: projectDescriptor,
          type: 'project',
        })
      } catch (err) {
        cli.action.stop('Failed')
        this.log(err)
      }
    })
  }

  async configureProject(project) {
    await this.withUserSession(async userSession => {
      await assertUserIsAdminOrContributor(userSession, { project })

      try {
        debug('Get last project descriptor')

        console.log('\x1Bc')

        const allMembers = await getMembers(userSession, { project })

        const projectMembersDescriptor = await getLatestMembersDescriptor(
          userSession,
          {
            project,
            type: 'members',
          }
        )

        await this.configureMembers({
          allMembers,
          projectMembers: projectMembersDescriptor.content,
          envsMembers: { [project]: projectMembersDescriptor.content },
          project: true,
          env: project,
          currentStep: 1,
          type: 'project',
        })

        await this.saveChanges(project, projectMembersDescriptor)
      } catch (err) {
        this.log(chalk.bold(err))
      }
    })
  }

  async run() {
    const { args } = this.parse(ProjectCommand)
    const project = await this.getProjectName()

    try {
      if (args.action) {
        if (args.action === 'config') {
          this.configureProject(project)
        }
      }
    } catch (error) {
      this.log(`${chalk.red(error)}`)
    }
  }
}

ProjectCommand.description = `Manage project.`

ProjectCommand.args = [
  {
    name: 'action',
    required: false, // make the arg required with `required: true`
    description: 'Configure project members', // help description
    hidden: false,
  },
]

ProjectCommand.examples = [chalk.yellow('$ ks env config')]

module.exports = ProjectCommand
