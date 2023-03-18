import Head from "next/head"
import NextNProgress from "nextjs-progressbar"

export default function LoadingLayout() {
  return (
    <>
      <Head>
        <title>Loading...</title>
      </Head>
      <NextNProgress color="red" />
    </>
  )
}
