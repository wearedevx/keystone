const isEmail = require('is-email')
const { inviteMember } = require('../../invitation')
const { assertUserIsAdmin } = require('../../member')

const invite = async (
  userSession,
  { from, project, emails, role = 'reader' }
) => {
  // check if project exists
  await assertUserIsAdmin(userSession, { project })

  if (!isEmail(from)) {
    throw new Error(`Your email address is invalid: ${from}`)
  }

  return inviteMember(userSession, { from, project, emails, role })
}

module.exports = { invite }
