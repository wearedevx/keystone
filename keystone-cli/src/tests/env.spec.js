require('./utils/mock')
jest.mock('../lib/blockstackLoader')
jest.mock('../lib/commands')

const fs = require('fs')
const { stdin } = require('mock-stdin')
const EnvCommand = require('../commands/env')
const { login, logout, runCommand } = require('./utils/helpers')

describe('Env Command', () => {
  let result

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

  it('should create a new environment after removing it', async () => {
    await login()
    const envName = 'test_env'
    // Remove the environment if exist, then created it
    await runCommand(EnvCommand, ['remove', envName])
    await runCommand(EnvCommand, ['new', envName])

    const envCreated = result.find(
      log => log.indexOf('successfully created') > -1
    )
    expect(envCreated).toBeDefined()
  }, 20000)
  it('should add a member to an environment', async () => {})
})
