require('./utils/mock')
const { stdin } = require('mock-stdin')
const fs = require('fs')
const path = require('path')
const { prepareEnvironment } = require('./utils')

jest.mock('../lib/blockstackLoader')
jest.mock('../lib/commands')

const MemberAddCommand = require('../commands/member/add')
const MemberRemoveCommand = require('../commands/member/rm')
const { runCommand } = require('./utils/helpers')

describe('Invite Command', () => {
  let result
  let io

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

  afterEach(() => {
    jest.restoreAllMocks()
    io.restore()
  })

  it('should add a member to a project', async () => {
    await prepareEnvironment()

    const username = 'keystone_test2.id.blockstack'
    await runCommand(MemberRemoveCommand, ['-u', username])

    // gen pub key file for new user
    fs.writeFileSync(path.join(__dirname, './hub', `${username}--public.key`))

    const projects = fs
      .readFileSync(
        path.join(
          __dirname,
          './hub/',
          `keystone_test1.id.blockstack--projects.json`
        )
      )
      .toString()
    fs.writeFileSync(
      path.join(__dirname, './hub/', `${username}--projects.json`),
      projects
    )

    await runCommand(MemberAddCommand, [username, 'test2@keystone.sh'])

    const invited = result.find(log =>
      log.indexOf('keystone_test2.id.blockstack added to')
    )
    expect(invited).toBeDefined()
  }, 20000)
})
