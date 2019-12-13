import React from 'react'
import useUser from '../hooks/useUser'
import invitations from '@keystone/core/dist/invitation'
console.log('TCL: invitations', invitations)

export default () => {
  const { loggedIn, userData } = useUser()
  console.log('TCL: loggedIn', loggedIn)
  return (
    <div className="flex flex-col items-center ">
      <div className="shadow-md rounded p-4 bg-white w-2/4">
        <h1 className="text-xl">
          <span
            role="img"
            aria-label="A party popper, as explodes in a shower of confetti and streamers at a celebration"
          >
            🎉
          </span>
          You've been invited by ___ to join ___
        </h1>
      </div>
      <div className="my-4 flex flex-row w-2/4 justify-end">
        <div className="rounded font-bold text-white bg-primary py-1 px-4 shadow-md text-center mx-2">
          Join
        </div>
        <div className="rounded font-bold text-gray-300 bg-secondary py-1 px-4 shadow-md text-center">
          Decline
        </div>
      </div>
      {/* <InvitationBannerJoin project={project} from={from} id={id} /> */}
      {/* <InvitationBoard project={project} /> */}
    </div>
  )
}
