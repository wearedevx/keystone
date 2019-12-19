const { Command, flags } = require('@oclif/command')
const { cli } = require('cli-ux')
const chalk = require('chalk')
const normalize = require('normalize-path')
const treeify = require('treeify')
const inquirer = require('inquirer')
const path = require('path')
const { isOneOrMoreAdmin, setMembersToEnvs } = require('@keystone/core/lib/env')
const { logo } = require('../lib/ux')
// const { getProjectDescriptor } = require('../lib/core')
const { getSession, getProjectConfig } = require('../lib/blockstackLoader')

function promptEnvToChange(envs, project = false) {
  return inquirer.prompt([
    {
      type: 'list',
      name: 'env',
      message: `Which ${
        project ? 'project' : 'environment'
      } you want to configure?`,
      choices: [
        ...envs,
        new inquirer.Separator(),
        chalk.green('Save & quit'),
        chalk.red('Cancel'),
      ],
    },
  ])
}

function promptUsersToChange(role, env, allMembers, envsMembers, envs) {
  return inquirer.prompt([
    {
      type: 'checkbox',
      name: 'members',
      message: `Which users you want to make ${chalk.green(
        role
      )} for the ${chalk.green(env)} environment?`,
      choices: [
        ...allMembers.map(user => {
          if (
            envsMembers[env][role].find(
              member => member.blockstack_id === user.blockstack_id
            )
          ) {
            return { name: user.blockstack_id, value: user, checked: true }
          }
          return { name: user.blockstack_id, value: user }
        }),
      ],
      default: [envs],
    },
  ])
}

function promptRoleToChange(envsMembers, env) {
  return inquirer.prompt([
    {
      type: 'list',
      name: 'role',
      message: 'Which role you want to set?',
      choices: [
        {
          value: 'readers',
          name: `readers ${chalk.gray(`(${envsMembers[env].readers.length})`)}`,
        },
        {
          value: 'contributors',
          name: `contributors ${chalk.gray(
            `(${envsMembers[env].contributors.length})`
          )}`,
        },
        {
          value: 'admins',
          name: `admins ${chalk.gray(`(${envsMembers[env].admins.length})`)}`,
        },
        new inquirer.Separator(),
        chalk.red('Cancel'),
      ],
    },
  ])
}

