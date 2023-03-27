export interface ICardProps<T> {
  data: T
  refetchData: () => Promise<void>
}
