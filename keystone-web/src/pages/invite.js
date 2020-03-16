import React, { useState } from 'react'
import useUser from '../hooks/useUser'
import queryString from 'query-string'
import ErrorCard from '../components/cards/error'
import BaseCard from '../components/cards/base'
import Button from '../components/button'
import { getNameAndUUID } from '@keystone.sh/core/lib/projects'
import { acceptInvite } from '@keystone.sh/core/lib/invitation'
import WithLoggin from '../components/withLoggin'

const TitlePromptInvite = ({ project }) => (
  <>
    You are invited to project <strong>{project}</strong>.
  </>
)

const join = async (
  userSession,
  { name, from, blockstackId, userEmail, onError, onDone, onSuccess }
) => {
  try {
    const projects = await acceptInvite(userSession, {
      name,
      from,
      blockstackId,
      userEmail,
    })
    onSuccess(projects)
    return projects
  } catch (error) {
    onError(error)
  } finally {
    onDone()
  }
}

const PromptInvite = ({
  project,
  uuid,
  adminUserEmail,
  blockstackId,
  userEmail,
}) => {
  const { userSession, signUserOut } = useUser()
  const [joining, setJoining] = useState(false)
  const [error, setError] = useState(false)
  const [success, setSuccess] = useState(false)

  const { username = 'username missing' } = userSession.loadUserData()

  return (
    <>
      {!error && !success && (
        <>
          <BaseCard title={<TitlePromptInvite project={project} />}>
            {joining && <p>Updating your projects list, please wait...</p>}

            {!joining && (
              <p className="text-gray-800">
                This invite is sent by{' '}
                <span className="font-bold text-sm font-mono text-black">
                  {adminUserEmail}
                </span>
                .<br /> Click join to accept or ignore if you don't know the
                sender.
              </p>
            )}

            <p className="font-mono text-blue-700 mt-6 text-xs">
              <span className="font-bold">Project ID:</span> {uuid}
            </p>
          </BaseCard>
          <div className="my-4 flex flex-row w-2/4 justify-center">
            <div className="flex flex-col items-center">
              <Button
                disabled={joining || success}
                onClick={async () => {
                  setJoining(true)
                  await join(userSession, {
                    name: `${project}/${uuid}`,
                    from: adminUserEmail,
                    blockstackId,
                    userEmail,
                    onError: e => setError(e.message),
                    onDone: () => setJoining(false),
                    onSuccess: () => setSuccess(true),
                  })
                }}
              >
                Accept and join "{project}"
              </Button>
              <div className="text-xs text-gray-600 mt-1 flex flex-col items-center">
                <div>connected with {username}</div>
                <div
                  className="underline font-bold cursor-pointer"
                  onClick={() => signUserOut()}
                >
                  Log out
                </div>
              </div>
            </div>
          </div>
        </>
      )}

      {error && (
        <ErrorCard title={error}>
          Refresh and retry or open an issue on Github.
        </ErrorCard>
      )}

      {success && (
        <SuccessCard title={`An email has been sent to ${adminUserEmail}`}>
          This user will confirm your membership and encrypt the projects files
          for you.
        </SuccessCard>
      )}
    </>
  )
}

export default () => {
  const { project, id, from, to } =
    (typeof location !== 'undefined' && queryString.parse(location.search)) ||
    {}
  let missingParams = !project || !id || !from || !to
  let projectName,
    projectUUID = null
  try {
    ;[projectName, projectUUID] = getNameAndUUID(project)
  } catch (error) {
    missingParams = true
  }

  return (
    <div className="flex flex-col items-center ">
      <WithLoggin redirectURI="/invite">
        {missingParams && (
          <ErrorCard
            title={'Your link is malformed. Please open an issue on GitHub.'}
          >
            Or check that the link in your browser is the same than the link you
            received in your mailbox.
          </ErrorCard>
        )}

        {!missingParams && (
          <PromptInvite
            project={projectName}
            uuid={projectUUID}
            adminUserEmail={decodeURIComponent(from)}
            userEmail={decodeURIComponent(to)}
            blockstackId={id}
          />
        )}
      </WithLoggin>
    </div>
  )
}
