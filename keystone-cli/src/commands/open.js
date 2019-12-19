// const fs = require('fs')
// const { flags } = require('@oclif/command')
// const { cli } = require('cli-ux')
// const chalk = require('chalk')
// const openEditor = require('open')

// const { getFiles } = require('../lib/core/file')
const { CommandSignedIn } = require('../lib/commands')

// const fsp = fs.promises
class OpenCommand extends CommandSignedIn {
//   async open(files, project, directory = '') {
//     await this.withUserSession(async userSession => {
//       try {
//         cli.action.start(`Fetching`)
//         const fetchedFiles = await getFiles(userSession, { project, files })
//         await Promise.all(
//           fetchedFiles.map(async file => {
//             if (file.fetched) {
//               // handle path better than this naive guess
//               const path = CommandSignedIn.normalizePath(directory, file.name)

//               try {
//                 await fsp.writeFile(path, file.descriptor.content)
//                 try {
//                   await openEditor(path)
//                 } catch (error) {
//                   this.log(
//                     `▻ Can't open your editor for ${chalk.yellow(
//                       file.name
//                     )} ${chalk.red.bold('✗')}`,
//                     file.error
//                   )
//                 }
//               } catch (error) {
//                 this.log(
//                   `▻ Can't write ${chalk.yellow(
//                     file.name
//                   )} to ${path} ${chalk.red.bold('✗')}`
//                 )
//               }
//               this.log(
//                 `▻ ${chalk.yellow(
//                   file.name
//                 )} written to ${path} ${chalk.green.bold('✓')}`
//               )
//             } else {
//               this.log(
//                 `▻ Can't open ${chalk.yellow(file.name)} ${chalk.red.bold(
//                   '✗'
//                 )}`,
//                 file.error
//               )
//               cli.action.stop(`Failed\n`)
//             }
//           })
//         )
//       } catch (error) {
//         console.error(error)
//         cli.action.stop('failed\n')
//         this.log(`${error.message}`)
//       }
//     })
//   }

//   async run() {
//     const { argv, flags } = this.parse(OpenCommand)

//     try {
//       const project = await this.getProjectName(flags)
//       const directory = await this.getDefaultDirectory(flags)
//       await this.open(argv, project, directory)
//     } catch (error) {
//       this.log(error.message)
//     }
//   }
}

// OpenCommand.description = `Fetch one or more files.
// `

// OpenCommand.examples = [
//   `${chalk.yellow('$ ks open my-file')} ${chalk.gray.italic(
//     '#open my-file from the project set in .ksconfig'
//   )}`,
//   `${chalk.yellow(
//     '$ ks open my-file --project=my-project'
//   )} ${chalk.gray.italic('#open my-file from my-project')}`,
//   `${chalk.yellow('$ ks open my-file --directory=config/')} ${chalk.gray.italic(
//     '#open my-file and copy to directory config/'
//   )}`,
// ]

// OpenCommand.args = [
//   {
//     name: 'files',
//     required: false, // make the arg required with `required: true`
//     description: 'Open a file in your default editor', // help description
//     hidden: false,
//   },
// ]

// OpenCommand.flags = {
//   ...CommandSignedIn.flags,
//   directory: flags.string({
//     char: 'd',
//     multiple: false,
//     description: 'Set the destination folder',
//   }),
// }

// OpenCommand.strict = false

module.exports = OpenCommand
