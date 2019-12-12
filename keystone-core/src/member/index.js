const {
  getLatestMembersDescriptor,
  updateDescriptor,
  updateDescriptorForMembers,
} = require('../descriptor')
const KeystoneError = require('../error')
const { ROLES, ERROR_CODES } = require('../constants')

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

const createMembersDescriptor = (userSession, { project, env }) => {
  const { username } = userSession.loadUserData()

  const members = {
    [ROLES.ADMINS]: [{ blockstack_id: username }],
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

const addMember = async (userSession, { project, env, member, role }) => {
  const membersDescriptor = await getLatestMembersDescriptor(userSession, {
    project,
    env,
  })

  const allMembers = Object.keys(ROLES).reduce((members, r) => {
    return [...members, membersDescriptor.content[r]]
  }, [])

  if (allMembers.find(m => m === member)) {
    throw new KeystoneError(
      ERROR_CODES.InvitationFailed,
      'User already in the project',
      member
    )
  }

  membersDescriptor.content[role].push({ blockstack_id: member })

  return updateDescriptor(userSession, {
    project,
    env,
    type: 'members',
    name: 'members',
    content: membersDescriptor.content,
    membersDescriptor,
  })
}

module.exports = {
  createMembersDescriptor,
  addMember,
  assertUserIsAdminOrContributor,
  assertUserIsAdmin,
}
