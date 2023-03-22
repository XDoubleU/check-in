import CustomPagination, {
  type CustomPaginationProps
} from "@/components/CustomPagination"
import SchoolCard from "@/components/cards/SchoolCard"
import AdminLayout from "@/layouts/AdminLayout"
import { type CreateSchoolDto, type School } from "types-custom"
import { createSchool, getAllSchools } from "my-api-wrapper"
import { useRouter } from "next/router"
import { useCallback, useEffect, useState } from "react"
import { Form } from "react-bootstrap"
import { useForm } from "react-hook-form"
import CreateModal from "@/components/modals/CreateModal"

interface SchoolList {
  schools: School[]
  pagination: CustomPaginationProps
}

// eslint-disable-next-line max-lines-per-function
export default function SchoolList() {
  const router = useRouter()

  const [schoolList, setSchoolList] = useState<SchoolList>({
    schools: [],
    pagination: {
      current: 0,
      total: 0
    }
  })

  const form = useForm<CreateSchoolDto>()

  const fetchData = useCallback(async () => {
    if (!router.isReady) return

    const page = router.query.page
      ? parseInt(router.query.page as string)
      : undefined

    const response = await getAllSchools(page)

    setSchoolList({
      schools: response.data?.schools ?? Array<School>(),
      pagination: {
        current: response.data?.page ?? 1,
        total: response.data?.totalPages ?? 1
      }
    })
  }, [router])

  useEffect(() => {
    void fetchData()
  }, [fetchData])

  const handleCreate = (data: CreateSchoolDto) => {
    return createSchool(data)
  }

  return (
    <AdminLayout title="Schools">
      <CreateModal<CreateSchoolDto, School>
        form={form}
        handler={handleCreate}
        refetchData={fetchData}
        typeName="school"
      >
        <Form.Group className="mb-3">
          <Form.Label>Name</Form.Label>
          <Form.Control
            type="text"
            placeholder="Name"
            required
            {...form.register("name")}
          ></Form.Control>
        </Form.Group>
      </CreateModal>

      <br />

      <div className="min-vh-51">
        {schoolList.schools.length == 0 ? "Nothing to see here." : ""}

        {schoolList.schools.map((school) => {
          return (
            <SchoolCard
              id={school.id}
              key={school.id}
              name={school.name}
              refetchData={fetchData}
            />
          )
        })}
      </div>

      <CustomPagination
        current={schoolList.pagination.current}
        total={schoolList.pagination.total}
      />
    </AdminLayout>
  )
}
