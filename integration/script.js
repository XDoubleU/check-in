const webSocket = new WebSocket("ws://localhost:8000")

eventSource.onopen = async () => {
  const response = await fetch("http://localhost:8000/locations/ws/")
  const data = await response.json()
  
  data.forEach((location) => {
    fill(location)
  })
}

webSocket.onmessage = (event) => {
  const data = JSON.parse(event.data)
  fill(data)
}

window.onbeforeunload = () => {
  webSocket.close()
}

let hasLiveDotStyle = false
function fill(location){
  let element = document.getElementById(location.normalizedName);
  
  if(element){
      const capacity = document.createElement("p")
      const yesterdayFullAt = document.createElement("p")
      
      capacity.innerHTML = `${getLiveDot()} Capaciteit: ${location.available}/${location.capacity}`

      let output = "Gisteren niet volzet"
      if (location.yesterdayFullAt) {
        const time = new Date(location.yesterdayFullAt).toLocaleTimeString([], {
          timeStyle: "short",
          hourCycle: "h23"
        })

        output = `Gisteren vol om: ${time}`
      }

      yesterdayFullAt.innerHTML = output

      element.innerHTML = ""
      element.appendChild(capacity)
      element.appendChild(yesterdayFullAt)
  }
}

function getLiveDot() {
  createLiveDotStyle()

  return `
  <svg height="10" width="10" class="blinking" viewBox="40 45 20 10">
    <circle cx="50" cy="50" r="10" fill="red" /> 
  </svg>
  `
}

function createLiveDotStyle() {
  if (hasLiveDotStyle) {
    return
  }

  hasLiveDotStyle = true

  const styleNode = document.createElement("style")

  const style = `
  .blinking {
    -webkit-animation: 1s blink ease infinite;
    -moz-animation: 1s blink ease infinite;
    -ms-animation: 1s blink ease infinite;
    -o-animation: 1s blink ease infinite;
    animation: 1s blink ease infinite;
    
  }

  @keyframes "blink" {
    from, to {
      opacity: 0;
    }
    50% {
      opacity: 1;
    }
  }

  @-moz-keyframes blink {
    from, to {
      opacity: 0;
    }
    50% {
      opacity: 1;
    }
  }

  @-webkit-keyframes "blink" {
    from, to {
      opacity: 0;
    }
    50% {
      opacity: 1;
    }
  }

  @-ms-keyframes "blink" {
    from, to {
      opacity: 0;
    }
    50% {
      opacity: 1;
    }
  }

  @-o-keyframes "blink" {
    from, to {
      opacity: 0;
    }
    50% {
      opacity: 1;
    }
  }
  `

  if(!!(window.attachEvent && !window.opera)) {
    styleNode.styleSheet.cssText = style
  } else {
    var styleText = document.createTextNode(style)
    styleNode.appendChild(styleText)
  }
  document.getElementsByTagName('head')[0].appendChild(styleNode)
}