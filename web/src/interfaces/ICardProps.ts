export interface ICardProps<T> {
  data: T
  fetchData: () => Promise<void>
}
