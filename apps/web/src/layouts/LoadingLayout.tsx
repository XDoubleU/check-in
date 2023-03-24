import Loader from "@/components/Loader"
import Head from "next/head"

export default function LoadingLayout() {
  return (
    <>
      <Head>
        <title>Loading...</title>
      </Head>
      <Loader />
    </>
  )
}
