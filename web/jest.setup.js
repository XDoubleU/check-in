import '@testing-library/jest-dom'
import { mocked } from "jest-mock"
import { noUserMock, DefaultLocation } from "mocks"
import { getMyUser, getDataForRangeChart, getAllLocations } from "api-wrapper"

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

mocked(getAllLocations).mockImplementation(() => {
  return Promise.resolve({
    ok: true,
    data: [DefaultLocation]
  })
})

mocked(getDataForRangeChart).mockImplementation(() => {
  return Promise.resolve({
    ok: true,
    data: {
      "2023-08-24": {
        capacity: 10,
        schools: {
          "Andere": 5
        }
      }
    }
  })
})
