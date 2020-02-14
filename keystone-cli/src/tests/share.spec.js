require('./utils/mock')
const { prepareEnvironment } = require('./utils')
jest.mock('../lib/blockstackLoader')
jest.mock('../lib/commands')

const fs = require('fs')
const path = require('path')
const uuid = require('uuid/v4')
const { stdin } = require('mock-stdin')
const ShareCommand = require('../commands/share')
const PullCommand = require('../commands/pull')
const PushCommand = require('../commands/push')
const { login, logout, runCommand } = require('./utils/helpers')

const fsp = fs.promises

describe('Share Command', () => {
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

  it('should create a shared user', async () => {
    await prepareEnvironment()

    await login()

    await runCommand(ShareCommand, ['default'])

    const sharedUserCreated = result.find(
      log => log.indexOf(`This token must be stored`) > -1
    )
    expect(sharedUserCreated).toBeDefined()
  }, 20000)

  it('Should pull files with shared user token', async () => {
    await prepareEnvironment()
    await login()
    const pathToFile = path.join(__dirname, './local/bar.txt')

    // Create shared user token
    await runCommand(ShareCommand, ['default'])

    // Create a file and push
    const uid = uuid()
    await fsp.writeFile(pathToFile, uid)
    await runCommand(PushCommand, ['bar.txt'])

    await logout()

    // store token generated
    const token = result.find(log => log.indexOf(`▻`) > -1).split(' ')[1]

    process.env.KEYSTONE_SHARED = token.replace(/.[[0-9]+m/g, '') // remove chalk colors

    // pull files using the token
    await runCommand(PullCommand, [])

    process.env.KEYSTONE_SHARED = null

    const content = (await fsp.readFile(pathToFile)).toString()
    expect(content).toEqual(uid)
  })
})
