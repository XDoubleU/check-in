export type DeepRequired<T> = Required<{
  [P in keyof T]-?: T[P] extends object | undefined
    ? DeepRequired<Required<T[P]>>
    : T[P]
}>

export type WithRequired<T, K extends keyof T> = T & {
  [P in K]-?: NonNullable<T[P]>
}

export type PartialBy<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>
