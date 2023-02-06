import Head from "next/head"
import { Oval } from "react-loader-spinner"


export default function Spinner(){
  return (
    <>
      <Head>
        <title>Loading...</title>
      </Head>
      <Oval
              height={80}
              width={80}
              color="red"
              wrapperStyle={{"justifyContent": "center", "alignItems": "center", "height": "100vh"}}
              wrapperClass=""
              visible={true}
              ariaLabel='oval-loading'
              secondaryColor="white"
              strokeWidth={2}
              strokeWidthSecondary={2} 
            />
    </>
  )
}