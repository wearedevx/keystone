import React, { useState, useEffect } from 'react'
import useUser from '../hooks/useUser'
import queryString from 'query-string'
import ErrorCard from '../components/cards/error'
import BaseCard from '../components/cards/base'
import Button from '../components/button'
import ButtonWithLoggout from '../components/buttonWithLoggout'
import { ROLES } from '@keystone.sh/core/lib/constants'
import { getNameAndUUID } from '@keystone.sh/core/lib/projects'
import { acceptInvite } from '@keystone.sh/core/lib/invitation'
import WithLoggin from '../components/withLoggin'
import {
  getLatestMembersDescriptor,
  getMembers,
  getLatestProjectDescriptor,
} from '@keystone.sh/core/lib/descriptor'
import {
  setMembersToEnvs,
  isOneOrMoreAdmin,
} from '@keystone.sh/core/lib/env/configure'
import configureEnv from '@keystone.sh/core/lib/commands/env/config'
import { add } from '@keystone.sh/core/lib/commands/add'
import { getPubkey } from '@keystone.sh/core/lib/file/gaia'

import SuccessCard from '../components/cards/success'
import ProjectId from '../components/projectId'

const TitlePrompt = ({ project, env }) => (
  <SuccessCard title={`Configuring ${project}`}></SuccessCard>
)

const getProjectDetails = async (userSession, { project }) => {
  const projectDescriptor = await getLatestProjectDescriptor(userSession, {
    project,
  })

  const projectMembersDescriptor = await getLatestMembersDescriptor(
    userSession,
    { project }
  )

  const projectMembers = [
    ...projectMembersDescriptor.content[ROLES.ADMINS],
    ...projectMembersDescriptor.content[ROLES.CONTRIBUTORS],
    ...projectMembersDescriptor.content[ROLES.READERS],
  ]

  const envMembersDescriptors = await Promise.all(
    projectDescriptor.content.env.map(async env => {
      const descriptor = await getLatestMembersDescriptor(userSession, {
        project,
        env,
      })

      return { env, descriptor }
    })
  )

  const envsMembers = envMembersDescriptors.reduce(
    (acc, { env, descriptor }) => {
      acc[env] = descriptor.content
      return acc
    },
    {}
  )

  const allMembers = await getMembers(userSession, { project })
  return {
    envsMembers,
    allMembers,
    projectDescriptor,
    envMembersDescriptors,
    projectMembersDescriptor,
    projectMembers,
  }
}

const ChooseEnv = ({ setEnvironment, environments, blockstackId }) => {
  return (
    <>
      <p>
        Which environments <strong>{blockstackId}</strong> should have access
        to?
      </p>
      <div className="flex flex-row mt-4 justify-center">
        {environments.map(env => (
          <Button key={env} onClick={() => setEnvironment(env)}>
            {env}
          </Button>
        ))}
      </div>
    </>
  )
}

const ChooseRole = ({ setRole, environment, blockstackId }) => {
  return (
    <>
      <p>
        On environnment <strong>{environment}</strong>, which group{' '}
        <strong>{blockstackId}</strong> should be part of?
      </p>
      <div className="flex flex-row mt-4 justify-center">
        {Object.keys(ROLES).map(role => (
          <Button key={role} onClick={() => setRole(ROLES[role])}>
            {ROLES[role]}
          </Button>
        ))}
      </div>
    </>
  )
}

const Confirm = ({ onReset, onConfirm, role, blockstackId }) => {
  return (
    <>
      <p>
        <strong>{blockstackId}</strong> will become a member of the{' '}
        <strong>{role}</strong> group. Confirm?
      </p>
      <div className="flex flex-row mt-6 justify-center">
        <Button onClick={onConfirm}>Confirm</Button>
        <Button onClick={onReset} type="secondary">
          Reset
        </Button>
      </div>
    </>
  )
}

