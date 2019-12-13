import React from 'react'
import { Root, Routes } from 'react-static'
import { Router } from './components/Router'
import Logo from './components/logo'
import GithubIcon from './svgs/github'
import useBlockstack from './hooks/useBlockstack'
// import { descriptor } from '@keystone/core'
// console.log('TCL: kcore', descriptor)

import './app.css'

function App() {
  useBlockstack()
  return (
    <Root>
      <nav className="p-2 flex flex-col lg:flex-row lg:p-4 self-center items-center">
        <div className="flex-grow flex self-center items-center">
          <Logo scale={0.8} />
        </div>
        <div className="flex flex-col lg:flex-row items-center">
          <h1 className=" text-gray-800 text-sm text-center mb-6 mt-2 lg:text-left lg:mb-0 lg:mt-0">
            <span>Open-source decentralized system </span>
            <span className="font-bold inline-block lg:inline">
              for developers to store, share and use secrets safely
            </span>
            .
          </h1>
          <div className="px-2 text-2xl mb-6 lg:mr-4 lg:mb-0 lg:px-4">
            <GithubIcon />
          </div>
        </div>
        {/* <Link to="/">Home</Link>
        <Link to="/about">About</Link>
        <Link to="/blog">Blog</Link>
        <Link to="/dynamic">Dynamic</Link> */}
      </nav>
      <div className="flex flex-col items-center mx-4 lg:mx-10">
        <div className="h-1 bg-softline w-full mb-10" />
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
