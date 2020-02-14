// const CatCommand = require('../commands/cat')
// const PushCommand = require('../commands/push')
// const { login, logout, runCommand } = require('./helpers')

// describe('Cat Command', () => {
//   let result

//   beforeEach(() => {
//     // catch everything on stdout
//     // and put it in result
//     result = []
//     jest
//       .spyOn(process.stdout, 'write')
//       .mockImplementation(val => result.push(val))
//   })

//   afterEach(() => jest.restoreAllMocks())

//   it('Cat - Not connected', async () => {
//     await logout()
//     await runCommand(CatCommand, ['test.txt'])
//     expect(result[result.length - 1]).toContain(
//       `You're not connected, please sign in first`
//     )
//   })

//   it('should print content of file', async () => {
//     await login()
//     await runCommand(PushCommand, ['src/tests/test.txt'])
//     await runCommand(CatCommand, ['test.txt'])
//     const fileContent = result.find(
//       log => log.indexOf('echo') > -1 || log.indexOf('undefined')
//     )
//     expect(fileContent).toBeDefined()
//   })

//   it("should not get the file because don't exist", async () => {
//     await login()
//     await runCommand(CatCommand, ['not_existing.txt'])

//     const fileContent = result.find(log => log.indexOf('returned 404') > -1)
//     expect(fileContent).toBeDefined()
//   })
// })
