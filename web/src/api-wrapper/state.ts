import {
  type State,
  type StateDto
} from "./types/apiTypes"
import { fetchHandler } from "./fetchHandler"
import { type APIResponse } from "./types"

const STATE_ENDPOINT = "state"

export async function getState(): Promise<APIResponse<State>> {
  return await fetchHandler(`${STATE_ENDPOINT}`, undefined, undefined)
}

export async function updateState(
  updateStateDto: StateDto
): Promise<APIResponse<State>> {
  return await fetchHandler(
    STATE_ENDPOINT,
    "PATCH",
    updateStateDto
  )
}
