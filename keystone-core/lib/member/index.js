const {
  getLatestMembersDescriptor,
  updateDescriptor,
  updateDescriptorForMembers,
} = require('../descriptor')
const KeystoneError = require('../error')
const { ROLES, ERROR_CODES, SHARED_MEMBER } = require('../constants')
const { getPubkey } = require('../file/gaia')

const doesUserHasRole = async (userSession, { project, env }, roles) => {
  const { username } = userSession.loadUserData()
  const memberDescriptor = await getLatestMembersDescriptor(userSession, {
    project,
    env,
  })

  return roles.reduce((hasRole, role) => {
    return (
      hasRole ||
      memberDescriptor.content[role].find(a => a.blockstack_id === username)
    )
  }, false)
}

const assertUserIsAdmin = async (userSession, { project, env }) => {
  if (!(await doesUserHasRole(userSession, { project, env }, [ROLES.ADMINS])))
    throw new KeystoneError('NeedToBeAdmin')
}

const assertUserIsAdminOrContributor = async (
  userSession,
  { project, env }
) => {
  if (
    !(await doesUserHasRole(userSession, { project, env }, [
      ROLES.ADMINS,
      ROLES.CONTRIBUTORS,
    ]))
  ) {
    throw new KeystoneError('NeedToBeAdminOrContributor')
  }
}

const createMembersDescriptor = async (userSession, { project, env }) => {
  const { username } = userSession.loadUserData()

  const publicKey = await getPubkey(userSession, { blockstack_id: username })

  const members = {
    [ROLES.ADMINS]: [{ blockstack_id: username, publicKey }],
    [ROLES.CONTRIBUTORS]: [],
    [ROLES.READERS]: [],
  }

  const membersDescriptor = { content: members }

  return updateDescriptorForMembers(userSession, {
    env,
    project,
    type: 'members',
    membersDescriptor,
    content: members,
    name: 'members',
  })
}

const addMember = async (
  userSession,
  { project, env, member, role, publicKey }
) => {
  if (!publicKey) {
    publicKey = await getPubkey(userSession, { blockstack_id: member })
  }
  const membersDescriptor = await getLatestMembersDescriptor(userSession, {
    project,
    env,
  })

  const allMembers = Object.values(ROLES).reduce((members, r) => {
    return [...members, ...membersDescriptor.content[r]]
  }, [])

  if (
    member !== SHARED_MEMBER &&
    allMembers.find(m => m.blockstack_id === member)
  ) {
    throw new KeystoneError(
      ERROR_CODES.InvitationFailed,
      'User already in the project',
      member
    )
  }

  const newMember = { blockstack_id: member, publicKey }
  membersDescriptor.content[role] = membersDescriptor.content[role].filter(
    m => m.blockstack_id !== SHARED_MEMBER
  )

  membersDescriptor.content[role].push(newMember)

  const newMembersDescriptor = await updateDescriptor(userSession, {
    project,
    env,
    type: 'members',
    name: 'members',
    content: membersDescriptor.content,
    membersDescriptor,
  })

  // If the member is added to an env
  if (env) {
    const envDescriptor = await updateDescriptor(userSession, {
      project,
      env,
      type: 'env',
      name: env,
      membersDescriptor: newMembersDescriptor,
      updateAnyway: true,
    })

    await Promise.all(
      envDescriptor.content.files.map(async ({ name: filename }) => {
        return updateDescriptor(userSession, {
          project,
          env,
          type: 'file',
          name: filename,
          membersDescriptor: newMembersDescriptor,
          updateAnyway: true,
        })
      })
    )
  }
  return { member, membersDescriptor }
}

module.exports = {
  createMembersDescriptor,
  addMember,
  assertUserIsAdminOrContributor,
  assertUserIsAdmin,
}
