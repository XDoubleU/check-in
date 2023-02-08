import { prisma } from "@/common/prisma"
import CustomButton from "@/components/CustomButton"
import CustomPagination from "@/components/CustomPagination"
import SchoolCard from "@/components/SchoolCard"
import AdminLayout from "@/layouts/AdminLayout"
import LoadingLayout from "@/layouts/LoadingLayout"
import { School } from "@prisma/client"
import { GetServerSidePropsContext } from "next"
import { useSession } from "next-auth/react"
import { useRouter } from "next/router"
import { FormEvent, useState } from "react"
import { Col, Form, Modal } from "react-bootstrap"

type SchoolListProps = {
  schools: School[],
  totalSchools: number,
  pageSize: number,
  currentPage: number
}

export default function SchoolList({schools, totalSchools, pageSize, currentPage}: SchoolListProps) {
  const router = useRouter()

  const [addInfo, setAddInfo] = useState({ name: "" })
  const [showAdd, setShowAdd] = useState(false)
  const handleCloseAdd = () => setShowAdd(false)
  const handleShowAdd = () => setShowAdd(true)

  const {data, status} = useSession({
    required: true
  })

  const handleAdd = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const data = {
      ...addInfo
    }

    const response = await fetch("/api/schools", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data)
    })

    if (response.status < 300) {
      router.replace(router.asPath)
      addInfo.name = ""
      handleCloseAdd()
    }
  }

  if (status == "loading") {
    return <LoadingLayout/>
  }

  return (
    <AdminLayout title="Schools" isAdmin={data.user.isAdmin}>
      <h1>Schools</h1>
      <br/>
      
      <Modal show={showAdd} onHide={handleCloseAdd}>
        <Modal.Body>
        <Modal.Title>Add school</Modal.Title>
        <br/>
          <Form onSubmit={handleAdd}>
            <Form.Group className="mb-3">
              <Form.Label>Name</Form.Label>
              <Form.Control type="text" placeholder="Name" value={addInfo.name} onChange={({ target}) => setAddInfo({ ...addInfo, name: target.value })}></Form.Control>
            </Form.Group>
            <CustomButton type="button" style={{"float": "left"}}>Cancel</CustomButton>
            <CustomButton type="submit" style={{"float": "right"}}>Add</CustomButton>
          </Form>
        </Modal.Body>
      </Modal>

      <Col size={2}>
        <CustomButton onClick={handleShowAdd}>
          Add
        </CustomButton>
      </Col>

      <br/>

      <div className="min-vh-51">
        {
          (schools === undefined || schools.length == 0) ? "Nothing to see here." : ""
        }

        {
          schools.map((school) => {
            return <SchoolCard id={school.id} key={school.id} title={school.name} />
          })
        }
      </div>

      <CustomPagination current={currentPage} total={totalSchools} pageSize={pageSize} />
      
    </AdminLayout>
  )  
}

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const pageSize = 4
  let currentPage = parseInt(context.query.page as string ?? "1")

  const totalSchools = await prisma.school.count({
    where: {
      id: {
        not: 1
      }
    }
  })

  const totalPages = Math.ceil(totalSchools/pageSize)
  if ( currentPage > totalPages ) {
    currentPage--
    context.query.page = currentPage.toString()
  }

  console.log(currentPage)
  const response = await fetch(`${process.env.NEXTAUTH_URL}/api/schools?page=${currentPage}&pageSize=${pageSize}`)
  const schools = await response.json()

  return {
    props: {
      schools,
      totalSchools,
      pageSize,
      currentPage
    }
  }
}