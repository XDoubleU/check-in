import CustomPagination, {
  type CustomPaginationProps
} from "../components/CustomPagination"
import Loader from "../components/Loader"
import type APIResponse from "my-api-wrapper/dist/src/types/apiResponse"
import { useRouter } from "next/router"
import {
  type ReactNode,
  useCallback,
  useEffect,
  type Dispatch,
  type SetStateAction
} from "react"
import { type FieldValues, type UseFormReturn } from "react-hook-form"
import ManagerLayout from "./AdminLayout"
import { type ICreateModalProps } from "../interfaces/ICreateModalProps"
import { type ICardProps } from "../interfaces/ICardProps"

export interface List<T> {
  data: T[] | undefined
  pagination: CustomPaginationProps
}

interface ListViewLayoutProps<
  T,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  U extends List<any>,
  V extends FieldValues
> {
  title: string
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  form: UseFormReturn<any>
  list: List<T>
  setList: Dispatch<SetStateAction<List<T>>>
  apiCall: (page?: number) => Promise<APIResponse<U>>
  preprocessList?: (data: U) => Promise<List<T>>
  createModal?: (props: ICreateModalProps<V>) => ReactNode
  card: (props: ICardProps<T>) => ReactNode
}

// eslint-disable-next-line max-lines-per-function
export default function ListViewLayout<
  T,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  U extends List<any>,
  V extends FieldValues
>({
  title,
  form,
  list,
  setList,
  createModal,
  apiCall,
  preprocessList,
  card
}: ListViewLayoutProps<T, U, V>) {
  const router = useRouter()

  const fetchData = useCallback(async () => {
    if (!router.isReady) return

    const page = router.query.page
      ? parseInt(router.query.page as string)
      : undefined

    const response = await apiCall(page)
    if (!response.data) return

    if (
      response.data.pagination.total !== 0 &&
      response.data.pagination.current > response.data.pagination.total
    ) {
      await router.push(
        `${router.pathname}?page=${response.data.pagination.total}`
      )
    }

    const data = response.data
    if (preprocessList) {
      setList(await preprocessList(data))
    } else {
      setList(data as unknown as List<T>)
    }
  }, [apiCall, preprocessList, router, setList])

  useEffect(() => {
    void fetchData()
  }, [fetchData])

  return (
    <ManagerLayout
      title={title}
      titleButton={createModal?.({ form, fetchData })}
    >
      <br />

      <div className="min-vh-51">
        {!list.data && <Loader />}

        {list.data && list.data.length == 0 ? "Nothing to see here." : ""}

        {list.data?.map((item, index) => {
          return (
            <div key={index}>
              {card({ data: item, refetchData: fetchData })}
            </div>
          )
        })}
      </div>

      <CustomPagination
        current={list.pagination.current}
        total={list.pagination.total}
      />
    </ManagerLayout>
  )
}
