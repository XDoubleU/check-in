import '@testing-library/jest-dom'
import { mocked } from "jest-mock"
import { noUserMock } from "mocks"
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

delete window.ResizeObserver;
window.ResizeObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));

jest.mock("api-wrapper")
mocked(getMyUser).mockImplementation(noUserMock)

