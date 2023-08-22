import CustomPagination, {
  type CustomPaginationProps
} from "components/CustomPagination"
import Loader from "components/Loader"
import { type APIResponse } from "api-wrapper"
import { useRouter } from "next/router"
import {
  type ReactNode,
  useCallback,
  useEffect,
  type Dispatch,
  type SetStateAction
} from "react"
import { type FieldValues, type UseFormReturn } from "react-hook-form"
import ManagerLayout from "./ManagerLayout"
import { type ICardProps } from "interfaces/ICardProps"
import { type ICreateModalProps } from "interfaces/ICreateModalProps"
import { type WithRequired } from "typing-helpers"

export interface List<T> {
  data: T[] | undefined
  pagination: CustomPaginationProps
}

function isList<T>(list: List<T> | T[]): list is List<T> {
  return "pagination" in list
}

interface ListViewLayoutProps<
  T extends { id: string | number },
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  U extends List<T> | T[],
  V extends FieldValues,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  W extends any[],
  X
> {
  title: string
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  form?: UseFormReturn<any>
  list: U
  setList: Dispatch<SetStateAction<U>>
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  apiCall: (
    ...args: W
  ) => Promise<APIResponse<WithRequired<List<X>, "data"> | X[]>>
  apiCallArgs?: W
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  preprocessList?: (data: WithRequired<List<X>, "data">) => Promise<U>
  createModal?: (props: ICreateModalProps<V>) => ReactNode
  card: (props: ICardProps<T>) => ReactNode
}

// eslint-disable-next-line max-lines-per-function
export default function ListViewLayout<
  T extends { id: string | number },
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  U extends List<T> | T[],
  V extends FieldValues,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  W extends any[],
  X
>({
  title,
  form,
  list,
  setList,
  createModal,
  apiCall,
  apiCallArgs,
  preprocessList,
  card
}: ListViewLayoutProps<T, U, V, W, X>) {
  const router = useRouter()

  const fetchData = useCallback(async () => {
    if (!router.isReady) return

    const page = router.query.page
      ? parseInt(router.query.page as string)
      : undefined

    const args = apiCallArgs ?? ([] as unknown as W)

    if (isList(list)) {
      args.push(page)
    }

    // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
    const response = await apiCall(...args)

    if (!response.data) return

    if (
      isList(response.data) &&
      response.data.pagination.total !== 0 &&
      response.data.pagination.current > response.data.pagination.total
    ) {
      await router.push({
        pathname: router.pathname,
        query: { page: response.data.pagination.total }
      })
    }

    const data = response.data
    if (isList(data) && preprocessList) {
      setList(await preprocessList(data))
    } else {
      setList(data as unknown as U)
    }
  }, [apiCall, apiCallArgs, list, preprocessList, router, setList])

  useEffect(() => {
    void fetchData()
  }, [fetchData])

  return (
    <ManagerLayout
      title={title}
      titleButton={form ? createModal?.({ form, fetchData }) : undefined}
    >
      <br />

      <div className="min-vh-51">
        {!((isList(list) && list.data) || list) && (
          <Loader message="Fetching data." />
        )}

        {(isList(list) && list.data?.length == 0) ||
        (!isList(list) && list.length == 0)
          ? "Nothing to see here."
          : ""}

        {isList(list)
          ? list.data?.map((item) => {
              return (
                <div key={item.id}>
                  {card({ data: item, fetchData: fetchData })}
                </div>
              )
            })
          : list.map((item) => {
              return (
                <div key={item.id}>
                  {card({ data: item, fetchData: fetchData })}
                </div>
              )
            })}
      </div>

      {isList(list) ? (
        <CustomPagination
          current={list.pagination.current}
          total={list.pagination.total}
        />
      ) : (
        <></>
      )}
    </ManagerLayout>
  )
}
