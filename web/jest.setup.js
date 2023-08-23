import { signOut, signIn, getMyUser } from "api-wrapper"
import '@testing-library/jest-dom'

jest.mock('next/router', () => require('next-router-mock'))

jest.mock('next/head', () => {
  return {
    __esModule: true,
    default: ({ children }) => {
      return <>{children}</>;
    },
  }
})

jest.mock("./src/api-wrapper")
signOut.mockImplementation(() => Promise.resolve(undefined))
signIn.mockImplementation((signInDto) => {
  if(signInDto.username === "validusername" && signInDto.password === "validpassword") {
    return Promise.resolve({
      ok: true,
      data: {
        username: "validusername"
      }
    })
  }

  return Promise.resolve({
    ok: false,
    message: "Invalid credentials"
  })
})

getMyUser.mockImplementation(() => Promise.resolve({
  ok: true,
  data: {
    username: "validusername"
  }
}))