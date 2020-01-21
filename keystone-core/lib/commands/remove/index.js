const { deepCopy } = require('../../utils')

const {
  getLatestMembersDescriptor,
  updateDescriptor,
} = require('../../descriptor')
const { ROLES } = require('../../constants')

const removeFromProject = async (userSession, { project, user }) => {
  const membersDescriptor = await getLatestMembersDescriptor(userSession, {
    project,
  })

  const members = deepCopy(membersDescriptor.content)

  Object.values(ROLES).forEach(role => {
    const userIndex = members[role].findIndex(u => u.blockstack_id === user)
    if (userIndex >= 0) {
      members[role].splice(userIndex, 1)
    }
  })

  await updateDescriptor(userSession, {
    type: 'members',
    project,
    content: members,
    membersDescriptor,
  })
  return { removed: true }
}

module.exports = {
  removeFromProject,
}
