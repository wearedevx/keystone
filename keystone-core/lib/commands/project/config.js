const { updateDescriptor } = require('../../descriptor')

const config = (userSession, { project, descriptor }) => {
  console.log("TCL: config -> descriptor", descriptor)
  return updateDescriptor(userSession, {
    descriptorPath: descriptor.path,
    project,
    type: 'members',
    content: descriptor.content,
    name: descriptor.name,
  })
}

module.exports = config
