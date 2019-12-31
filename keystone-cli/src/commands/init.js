const chalk = require('chalk')
const inquirer = require('inquirer')
const init = require('@keystone/core/lib/commands/init')
const { CommandSignedIn } = require('../lib/commands')

class InitCommand extends CommandSignedIn {
  async prompt() {
    const answer = await inquirer.prompt([
      {
        name: 'project_name',
        message: `What's your project name?`,
      },
    ])
    if (answer.project_name.length > 1) {
      await this.initProject(answer.project_name)
    } else {
      await this.prompt()
    }
  }

  async initProject(project, overwrite = false) {
    return new Promise(async (resolve, reject) => {
      await this.withUserSession(async userSession => {
        try {
          const projectWithId = await init(userSession, { project, overwrite })
          this.log(
            `▻ Project ${chalk.bold(
              projectWithId
            )} successfully created ${chalk.green.bold('✓')}`
          )
          this.log(
            `▻ You can add files with: ${chalk.yellow(`$ ks push my-file`)}`
          )
          resolve()
        } catch (error) {
          switch (error.code) {
            case 'ProjectNameExists':
              {
                const choice = await inquirer.prompt([
                  {
                    type: 'list',
                    name: 'project',
                    message:
                      'One or more projects have the same name \n Choose one or pick another name !',
                    choices: [...error.data, 'Type a new name'],
                  },
                ])

                await this.initProject(choice.project)
              }
              break
            case 'ConfigFileExists':
              {
                const overwriteConfig = await CommandSignedIn.confirm(
                  'A config file already exists, overwrite?'
                )
                if (overwriteConfig) {
                  await this.initProject(project, overwriteConfig)
                }
                resolve()
              }
              break
            default:
              reject(error)
          }
        }
      })
      resolve()
    })
  }

  async run() {
    try {
      const { args } = this.parse(InitCommand)
      if (args.project_name) {
        await this.initProject(args.project_name)
      } else {
        await this.prompt()
      }
    } catch (error) {
      console.error(error)
      this.log(`${error.message}`)
    }
  }
}

InitCommand.description = `Create Keystone config file
`

InitCommand.args = [
  {
    name: 'project_name',
    required: false, // make the arg required with `required: true`
    description: 'Your project name', // help description
    hidden: false,
  },
]

InitCommand.examples = [chalk.yellow('$ ks init project-name')]

module.exports = InitCommand
