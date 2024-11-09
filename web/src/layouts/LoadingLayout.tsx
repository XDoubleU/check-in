import Loader, { type LoaderProps } from "components/Loader"
import Head from "next/head"

export default function LoadingLayout(props: Readonly<LoaderProps>) {
  return (
    <>
      <Head>
        <title>Loading...</title>
      </Head>
      <Loader {...props} />
    </>
  )
}