const CommandSignedIn = class extends Command {
  async configureMembers(props) {
    try {
      let { envsMembers, currentStep, env, role } = props
      const { allMembers, project } = props

      const envs = Object.keys(envsMembers)

      // Nicely log env config with color
      let members = JSON.parse(JSON.stringify(envsMembers))

      members = Object.keys(members).reduce((acc, role) => {
        acc[chalk.red(role)] = Object.keys(members[role]).reduce(
          (acc2, curr2) => {
            acc2[chalk.blue(curr2)] = members[role][curr2].map(x =>
              chalk.green(x.blockstack_id)
            )
            return acc2
          },
          {}
        )
        return acc
      }, {})

      console.log(treeify.asTree(members, true))
      if (!isOneOrMoreAdmin(envsMembers)) {
        this.log(
          chalk.red(
            '/!\\ There must be one or more admin in each environment\n'
          )
        )
      }

      if (currentStep === 0) {
        const result = await promptEnvToChange(envs, project)

        if (/.*cancel.*/i.test(result.env)) {
          if (
            await CommandSignedIn.confirm(
              chalk.red('Sure you want to quit and abort all changes ?')
            )
          ) {
            process.exit(0)
          }
        } else if (/.*save.*/i.test(result.env)) {
          if (
            await CommandSignedIn.confirm(
              chalk.green('Sure you want to save your changes and quit ?')
            )
          ) {
            if (isOneOrMoreAdmin(envsMembers)) {
              return envsMembers
            }
          }
        } else {
          env = result.env
          currentStep = 1
        }
      }
      if (currentStep === 1) {
        const result = await promptRoleToChange(envsMembers, env)
        role = result.role

        if (/.*ancel.*/i.test(role)) {
          currentStep = 0
        } else {
          currentStep = 2
        }
      }
      if (currentStep === 2) {
        let { members } = await promptUsersToChange(
          role,
          env,
          allMembers,
          envsMembers,
          envs
        )
        console.log('TCL: extends -> configureMembers -> members', members)

        currentStep = 0

        // members = members.map(member => ({
        //   blockstack_id: member,
        // }))
        envsMembers = setMembersToEnvs({ envsMembers, members, role, env })
      }
      console.log('\x1Bc')

      return this.configureMembers({
        allMembers,
        envsMembers,
        envs,
        currentStep,
        env,
        role,
        project,
      })
    } catch (err) {
      console.error(err)
    }
  }

  // async fetchProject(userSession, project) {
  //   try {
  //     cli.action.start('Fetching project')
  //     return await getProjectDescriptor(userSession, { project })
  //   } catch (error) {
  //     cli.action.stop('Failed')
  //     this.log(`▻ ${error.message} ${chalk.red.bold('✗')}`)
  //     // end the cli
  //     return false
  //   }
  // }

  async getConfigFolderPath() {
    try {
      const { absoluteProjectPath } = await getProjectConfig()
      return absoluteProjectPath
    } catch (error) {
      console.error(error)
      process.exit(2)
    }
  }

  async getProjectEnv() {
    try {
      const {
        config: { env },
      } = await getProjectConfig()
      return env
    } catch (error) {
      this.log(
        `▻ Keystone config file is missing or malformed, please start with ${chalk.yellow(
          `$ ks init`
        )}`
      )
      process.exit(1)
    }
  }

  async getProjectName() {
    try {
      const {
        config: { project },
      } = await getProjectConfig()

      return project
    } catch (error) {
      this.log(
        `▻ Keystone config file is missing or malformed, please start with ${chalk.yellow(
          `$ ks init`
        )}`
      )
      process.exit(1)
    }
  }

  async getDefaultDirectory(flags) {
    if (flags && flags.directory) {
      return flags.directory
    }
    try {
      const { config } = await getProjectConfig()
      return config.directory
    } catch (error) {
      this.log(
        `▻ Keystone config file is missing or malformed, please start with ${chalk.yellow(
          `$ ks init`
        )}`
      )
      process.exit(0)
    }
  }

  async withUserSession(callback, options = {}) {
    const { configPath } = options
    const userSession = await this.getUserSession(configPath)
    if (userSession) {
      await callback(userSession)
    }
  }

  async getUserSession(configPath) {
    try {
      const config = configPath || this.config.configDir
      const userSession = await getSession(config)
      if (userSession && userSession.isUserSignedIn()) {
        return userSession
      }
      this.log(
        `▻ You're not connected, please sign in first: ${chalk.yellow(
          `$ ks login`
        )}`
      )
    } catch (error) {
      this.log(
        `▻ You're not connected, please sign in first: ${chalk.yellow(
          `$ ks login`
        )}`
      )
    }
    return false
  }

  async init() {
    this.log(logo)
  }

  async getFileRelativePath(filePath) {
    const absoluteFilePath = path.resolve(filePath)
    const absoluteConfigFolderPath = await this.getConfigFolderPath()

    if (absoluteFilePath.indexOf(absoluteConfigFolderPath) === -1) {
      throw new Error(
        `${filePath} is not in the keystone project ${absoluteConfigFolderPath}`
      )
    }

    const relativeFilePath = absoluteFilePath.replace(
      `${absoluteConfigFolderPath}/`,
      ''
    )

    return relativeFilePath
  }

  async getAbsolutePath(filePath) {
    const absoluteConfigFolderPath = await this.getConfigFolderPath()
    const absoluteFilePath = path.join(absoluteConfigFolderPath, filePath)
    return absoluteFilePath
  }
}

CommandSignedIn.normalizePath = (directory, filename) => {
  let dir = directory
  if (!directory) {
    dir = './'
  }
  return normalize(`${dir}/${filename}`)
}

CommandSignedIn.confirm = async message => {
  const answer = await inquirer.prompt([
    {
      name: 'confirm',
      type: 'confirm',
      message,
    },
  ])
  return answer.confirm
}

CommandSignedIn.flags = {
  project: flags.string({
    char: 'p',
    multiple: false,
    description: 'Set the project',
  }),
  config: flags.string({
    char: 'c',
    multiple: false,
    description: 'Set the path to the blockstack session file',
  }),
}

module.exports = {
  CommandSignedIn,
}
