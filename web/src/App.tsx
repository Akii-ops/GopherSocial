import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'



export const API_URL = import.meta.env.VITE_API_URL || "http://localhost:9999/v1"


function App() {
  const [count, setCount] = useState(0)

  return (
    <>
     <div>APP home screen</div>
    </>
  )
}

export default App
