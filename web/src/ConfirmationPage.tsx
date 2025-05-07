import { useNavigate, useParams } from "react-router-dom"
import { API_URL } from "./App"

export const ConfirmationPage = () =>{
    const {token = ''} = useParams()
    const redirect = useNavigate()



    const handleConfirmation = async ()=>{
        const response = await fetch (`${API_URL}/users/activate/${token}`,{
            method: "PUT"
        })

        if (response.ok){
            // redirect to  / page
            redirect("/")

        }else{
            //handle error
        }
    }


    return(
        <div>
            <h2>Confirmation</h2>
            <button onClick={handleConfirmation}>Click to Confirma</button>
        </div>
    )

}