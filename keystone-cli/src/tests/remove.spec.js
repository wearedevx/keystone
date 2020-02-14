// const { stdin } = require('mock-stdin')
// const fs = require('fs')

// const RemoveCommand = require('../commands/remove')
// const InitCommand = require('../commands/init')
// const { login, logout, runCommand } = require('./helpers')

// describe('Remove Command', () => {
//   let result
//   let io
//   const keys = {
//     up: '\x1B\x5B\x41',
//     down: '\x1B\x5B\x42',
//     enter: '\x0D',
//     space: '\x20',
//   }
//   beforeEach(() => {
//     // catch everything on stdout
//     // and put it in result
//     result = []
//     jest
//       .spyOn(process.stdout, 'write')
//       .mockImplementation(val => result.push(val))
//     io = stdin()
//   })

//   afterEach(() => {
//     jest.restoreAllMocks()
//     io.restore()
//   })

//   it('Remove - Not connected', async () => {
//     await logout()
//     await runCommand(RemoveCommand, ['--project=new-project'])
//     expect(result[result.length - 1]).toContain(
//       `You're not connected, please sign in first`
//     )
//   }, 20000)

//   it('should remove a project after initializing it and delete config file', async () => {
//     await login()

//     const existingConfig = fs.readFileSync('.ksconfig')

//     const sendKeystrokes = async () => {
//       io.send(keys.enter)
//     }
//     setTimeout(() => sendKeystrokes().then(), 500)
//     await runCommand(InitCommand, ['new-project'])

//     const createdProject = result.find(log =>
//       /.*Project .* successfully created/g.test(log)
//     )

//     const interval = setInterval(() => {
//       if (result.find(log => log.indexOf('Remove keystone config') > -1)) {
//         sendKeystrokes().then()
//         clearInterval(interval)
//       }
//     }, 500)
//     await runCommand(RemoveCommand, [
//       `--project=${createdProject.replace(/.[[0-9]+m/g, '').split(' ')[2]}`,
//     ])
//     const removed = result.find(log =>
//       log.indexOf('Project .* successfully removed')
//     )
//     expect(removed).toBeDefined()
//     expect(fs.existsSync('.ksconfig')).toBeFalsy()

//     fs.writeFile('.ksconfig', existingConfig)
//   }, 20000)
// })
