require('./utils/mock')
const fs = require('fs')

const PullCommand = require('../commands/pull')
const { login, runCommand } = require('./utils/helpers')

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

  fit('should pull a file from hub', async () => {
    await login()

    await runCommand(PullCommand)

    const pulledFile = result.find(log => log.indexOf('pushed') > -1)
    expect(pulledFile).toBeDefined()
  })
})
