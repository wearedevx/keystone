const { updateDescriptor } = require('../../descriptor')

const config = (userSession, { project, descriptor }) => {
  return updateDescriptor(userSession, {
    descriptorPath: descriptor.path,
    project,
    type: 'project',
    content: descriptor.content,
    name: descriptor.name,
  })
}

module.exports = config
