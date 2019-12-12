const fs = require('fs')
const uuid = require('uuid/v4')
const PushCommand = require('../commands/push')

const { login, logout, runCommand } = require('./helpers')

describe('Push Command', () => {
  let result

  beforeEach(() => {
    // catch everything on stdout
    // and put it in result
    result = []
    jest
      .spyOn(process.stdout, 'write')
      .mockImplementation(val => result.push(val))
  })

  afterEach(() => jest.restoreAllMocks())

  it('Push - Not connected', async () => {
    await logout()
    fs.writeFile(`test.txt`, 'echo')
    await runCommand(PushCommand, ['test.txt'])
    expect(result[result.length - 1]).toContain(
      `You're not connected, please sign in first`
    )
  })
  it('should create a file in current working project', async () => {
    await login()
    const uid = uuid()
    fs.writeFile(`test-${uid}.txt`, 'echo')
    await runCommand(PushCommand, [`test-${uid}.txt`])
    fs.unlinkSync(`test-${uid}.txt`)
    const pushedFile = result.find(log => log.indexOf('pushed') > -1)
    expect(pushedFile).toBeDefined()
  })
  it('should not create a file in current working project because already exist', async () => {
    await login()
    await runCommand(PushCommand, [`./src/tests/test.txt`])
    await runCommand(PushCommand, [`./src/tests/test.txt`])
    const pushedFile = result.find(
      log =>
        log.indexOf(
          'A version of this file with the same content already exists.'
        ) > -1
    )
    expect(pushedFile).toBeDefined()
  })
})
