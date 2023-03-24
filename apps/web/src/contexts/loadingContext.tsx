import React, {
  type Dispatch,
  type SetStateAction,
  type ReactNode,
  useState
} from "react"

export interface LoadingContextProps {
  loading: boolean
  setLoading: Dispatch<SetStateAction<boolean>>
}

interface Props {
  children: ReactNode
}

export const LoadingContext = React.createContext<LoadingContextProps>({
  loading: true,
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  setLoading: () => {}
})

export const LoadingProvider = ({ children }: Props) => {
  const [loading, setLoading] = useState(true)

  return (
    <LoadingContext.Provider
      value={{
        loading: loading,
        setLoading: setLoading
      }}
    >
      {children}
    </LoadingContext.Provider>
  )
}
