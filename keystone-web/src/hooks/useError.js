import React, { useState, useEffect } from "react"
import useStore from "../store"

export const useError = () => {
  const { error, setError } = useStore(state => ({
    setError: state.setError,
  }))

  return { setError, error }
}
