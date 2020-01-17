const fs = require('fs')
const uuid = require('uuid/v4')
// const { ERROR_CODES } = require('../../constants')

const KeystoneError = require('../../error')
const {
  getProjects,
  findProjectByName,
  findProjectByUUID,
  createProject,
} = require('../../projects')
const { writeFileToDisk, getKeystoneFolder } = require('../../file')
const { createEnv } = require('../../env')

const {
  getDescriptor,
  updateDescriptorForMembers,
  uploadDescriptorForEveryone,
  extractMembersByRole,
} = require('../../descriptor')

const { createMembersDescriptor } = require('../../member')

const { KEYSTONE_ENV_CONFIG_PATH } = require('../../constants')

const init = async (userSession, { project, overwrite = false }) => {
  try {
    const projects = await getProjects(userSession)

    // is there a project with the same name?
    const projectExists = findProjectByName(projects, project)
    if (projectExists.length > 0) {
      throw new KeystoneError(
        'ProjectNameExists',
        'One or more projects have the same name',
        projectExists
      )
    }

    let projectFullname = `${project}/${uuid()}`
    let createDefault = true
    const ksconfigDescriptor = {
      name: '.ksconfig',
      content: {
        project: projectFullname,
      },
    }

    const envConfigDescriptor = {
      name: KEYSTONE_ENV_CONFIG_PATH,
      content: {
        env: 'default',
      },
    }
    const projectWithUUID = findProjectByUUID(projects, project)
    if (projectWithUUID) {
      projectFullname = project
      ksconfigDescriptor.content = {
        project,
      }
      createDefault = false
    }

    if (!fs.existsSync('.ksconfig') || overwrite) {
      await writeFileToDisk(ksconfigDescriptor, './')
      await writeFileToDisk(envConfigDescriptor, getKeystoneFolder('.'))

      if (createDefault) {
        const projectDescriptor = await createProject(userSession, {
          name: projectFullname,
        })

        await createMembersDescriptor(userSession, { project: projectFullname })

        // create default environnment
        await createEnv(userSession, {
          env: 'default',
          projectDescriptor,
          // members,
        })
      } else {
        try {
          const projectMembersDescriptor = await getDescriptor(userSession, {
            project,
            type: 'members',
            origin: projectWithUUID.createdBy,
          })
          await uploadDescriptorForEveryone(userSession, {
            members: extractMembersByRole(projectMembersDescriptor, [
              'admins',
              'contributors',
              'readers',
            ]),
            descriptor: projectMembersDescriptor,
            type: 'members',
          })
          const projectDescriptor = await getDescriptor(userSession, {
            project,
            type: 'project',
            origin: projectWithUUID.createdBy,
          })
          await uploadDescriptorForEveryone(userSession, {
            members: extractMembersByRole(projectMembersDescriptor, [
              'admins',
              'contributors',
              'readers',
            ]),
            descriptor: projectDescriptor,
            type: 'project',
          })
        } catch (err) {
          console.log(err)
          throw new KeystoneError(
            'FailedToFetch',
            'Cannot fetch descriptor from project host.',
            projectWithUUID.createdBy
          )
        }
      }
    } else {
      throw new KeystoneError(
        'ConfigFileExists',
        'Config file .ksconfig already exist'
      )
    }
    return project
  } catch (error) {
    throw error
  }
}

module.exports = init
