import Head from "next/head"

export default function Error() {
  const style: React.CSSProperties = {
    position: "fixed",
    top: "50%",
    left: "50%",
    transform: "translate(-50%, -50%)",
    textAlign: "center"
  }

  return (
    <>
      <Head>
        <title>Something went wrong</title>
      </Head>
      <div style={style}>
        <h1>Oops, something went wrong</h1>
      </div>
    </>
  )
}
