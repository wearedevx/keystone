const fs = require('fs')

const FetchCommand = require('../commands/fetch')
const PushCommand = require('../commands/push')
const { login, logout, runCommand } = require('./helpers')

describe('Fetch Command', () => {
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

  it('Fetch - Not connected', async () => {
    await logout()
    await runCommand(FetchCommand, [])
    expect(result[result.length - 1]).toContain(
      `You're not connected, please sign in first`
    )
  })

  it('should fetch every files in current project', async () => {
    await login()
    await runCommand(PushCommand, ['src/tests/test.txt'])
    await runCommand(FetchCommand, [])
    const writtenFile = result.find(log => log.indexOf('written to') > -1)
    expect(writtenFile).toBeDefined()

    expect(fs.existsSync('test.txt')).toBeTruthy()
    fs.unlinkSync('test.txt')
  }, 20000)

  it('should fetch one selected file in current project', async () => {
    await login()
    await runCommand(PushCommand, ['src/tests/test.txt'])
    await runCommand(FetchCommand, ['test.txt'])
    const writtenFile = result.find(log => log.indexOf('written to') > -1)
    expect(writtenFile).toBeDefined()

    let existedFile
    if (fs.existsSync('test.txt')) {
      existedFile = true
    } else {
      existedFile = false
    }
    fs.unlinkSync('test.txt')
    expect(existedFile).toBeTruthy()
  }, 20000)

  it('should not fetch because file does not exist', async () => {
    await login()

    await runCommand(FetchCommand, ['not_existing.txt'])
    const err = result.find(
      log => log.indexOf('Unable to fetch file from storage') > -1
    )
    expect(err).toBeDefined()
  }, 20000)

  it('should fetch one selected file in selected directory', async () => {
    await login()
    fs.writeFile('src/tests/test-fetched.txt', 'echo')
    await runCommand(PushCommand, ['src/tests/test-fetched.txt'])
    fs.unlinkSync('src/tests/test-fetched.txt')
    await runCommand(FetchCommand, [
      'test-fetched.txt',
      '--directory=src/tests',
    ])
    const writtenFile = result.find(log => log.indexOf('written to') > -1)
    expect(fs.existsSync('src/tests/test-fetched.txt')).toBeTruthy()
    expect(writtenFile).toBeDefined()
    fs.unlinkSync('src/tests/test-fetched.txt')
  }, 20000)
})
