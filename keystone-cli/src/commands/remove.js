const { flags } = require('@oclif/command')
const chalk = require('chalk')
const fs = require('fs')

const { removeProject } = require('../lib/core/project')
const { CommandSignedIn } = require('../lib/commands')

class RemoveCommand extends CommandSignedIn {
  async removeProject(name, force) {
    await this.withUserSession(async userSession => {
      try {
        await removeProject(userSession, { name, force })
        this.log(
          `▻ Project ${chalk.bold(
            name
          )} successfully removed ${chalk.green.bold('✓')}\n`
        )
        const configFile = fs.readFileSync('.ksconfig')
        const { project } = JSON.parse(configFile)
        if (name === project) {
          const confirm = await CommandSignedIn.confirm(
            'Remove keystone config?'
          )
          if (confirm) {
            fs.unlinkSync('.ksconfig')
            this.log(`▻ Config file removed ${chalk.green.bold('✓')}\n`)
          }
        }
      } catch (error) {
        console.error(error)
        this.log(`▻ ${chalk.red(error.message)}\n`)
      }
    })
  }

  async run() {
    const { flags } = this.parse(RemoveCommand)
    console.log(flags)

    try {
      const project = await this.getProjectName(flags)
      const { force } = flags
      await this.removeProject(project, force)
    } catch (error) {
      this.log(error.message)
    }
  }
}

RemoveCommand.description = `Remove a project.
...
If you're an administrator, the project will be removed for everyone.\n
If you're a contributor or a reader, you will be removed from the project.
`

RemoveCommand.examples = [
  `${chalk.yellow('$ ks remove')} ${chalk.gray.italic(
    '#remove your project set in .ksconfig and all its files'
  )}`,
  `${chalk.yellow('$ ks remove --project=my-project')} ${chalk.gray.italic(
    '#remove your project called my-project'
  )}`,
]

RemoveCommand.flags = {
  ...CommandSignedIn.flags,
  force: flags.string({
    char: 'f',
    multiple: false,
    description:
      'Force the deletion. Beware, it might let some files from the project on your storage.',
  }),
}

module.exports = RemoveCommand