const PromptConfigure = ({
  project,
  projectName,
  uuid,
  blockstackId,
  email,
}) => {
  const { userSession } = useUser()
  const [error, setError] = useState(false)
  const [success, setSuccess] = useState(false)
  const [projectDetails, setProjectDetails] = useState({})
  const [environments, setEnvironments] = useState([])
  const [environment, setEnvironment] = useState(null)
  const [role, setRole] = useState(null)
  const [loading, setLoading] = useState(false)
  const [retry, setRetries] = useState(1)

  // get environnments
  useEffect(() => {
    const getData = async () => {
      setLoading('Retrieving your project... please wait...')
      try {
        const details = await getProjectDetails(userSession, { project })
        const { envsMembers } = details
        const envs = Object.keys(envsMembers)
        setEnvironments(envs)
        setProjectDetails(details)
      } catch (error) {
        setError(error.message)
      } finally {
        setLoading(false)
      }
    }

    getData()
  }, [retry])

  return (
    <>
      {error && (
        <>
          <ErrorCard title={error}>
            Please open an issue on Github if you think it's a bug.
          </ErrorCard>
          <ButtonWithLoggout
            onClick={() => {
              setError(null)
              setRetries(retry + 1)
            }}
          >
            Retry
          </ButtonWithLoggout>
        </>
      )}

      {!error && (
        <BaseCard
          title={<TitlePrompt project={projectName} env={environment} />}
        >
          {loading && <p>{loading}</p>}

          {!loading && (
            <>
              {!environment && (
                <ChooseEnv
                  blockstackId={blockstackId}
                  environments={environments}
                  setEnvironment={setEnvironment}
                />
              )}
              {environment && !role && (
                <ChooseRole
                  blockstackId={blockstackId}
                  environment={environment}
                  setRole={setRole}
                />
              )}
              {environment && role && (
                <Confirm
                  blockstackId={blockstackId}
                  environment={environment}
                  role={role}
                  onConfirm={async () => {
                    setLoading('Updating your project settings, please wait...')
                    try {
                      const {
                        envsMembers,
                        envMembersDescriptors,
                        projectMembers,
                      } = projectDetails

                      // Add user to the project
                      // if he's not already there
                      // - by default as a reader
                      const found = projectMembers.find(m => {
                        return m.blockstack_id === blockstackId
                      })

                      if (!found) {
                        await add(userSession, {
                          project,
                          invitee: {
                            blockstackId,
                            role: ROLES.READERS,
                            email,
                          },
                        })

                        const newProjectMembers = [
                          ...projectMembers,
                          { blockstack_id: blockstackId },
                        ]
                        // avoid duplicates at the project level
                        setProjectDetails({
                          ...projectDetails,
                          projectMembers: newProjectMembers,
                        })
                      }

                      const envMembersDescriptor = envMembersDescriptors.find(
                        envDescriptor => envDescriptor.env === environment
                      )
                      const newEnvsMembers = setMembersToEnvs({
                        envsMembers,
                        members: [
                          ...envsMembers[environment][role],
                          { blockstack_id: blockstackId },
                        ],
                        role,
                        env: environment,
                      })

                      if (!isOneOrMoreAdmin(newEnvsMembers)) {
                        setError(
                          'One or more admin member is required for this environment'
                        )
                      }

                      envMembersDescriptor.descriptor.content =
                        newEnvsMembers[environment]

                      // const newEnvMembersDescriptor = await Promise.all(
                      // envMembersDescriptor.map(async envDescriptor => {
                      // const envDescriptorClone = { ...envDescriptor }

                      for (const role of Object.values(ROLES)) {
                        const members = await Promise.all(
                          envMembersDescriptor.descriptor.content[role].map(
                            async member => {
                              if (!member.publicKey) {
                                member.publicKey = await getPubkey(
                                  userSession,
                                  member
                                )
                              }

                              return member
                            }
                          )
                        )

                        envMembersDescriptor.descriptor.content[role] = members
                      }

                      // return envDescriptorClone
                      //   })
                      // )

                      await configureEnv(userSession, {
                        project,
                        descriptors: [envMembersDescriptor],
                      })
                      setSuccess(
                        `${blockstackId} can pull files from environment ${environment}`
                      )
                      setRole(null)
                      setEnvironment(null)
                    } catch (error) {
                      console.error(error)
                      setError(error.message)
                    } finally {
                      setLoading(false)
                    }
                  }}
                  onReset={() => {
                    setRole(null)
                    setEnvironment(null)
                  }}
                />
              )}
            </>
          )}

          <ProjectId>{uuid}</ProjectId>
        </BaseCard>
      )}

      {success && (
        <>
          <SuccessCard title={`Project updated successfully`}>
            {success}
          </SuccessCard>
          <ButtonWithLoggout onClick={() => setSuccess(null)}>
            Continue configuration
          </ButtonWithLoggout>
        </>
      )}
    </>
  )
}

export default () => {
  const { project, id, email } =
    (typeof location !== 'undefined' && queryString.parse(location.search)) ||
    {}
  let missingParams = !project || !id || !email
  let projectName,
    projectUUID = null
  try {
    ;[projectName, projectUUID] = getNameAndUUID(project)
  } catch (error) {
    missingParams = true
  }

  return (
    <WithLoggin>
      {missingParams && (
        <ErrorCard
          title={'Your link is malformed. Please open an issue on GitHub.'}
        >
          Or check that the link in your browser is the same than the link you
          received in your mailbox.
        </ErrorCard>
      )}

      {!missingParams && (
        <PromptConfigure
          project={project}
          projectName={projectName}
          uuid={projectUUID}
          blockstackId={id}
          email={decodeURIComponent(email)}
        />
      )}
    </WithLoggin>
  )
}
