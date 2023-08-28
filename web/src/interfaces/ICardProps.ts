export interface ICardProps<T> {
  data: T
  readonly?: boolean
  fetchData: () => Promise<void>
}
