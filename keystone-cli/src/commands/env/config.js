const { cli } = require('cli-ux')
const chalk = require('chalk')
const { assertUserIsAdmin } = require('@keystone.sh/core/lib/member')
const {
  getLatestMembersDescriptor,
  getMembers,
  getLatestProjectDescriptor,
} = require('@keystone.sh/core/lib/descriptor')

const { config } = require('@keystone.sh/core/lib/commands/env')
const { ROLES } = require('@keystone.sh/core/lib/constants')

const { CommandSignedIn } = require('../../lib/commands')

class EnvCommand extends CommandSignedIn {
  async saveChanges(project, envsDescriptor, type) {
    await this.withUserSession(async userSession => {
      try {
        await config(userSession, {
          project,
          descriptors: envsDescriptor,
          type,
        })
      } catch (err) {
        cli.action.stop('failed')
        this.log(err)
      }
    })
  }

  async configureEnv(project) {
    await this.withUserSession(async userSession => {
      await assertUserIsAdmin(userSession, { project })
      const { username } = userSession.loadUserData()
      try {
        // Check if env already exists.
        const projectDescriptor = await getLatestProjectDescriptor(
          userSession,
          {
            project,
          }
        )

        const envMembersDescriptors = await Promise.all(
          projectDescriptor.content.env.map(async env => {
            const descriptor = await getLatestMembersDescriptor(userSession, {
              project,
              env,
            })

            return { env, descriptor }
          })
        )

        let envsMembers = envMembersDescriptors.reduce(
          (envs, { env, descriptor }) => {
            envs[env] = descriptor.content
            return envs
          },
          {}
        )

        // only environments admins can change users permissions
        // so we keep only environments where the user is admin
        envsMembers = Object.keys(envsMembers).reduce((envs, env) => {
          const isAdmin = envsMembers[env][ROLES.ADMINS].find(
            member => member.blockstack_id === username
          )
          if (isAdmin) {
            envs[env] = envsMembers[env]
          }
          return envs
        }, {})

        const allMembers = await getMembers(userSession, { project })

        this.log('\x1Bc')

        await this.configureMembers({
          allMembers,
          envsMembers,
          currentStep: 0,
          type: 'env',
        })

        await this.saveChanges(project, envMembersDescriptors)
      } catch (err) {
        this.log(chalk.bold(err))
      }
    })
  }

  async run() {
    const project = await this.getProjectName()
    try {
      this.configureEnv(project)
    } catch (error) {
      this.log(`${chalk.red(error)}`)
    }
  }
}

EnvCommand.description = `manage members role in project environments

roles can be the followings :
  reader: can only read files from the the environment and pull them locally
  contributor: can read, write and add new files to the environement
  admin: all the above plus ask people to join the project
`

EnvCommand.examples = [chalk.blue('$ ks env config')]

module.exports = EnvCommand
