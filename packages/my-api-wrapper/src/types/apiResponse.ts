export default interface APIResponse<T = undefined> {
  ok: boolean
  message?: string
  data?: T
}
