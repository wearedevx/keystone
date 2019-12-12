const axios = require('axios')
const isEmail = require('is-email')

const { getProjects, updateProjectsFromStatuses } = require('../project')
const { updateWorkspace, getPath, getFileFromGaia, writeFileOnGaia } = require('../file')

const KEYSTONE_MAIL = 'http://localhost:8080'
// const KEYSTONE_INVITE_URL = 'http://localhost:8000/invite'
const INVITATIONS_STORE = 'invitations.json'

const createInvitationsFile = async userSession => {
  try {
    await writeFileOnGaia(userSession, INVITATIONS_STORE, JSON.stringify([]))
    return []
  } catch (error) {
    throw error
  }
}

const getInvitations = async userSession => {
  try {
    const invitationsFile = await getFileFromGaia(
      userSession,
      INVITATIONS_STORE
    )
    if (invitationsFile) {
      return JSON.parse(invitationsFile)
    }
    throw new Error(`Invitations file not found`)
  } catch (error) {
    // if invitations not found, create the file
    const invitationsFile = await createInvitationsFile(userSession)
    return invitationsFile
  }
}

const addInvitee = async (userSession, { project, email }) => {
  // check if invitations.json exists
  try {
    let invitations = await getInvitations(userSession)
    const userData = userSession.loadUserData()

    // does the invitationn already exists?
    const alreadyInvited = invitations.findIndex(invitation => {
      return invitation.email === email && invitation.project === project
    })

    if (alreadyInvited > -1) {
      throw new Error(`${email} has already been invited to ${project}`)
    }

    invitations = [
      ...invitations,
      {
        email,
        project,
        blockstack_id: userData.username,
      },
    ]

    await updateInvitations(userSession, { invitations })
  } catch (error) {
    throw error
  }
}

const invite = async (
  userSession,
  { from, project, emails, role = 'reader' }
) => {
  // check if project exists
  try {
    // const projectDescriptor = await getProjectDescriptor(userSession, {
    //   project,
    // })
    const userData = userSession.loadUserData()
    // TODO : check if user is admin ??

    const isAdmin = true
    const fromEmailIsValid = isEmail(from)

    if (!fromEmailIsValid) {
      throw new Error(`Your email address is invalid: ${from}`)
    }

    if (isAdmin) {
      const emailsSent = await Promise.all(
        emails.map(async email => {
          // loosely check if format is an email
          if (isEmail(email)) {
            // const link = `${KEYSTONE_INVITE_URL}?from=${from}&id=${userData.username}&project=${project}`
            // console.log("should send email to", email, "with", link)
            try {
              await addInvitee(userSession, { project, email, role })

              await axios({
                method: 'post',
                url: KEYSTONE_MAIL,
                data: {
                  request: 'invite',
                  from,
                  id: userData.username,
                  project,
                  role,
                  email,
                },
              })

              return {
                email,
                sent: true,
              }
            } catch (error) {
              return {
                email,
                sent: false,
                error: error.message,
              }
            }
          } else {
            return {
              email,
              sent: false,
              error: `Email address ${email} is invalid`,
            }
          }
        })
      )

      return emailsSent
    }
    throw new Error(
      `You need to be an administrator of ${project} to invite people.`
    )
  } catch (error) {
    throw error
  }
}

const deleteInvites = async (userSession, { project, emails }) => {
  // check if invitations.json exists
  try {
    let invitations = await getInvitations(userSession)
    let deleted = []

    invitations = invitations.reduce((invites, invite) => {
      const foundInvite = emails.find(email => {
        return email == invite.email && project == invite.project
      })

      if (foundInvite) {
        deleted = [
          ...deleted,
          invitations.find(invite => {
            return invite.email == foundInvite && project == project
          }),
        ]
      }

      if (!foundInvite) {
        return [...invites, invite]
      }
      return invites
    }, [])

    await updateInvitations(userSession, { invitations })
    return deleted
  } catch (error) {
    throw error
  }
}

const updateInvitations = async (userSession, { invitations }) => {
  try {
    await writeFileOnGaia(
      userSession,
      INVITATIONS_STORE,
      JSON.stringify(invitations)
    )
    return []
  } catch (error) {
    throw error
  }
}

const acceptInvite = async (userSession, { invite }) => {
  const projects = await getProjects(userSession)
  const userData = userSession.loadUserData()
  const hasProject = projects.find(project => project.name === invite.name)

  if (hasProject) {
    throw new Error(
      `Project ${hasProject.name} is already in the user workspace`
    )
  }

  const projectsUpdated = await updateWorkspace(userSession, {
    name: invite.name,
    blockstack_id: invite.from,
    invitation: true,
    at: [invite.from],
  })

  // send an email to let the invitation creator know that the user accepted
  await axios({
    method: 'post',
    url: KEYSTONE_MAIL,
    data: {
      request: 'accept',
      to: invite.from,
      id: userData.username,
      project: invite.name,
      // role,
      // email
    },
  })

  return projectsUpdated
}



const checkInvitations = async userSession => {
  // get the projects
  const projects = await getProjects(userSession)
  const userData = userSession.loadUserData()
  const projectsStatuses = await Promise.all(
    projects.map(async project => {
      // check if invite still pending
      // TODO : remove true...
      if (project.pendingInvite || true) {
        const path = getPath({
          project: project.name,
          filename: null,
          blockstack_id: userData.username,
          type: 'project',
        })
        try {
          // great we have access!
          // TODO: instead of takin createdBy, use the at that references all contributors
          const projectDescriptor = await getFileFromGaia(userSession, path, {
            username: project.createdBy,
          })

          return {
            invite: 'fulfilled',
            project,
            projectDescriptor: JSON.parse(projectDescriptor),
          }
        } catch (error) {
          // too bad, still not validated or somehting happened
          return {
            invite: 'pending',
            project,
            error: error.message,
          }
        }
      } else {
        return {
          invite: 'fulfilled',
          project,
        }
      }
    })
  )

  const projectsUpdated = await updateProjectsFromStatuses(userSession, {
    projects,
    projectsStatuses,
  })
  return projectsUpdated
}

module.exports = {
  invite,
  deleteInvites,
  getInvitations,
  checkInvitations,
  acceptInvite,
}
