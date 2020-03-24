require('./utils/mock')

jest.mock('../lib/blockstackLoader')
jest.mock('../lib/commands')
const { stdin } = require('mock-stdin')
const fs = require('fs')
const invitationModule = require('@keystone.sh/core/lib/invitation')

// mock email sending
invitationModule.inviteMember = async (
  userSession,
  { from, project, emails, role = 'reader' }
) => {
  console.log(`SEND EMAIL from ${from} to ${emails} on ${project}`)
}

const MemberInviteCommand = require('../commands/member/invite')
const { runCommand } = require('./utils/helpers')

const { prepareEnvironment } = require('./utils')

describe('Invite Command', () => {
  let result
  let io

  const keys = {
    up: '\x1B\x5B\x41',
    down: '\x1B\x5B\x42',
    enter: '\x0D',
    space: '\x20',
  }

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

  it('should send an invitation', async () => {
    await prepareEnvironment()

    await runCommand(MemberInviteCommand, ['test2@keystone.shh', '--removal'])

    const interval = setInterval(() => {
      if (result.find(log => log.indexOf(`What's your email address?`) > -1)) {
        const sendKeystrokes = async () => {
          io.send(keys.enter)
        }
        setTimeout(() => sendKeystrokes().then(), 500)
        clearInterval(interval)
      }
    }, 500)

    await runCommand(MemberInviteCommand, ['test2@keystone.shh'])

    const invited = result.find(log =>
      log.indexOf('invitation as reader sent to')
    )
    expect(invited).toBeDefined()
  }, 20000)

  // it('should delete an invitation after sending it', async () => {
  //   await login()

  //   const interval = setInterval(() => {
  //     if (result.find(log => log.indexOf(`What's your email address?`) > -1)) {
  //       io.send('samuel@wearedevx.com')

  //       const sendKeystrokes = async () => {
  //         io.send(keys.enter)
  //       }
  //       setTimeout(() => sendKeystrokes().then(), 500)
  //       clearInterval(interval)
  //     }
  //   }, 500)

  //   await runCommand(MemberInviteCommand, ['abigael@wearedevx.com'])

  //   await runCommand(MemberInviteCommand, ['abigael@wearedevx.com', '--removal'])

  //   const removed = result.find(log => log.indexOf('has been deleted'))
  //   expect(removed).toBeDefined()
  // }, 20000)

  //   it('should send an invitation as contributor', async () => {
  //     await login()

  //     await runCommand(MemberInviteCommand, ['abigael@wearedevx.com', '--removal'])

  //     const interval = setInterval(() => {
  //       if (result.find(log => log.indexOf(`What's your email address?`) > -1)) {
  //         io.send('samuel@wearedevx.com')

  //         const sendKeystrokes = async () => {
  //           io.send(keys.enter)
  //         }
  //         setTimeout(() => sendKeystrokes().then(), 500)
  //         clearInterval(interval)
  //       }
  //     }, 500)

  //     await runCommand(MemberInviteCommand, [
  //       'abigael@wearedevx.com',
  //       '--role=contributor',
  //     ])

  //     const invited = result.find(log =>
  //       log.indexOf('invitation as contributor sent to')
  //     )
  //     expect(invited).toBeDefined()
  //   }, 20000)

  //   it('should not send an invitation because bad role name', async () => {
  //     await login()

  //     await runCommand(MemberInviteCommand, [
  //       'abigael@wearedevx.com',
  //       '--role=not_existing',
  //     ])

  //     const notInvited = result.find(log =>
  //       log.indexOf(
  //         'Expected --role=not_existing to be one of: reader, contributor, admin'
  //       )
  //     )
  //     expect(notInvited).toBeDefined()
  //   }, 20000)
})
