const {
  getLatestProjectDescriptor,
  uploadDescriptorForEveryone,
} = require('../../descriptor')
const { addMemberToProject } = require('../../invitation')
const { assertUserIsAdmin } = require('../../member')

const add = async (userSession, { project, invitee }) => {
  await assertUserIsAdmin(userSession, { project })

  const projectDescriptor = await getLatestProjectDescriptor(userSession, {
    project,
  })

  const memberAdded = await addMemberToProject(userSession, {
    project,
    invitee: {
      blockstack_id: invitee.blockstackId,
      email: invitee.email,
      role: invitee.role,
    },
  })

  // updateDescriptor(userSession, { project })

  await uploadDescriptorForEveryone(userSession, {
    members: [invitee.blockstackId],
    descriptor: projectDescriptor,
    type: 'project',
  })

  return { added: true, memberAdded }
}

module.exports = { add }
