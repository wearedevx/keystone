const treeify = require('treeify')

const { getLatestMembersDescriptor } = require('../../descriptor')

const listAllMembers = async (userSession, { project }) => {
  const projectMembersDescriptor = await getLatestMembersDescriptor(
    userSession,
    {
      project,
    }
  )

  console.log(treeify.asTree(projectMembersDescriptor.content, true))
}

const listEnvMembers = async (userSession, { project, env }) => {
  const envMembersDescriptor = await getLatestMembersDescriptor(userSession, {
    project,
    env,
  })

  console.log(treeify.asTree(envMembersDescriptor.content, true))
}

module.exports = {
  listAllMembers,
  listEnvMembers,
}
