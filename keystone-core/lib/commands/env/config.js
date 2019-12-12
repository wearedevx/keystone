const { updateDescriptor } = require('../../descriptor')
const { isAdminOrContributor } = require('../../member')
const config = (userSession, { project, descriptors }) => {
  return Promise.all(
    descriptors.map(({ descriptor, env }) => {
      updateDescriptor(userSession, {
        descriptorPath: descriptor.path,
        env,
        project,
        type: 'members',
        content: descriptor.content,
        name: descriptor.name,
        membersDescriptor: descriptor,
      })

      updateDescriptor(userSession, {
        // descriptorPath: descriptor.path,
        env,
        project,
        type: 'env',
        // content: descriptor.content,
        name: descriptor.name,
        membersDescriptor: descriptor,
        // content: descriptor.content
      })
    })
  )
}

module.exports = config
