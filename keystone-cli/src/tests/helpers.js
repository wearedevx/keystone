// const userFolder = require('user-home')
// const Config = require('@oclif/config')
// const fs = require('fs')
// const { createFolder, write, del } = require('../lib/cliStorage')
// const { uploadFile } = require('../lib/core')
// const { getSession } = require('../lib/blockstackLoader')

// // This file is required for testing the CLI
// // It's not versioned as it's a blockstack account linked to the Keystone app.
// // Generate your own:
// // 1) Login with the CLI
// // 2) Copy/paste the file called `session.json` from ~/.config/keystone-cli/
// // TODO: check if there's a session already existing and use it?
// const session = require('./blockstack_session.json')

// // use file API with promises - more elegant.
// const fsp = fs.promises

// const configPath = `${userFolder}/.config/keystone-cli`

// const checkConfigPath = async path => {
//   try {
//     await fsp.access(path, fs.constants.F_OK)
//     console.log(`path ${path} exist.`)
//   } catch (error) {
//     try {
//       await createFolder({ path })
//     } catch (error) {
//       throw error
//     }
//   }
//   return true
// }

// const login = async () => {
//   try {
//     if (await checkConfigPath(configPath)) {
//       await write({
//         path: `${configPath}/`,
//         filename: 'session.json',
//         content: JSON.stringify(session),
//       })
//     }
//   } catch (error) {
//     throw error
//   }
// }

// const logout = async () => {
//   try {
//     if (await checkConfigPath(configPath)) {
//       await del({
//         path: `${configPath}/`,
//         filename: 'session.json',
//       })
//     }
//   } catch (error) {
//     console.log(
//       "Tried to logout but it's already the case or there's an error",
//       error
//     )
//   }
// }

// const runCommand = async (Command, argv = []) => {
//   const defaultConfig = await Config.load()
//   defaultConfig.configDir = configPath
//   const command = new Command(argv, defaultConfig)
//   await command.run()
// }

// const getSessionWithConfig = async () => {
//   const userSession = await getSession(configPath)
//   if (userSession && userSession.isUserSignedIn()) {
//     return userSession
//   }
//   throw new Error("Can't retrieve user session.")
// }

// const putFile = async ({ path, content, encrypt }) => {
//   const userSession = await getSessionWithConfig()
//   await uploadFile(userSession, { path, content, encrypt })
// }

// const removeFile = async ({ path }) => {
//   const userSession = await getSessionWithConfig()
//   await userSession.deleteFile(path)
// }

// module.exports = {
//   login,
//   logout,
//   runCommand,
//   putFile,
//   getSessionWithConfig,
//   removeFile,
// }
