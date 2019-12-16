const {
  getLatestEnvDescriptor,
  incrementVersion,
  updateDescriptor,
} = require('../../descriptor')

const { deepCopy } = require('../../utils')

const deleteFiles = async (userSession, { project, env, files }) => {
  const { username } = userSession.loadUserData()

  const latestEnvDescriptor = await getLatestEnvDescriptor(userSession, {
    project,
    env,
    username,
  })
  const envDescriptorCloned = deepCopy(latestEnvDescriptor)

  envDescriptorCloned.content.files = latestEnvDescriptor.content.files.filter(
    envFile => !files.find(file => file === envFile.name)
  )

  if (
    envDescriptorCloned.content.files.length ===
    latestEnvDescriptor.content.files.length
  ) {
    throw new Error('No file to delete in keystone')
  }

  const updatedDescriptor = await updateDescriptor(userSession, {
    env,
    project,
    content: envDescriptorCloned.content,
    type: 'env',
  })
  return updatedDescriptor
}

module.exports = deleteFiles
