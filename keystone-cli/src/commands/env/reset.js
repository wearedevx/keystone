const chalk = require('chalk')
const inquirer = require('inquirer')

const { resetLocalFiles } = require('@keystone.sh/core/lib/file/disk')

const { CommandSignedIn } = require('../../lib/commands')

class EnvCommand extends CommandSignedIn {
  async resetEnv() {
    const absoluteProjectPath = await this.getConfigFolderPath()
    try {
      resetLocalFiles(absoluteProjectPath)
    } catch (err) {
      if (err.code === 'NoPendingModification') {
        console.log('No changes made to files.')
        process.exit(0)
      }
      if (err.code === 'PendingModification') {
        err.data.forEach(f => console.log(f.path, chalk.bold(f.status)))
        console.log('\n')
        const { confirm } = await inquirer.prompt([
          {
            type: 'confirm',
            name: 'confirm',
            message: `Are you sure you want to reset the following changes ?`,
          },
        ])
        if (confirm) resetLocalFiles(absoluteProjectPath, confirm)
        else process.exit(0)
      }
    }
  }

  async run() {
    await this.resetEnv()
  }
}

EnvCommand.description = `reset changes you made locally in tracked files`

EnvCommand.examples = [chalk.blue(`$ ks env reset`)]

module.exports = EnvCommand
