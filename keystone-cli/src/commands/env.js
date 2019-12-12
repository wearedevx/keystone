const { cli } = require('cli-ux')
const chalk = require('chalk')
const { flags } = require('@oclif/command')
const { assertUserIsAdminOrContributor } = require('../lib/core-2.0/member')
const {
  createEnv,
  removeEnvFiles,
  checkoutEnv,
} = require('../lib/core-2.0/env')
const {
  getLatestMembersDescriptor,
  getMembers,
  getLatestProjectDescriptor,
} = require('../lib/core-2.0/descriptor')
const {
  addEnvToProject,
  removeEnvFromProject,
} = require('../lib/core-2.0/projects')

const { CommandSignedIn } = require('../lib/commands')
const { config } = require('../lib/core-2.0/commands/env')

class EnvCommand extends CommandSignedIn {
  async saveChanges(project, envsDescriptor, type) {
    await this.withUserSession(async userSession => {
      try {
        config(userSession, { project, descriptors: envsDescriptor, type })
      } catch (err) {
        cli.action.stop('Failed')
        this.log(err)
      }
    })
  }

  /**
   * Create env
   */
  async newEnv(project, name) {
    await this.withUserSession(async userSession => {
      await assertUserIsAdminOrContributor(userSession, { project })

      if (userSession && userSession.isUserSignedIn()) {
        try {
          // Check if env already exists.
          const projectDescriptor = await getLatestProjectDescriptor(
            userSession,
            {
              project,
            }
          )

          // If not, create it.
          if (projectDescriptor.content.env.includes(name)) {
            throw new Error(`Env ${name} already exists.`)
          }

          await createEnv(userSession, {
            env: name,
            projectDescriptor,
          })

          await addEnvToProject(userSession, {
            projectDescriptor,
            env: name,
          })
          this.log(`▻ Environment ${chalk.bold(name)} successfully created`)
        } catch (err) {
          console.log(err)
          this.log(`▻ Environment creation failed : ${chalk.bold(err)}`)
        }
      }
    })
  }

  /**
   * Remove env.
   * @param {*} project
   * @param {*} name
   */
  async removeEnv(project, name) {
    await this.withUserSession(async userSession => {
      await assertUserIsAdminOrContributor(userSession, { project })

      const absoluteProjectPath = await this.getConfigFolderPath()

      await removeEnvFiles(userSession, {
        project,
        env: name,
        absoluteProjectPath,
      })

      // Remove en from project descriptor.
      await removeEnvFromProject(userSession, {
        project,
        env: name,
      })
    })
  }

  async configureEnv(project) {
    await this.withUserSession(async userSession => {
      await assertUserIsAdminOrContributor(userSession, { project })
      try {
        // Check if env already exists.
        const projectDescriptor = await getLatestProjectDescriptor(
          userSession,
          {
            project,
          }
        )

        const envMembversDescriptors = await Promise.all(
          projectDescriptor.content.env.map(async env => {
            const descriptor = await getLatestMembersDescriptor(userSession, {
              project,
              env,
            })

            return { env, descriptor }
          })
        )

        const envsMembers = envMembversDescriptors.reduce(
          (acc, { env, descriptor }) => {
            acc[env] = descriptor.content
            return acc
          },
          {}
        )

        const allMembers = await getMembers(userSession, { project })

        console.log('\x1Bc')

        await this.configureMembers({
          allMembers,
          envsMembers,
          currentStep: 0,
          type: 'env',
        })

        await this.saveChanges(project, envMembversDescriptors)
      } catch (err) {
        this.log(chalk.bold(err))
      }
    })
  }

  async checkout(project, env) {
    await this.withUserSession(async userSession => {
      try {
        const absoluteProjectPath = await this.getConfigFolderPath()

        const newConfig = await checkoutEnv(userSession, {
          project,
          env,
          absoluteProjectPath,
        })
        console.log(newConfig)
      } catch (err) {
        this.log(err)
      }
    })
  }

  async run() {
    const { args } = this.parse(EnvCommand)
    const project = await this.getProjectName()

    try {
      if (args.action) {
        if (args.action === 'config') {
          this.configureEnv(project)
        } else if (args.action === 'new') {
          if (args.env) {
            await this.newEnv(project, args.env)
          } else {
            throw new Error(`You need to give the name of the environment`)
          }
        } else if (args.action === 'remove') {
          if (args.env) {
            await this.removeEnv(project, args.env)
          } else {
            throw new Error(`You need to give the name of the environment`)
          }
        } else if (args.action === 'checkout') {
          console.log(args)
          if (args.env) {
            await this.checkout(project, args.env)
          } else {
            throw new Error(`You need to give the name of the environment`)
          }
        }
      }
    } catch (error) {
      this.log(`${chalk.red(error)}`)
    }
  }
}

EnvCommand.description = `Manage environments.
`

EnvCommand.args = [
  {
    name: 'action',
    required: false, // make the arg required with `required: true`
    description: 'Configure add or remove an environment', // help description
    hidden: false,
  },
  {
    name: 'env',
    required: false, // make the arg required with `required: true`
    description: 'Set working env', // help description
    hidden: false,
  },
]
EnvCommand.flags = {
  name: flags.string({
    char: 'n',
    multiple: false,
    description: `Enviroment name`,
  }),
}

EnvCommand.examples = [chalk.yellow('$ ks env config')]
EnvCommand.examples = [chalk.yellow('$ ks env new --name dev')]
EnvCommand.examples = [chalk.yellow('$ ks env remove --name dev')]

module.exports = EnvCommand
