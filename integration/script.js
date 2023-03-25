const eventSource = new EventSource("http://localhost:8000/sse/")

eventSource.onopen = async () => {
  const response = await fetch("http://localhost:8000/locations/sse/")
  const data = await response.json()
  
  data.forEach((location) => {
    fill(location)
  })
}

eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data)
  fill(data)
}

window.onbeforeunload = () => {
  eventSource.close()
}

function fill(location){
  let element = document.getElementById(location.normalizedName);
  
  if(element){
      element.innerHTML=`${location.available}/${location.capacity}`;
  }
}