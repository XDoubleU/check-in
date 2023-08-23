import '@testing-library/jest-dom'
import { mocked } from "jest-mock"
import { noUserMock } from "user-mocks"
import { getMyUser } from "api-wrapper"

jest.mock('next/router', () => require('next-router-mock'))

jest.mock('next/head', () => {
  return {
    __esModule: true,
    default: ({ children }) => {
      return <>{children}</>;
    },
  }
})

jest.mock("api-wrapper")
mocked(getMyUser).mockImplementation(noUserMock)