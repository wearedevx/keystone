import React from 'react'
import { Root, Routes } from 'react-static'
import { Router } from './components/Router'
import Logo from './components/logo'
import GithubIcon from './svgs/github'
import useBlockstack from './hooks/useBlockstack'
// import { descriptor } from '@keystone.sh/core'

import './app.css'

function App() {
  useBlockstack()
  return (
    <Root>
      <nav className="p-2 flex flex-col lg:flex-row lg:p-4 lg:mx-6 self-center items-center">
        <div className="flex-grow flex self-center items-center">
          <Logo />
        </div>
        <div className="flex flex-col lg:flex-row items-center">
          <h1 className=" text-gray-900 text-sm text-center mb-6 mt-2 lg:text-left lg:mb-0 lg:mt-0 lg:ml-12 font-mono">
            <span>Open-source decentralized system </span>
            <span className="font-bold inline-block lg:inline">
              for developers to store, share and use secrets.
            </span>
            .
          </h1>
          <div className="px-2 text-2xl mb-6 lg:mb-0">
            <a href="https://github.com/wearedevx/keystone" target="_blank">
              <GithubIcon />
            </a>
          </div>
        </div>
      </nav>
      <div className="flex flex-col items-center mx-4 lg:mx-10">
        <div className="h-1 bg-gray-300 w-full mb-10" />
      </div>
      <div className="mx-4 mb-40 lg:mx-10">
        <React.Suspense fallback={<em>Loading...</em>}>
          <Router>
            <Routes path="*" />
          </Router>
        </React.Suspense>
      </div>
    </Root>
  )
}

export default App
