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

interface ListViewLayoutProps<
  T extends { id: string | number },
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
  apiCall: (
    ...args: W
  ) => Promise<APIResponse<WithRequired<List<X>, "data">>>
  apiCallArgs?: W
  preprocessList?: (data: WithRequired<List<X>, "data">) => Promise<U>
  createModal?: (props: ICreateModalProps<V>) => ReactNode
  card: (props: ICardProps<T>) => ReactNode
}

// eslint-disable-next-line max-lines-per-function
export default function ListViewLayout<
  T extends { id: string | number },
  U extends List<T>,
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
    const page = router.query.page
      ? parseInt(router.query.page as string)
      : undefined

    const args = apiCallArgs ?? ([] as unknown as W)
    args.push(page)

    const response = await apiCall(...args)

    if (!response.data) return

    if (
      response.data.pagination.total !== 0 &&
      response.data.pagination.current > response.data.pagination.total
    ) {
      await router.push({
        pathname: router.pathname,
        query: { page: response.data.pagination.total }
      })
    }

    const data = response.data
    if (preprocessList) {
      setList(await preprocessList(data))
    } else {
      setList(data as unknown as U)
    }
  }, [apiCall, apiCallArgs, preprocessList, setList, router])

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
        {!list.data && (
          <Loader message="Fetching data." />
        )}

        {list.data?.length == 0
          ? "Nothing to see here."
          : ""}

        {list.data?.map((item) => {
              return (
                <div key={item.id}>
                  {card({ data: item, fetchData: fetchData })}
                </div>
              )
          })
        }
      </div>

      <CustomPagination
          current={list.pagination.current}
          total={list.pagination.total}
        />
    </ManagerLayout>
  )
}
