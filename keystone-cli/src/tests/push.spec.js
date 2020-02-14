const { writeFileToDisk } = require('./utils/mock')
const { prepareEnvironment } = require('./utils')
const fs = require('fs')
const pathUtil = require('path')
const PushCommand = require('../commands/push')
const { createDescriptor } = require('./utils')
const { login, logout, runCommand } = require('./utils/helpers')

jest.mock('../lib/blockstackLoader')
jest.mock('../lib/commands')
describe('Push Command', () => {
  let result

  beforeEach(() => {
    // catch everything on stdout
    // and put it in result
    result = []
    jest.spyOn(process.stdout, 'write').mockImplementation(val => {
      fs.appendFile('unit-test.log', val)
      result.push(val)
    })
  })

  afterEach(() => jest.restoreAllMocks())

  // it('Push - Not connected', async () => {
  //   await logout()
  //   fs.writeFile(`test.txt`, 'echo')
  //   await runCommand(PushCommand, ['test.txt'])
  //   expect(result[result.length - 1]).toContain(
  //     `You're not connected, please sign in first`
  //   )
  // })
  it('should create a file in current working project', async () => {
    await prepareEnvironment()
    await login()

    const fileDescriptor = createDescriptor({})
    await writeFileToDisk(fileDescriptor)
    await runCommand(PushCommand, [pathUtil.join(fileDescriptor.name)])

    const pushedFile = result.find(
      log => log.indexOf('pushed') > -1 || log.indexOf('already is the latest')
    )
    expect(pushedFile).toBeDefined()
  })
  // it('should not create a file in current working project because already exist', async () => {
  //   await login()
  //   await runCommand(PushCommand, [`./src/tests/test.txt`])
  //   await runCommand(PushCommand, [`./src/tests/test.txt`])
  //   const pushedFile = result.find(
  //     log =>
  //       log.indexOf(
  //         'A version of this file with the same content already exists.'
  //       ) > -1
  //   )
  //   expect(pushedFile).toBeDefined()
  // })
})
