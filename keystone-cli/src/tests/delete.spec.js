require('./utils/mock')
const { prepareEnvironment } = require('./utils')
jest.mock('../lib/blockstackLoader')
jest.mock('../lib/commands')

const fs = require('fs')
const { stdin } = require('mock-stdin')
const DeleteCommand = require('../commands/delete')
const PushCommand = require('../commands/push')
const PullCommand = require('../commands/pull')
const { login, logout, runCommand } = require('./utils/helpers')

describe('Delete Command', () => {
  let result
  // let io

  // const keys = {
  //   up: '\x1B\x5B\x41',
  //   down: '\x1B\x5B\x42',
  //   enter: '\x0D',
  //   space: '\x20',
  // }

  beforeEach(() => {
    // catch everything on stdout
    // and put it in result
    result = []
    jest.spyOn(process.stdout, 'write').mockImplementation(val => {
      fs.appendFile('unit-test.log', val)
      result.push(val)
    })

    io = stdin()
  })

  afterEach(() => jest.restoreAllMocks())
  // it('Delete - Not connected', async () => {
  // await logout()
  // fs.writeFile('test.txt', 'echo')
  // await runCommand(PushCommand, ['test.txt'])
  // await runCommand(DeleteCommand, ['test.txt'])
  // expect(result[result.length - 1]).toContain(
  // `You're not connected, please sign in first`
  // )
  // })
  it('should delete one file after pushing it', async () => {
    await login()
    // Prevent pull before you push error
    await runCommand(PullCommand, ['--force'])

    await runCommand(PushCommand, ['foo.txt'])
    await runCommand(DeleteCommand, ['foo.txt'])
    const deletedFile = result.find(
      log => log.indexOf('successfully deleted') > -1
    )
    expect(deletedFile).toBeDefined()
  }, 20000)

  // it('should delete all files in current project', async () => {
  //   await login()
  //   await runCommand(PushCommand, ['src/tests/test.txt'])

  //   const interval = setInterval(() => {
  //     if (result.find(log => log.indexOf('to delete all files from') > -1)) {
  //       const sendKeystrokes = async () => {
  //         io.send(keys.enter)
  //       }
  //       sendKeystrokes().then()
  //       clearInterval(interval)
  //     }
  //   }, 500)

  //   await runCommand(DeleteCommand, [])

  //   const deletedFile = result.find(log => log.indexOf('deleted') > -1)
  //   expect(deletedFile).toBeDefined()
  // }, 20000)

  // it('should fail because file does not exist', async () => {
  //   await login()

  //   await runCommand(DeleteCommand, ['not_existing.txt'])

  //   const deletedFile = result.find(log => log.indexOf('failed') > -1)
  //   expect(deletedFile).toBeDefined()
  // }, 20000)
})
