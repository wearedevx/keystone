const { updateDescriptor, getOwnDescriptor } = require('../../descriptor')

const { uploadFilesForNewMembers } = require('../../env/configure')

const config = (userSession, { project, descriptors }) => {
  return Promise.all(
    descriptors.map(async ({ descriptor, env }) => {
      const { content: members } = await getOwnDescriptor(userSession, {
        project,
        env,
        type: 'members',
      })
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

      // const previousMembers = members.reduce((members, role) => {
      //   members.push(...role)
      // }, [])
      // const currentMembers = descriptor.content.reduce((members, role) => {
      //   members.push(...role)
      // }, [])

      // const newMembers = !previousMembers.find(m =>
      //   currentMembers.find(me => me.blockstack_id === m.blockstack_id)
      // )

      // await uploadFilesForNewMembers(userSession, {
      //   project,
      //   env,
      //   members: newMembers,
      // })
      return updatedDescriptor
    })
  )
}

module.exports = config
