const { updateDescriptor, getOwnDescriptor } = require('../../descriptor')

const { uploadFilesForNewMembers } = require('../../env/configure')

const config = (userSession, { project, descriptors }) => {
  return Promise.all(
    descriptors.map(async ({ descriptor, env }) => {
      const updatedDescriptor = await Promise.all([
        updateDescriptor(userSession, {
          descriptorPath: descriptor.path,
          env,
          project,
          type: 'members',
          content: descriptor.content,
          name: descriptor.name,
          membersDescriptor: descriptor,
        }),
        updateDescriptor(userSession, {
          env,
          project,
          type: 'env',
          name: descriptor.name,
          membersDescriptor: descriptor,
          updateAnyway: true,
        }),
      ])

      await Promise.all(
        updatedDescriptor[1].content.files.map(async ({ name: filename }) => {
          return updateDescriptor(userSession, {
            project,
            env,
            type: 'file',
            name: filename,
            membersDescriptor: updatedDescriptor[0],
            updateAnyway: true,
          })
        })
      )

      return updatedDescriptor
    })
  )
}

module.exports = config
