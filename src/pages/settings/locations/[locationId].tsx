import { prisma } from "@/common/prisma"
import AdminLayout from "@/layouts/AdminLayout"
import LoadingLayout from "@/layouts/LoadingLayout"
import { authOptions } from "@/pages/api/auth/[...nextauth]"
import { Location } from "@prisma/client"
import { GetServerSidePropsContext } from "next"
import { getServerSession, Session } from "next-auth"
import { useSession } from "next-auth/react"

type LocationDetailProps = {
  location: Location
}

export default function LocationDetail({location}: LocationDetailProps) {
  const {data, status} = useSession({
    required: true
  })

  if (status == "loading") {
    return <LoadingLayout/>
  }

  return (
    <AdminLayout title="LocationDetail" isAdmin={data.user.isAdmin}>
      Logged in as {data.user.username}
      Viewing page for {location.id}
    </AdminLayout>
  )  
}

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const session = await getServerSession(context.req, context.res, authOptions) as Session
  const location = await prisma.location.findFirst({
    where: {
      id: context.query.locationId as string,
      userId: session.user.id
    }
  })

  if (location == null) {
    return {
      notFound: true
    }
  }

  return {
    props: {
      location
    }
  }
}