export interface APIResponse<T = undefined> {
  ok: boolean
  message?: unknown
  data?: T
}
