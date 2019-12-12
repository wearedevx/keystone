const {
  getLastEnvDescriptor,
  incrementVersion,
  uploadDescriptorForEveryone,
} = require('../../env')
const { deepCopy } = require('../../../utils')

const deleteFiles = async (userSession, { project, env, files }) => {
  const { username } = userSession.loadUserData()

  const latestEnvDescriptor = await getLastEnvDescriptor(userSession, {
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

  const newEnvDescriptor = incrementVersion({
    descriptor: envDescriptorCloned,
    author: username,
    previousDescriptor: latestEnvDescriptor,
    type: 'env',
  })

  return uploadDescriptorForEveryone(userSession, {
    type: 'env',
    descriptor: newEnvDescriptor,
    envDescriptor: newEnvDescriptor,
  })
}

module.exports = deleteFiles
