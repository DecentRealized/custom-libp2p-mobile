import type {NextPage} from 'next'
import {HelloWorldSample} from "./caller";
import {useEffect, useState} from "react";

const Home: NextPage = () => {
    const [state, setState] = useState({
        error: "",
        response: "",
    })

    return (
        <div style={{display: "flex", flexDirection: "column", height: "100vh",
            justifyContent: "center", alignItems: "center"}}>
            <button onClick={async () => {
                const result = await HelloWorldSample()
                setState({error: result.error(), response: result.output()})
            }}>
                Hello World
            </button>
            {
                state.error && <span>Error From WASM: <b>{state.error}</b></span>
            }
            {
                state.response && <span>Response From WASM: <b>{state.response}</b></span>
            }
        </div>
    )
}

export default Home
