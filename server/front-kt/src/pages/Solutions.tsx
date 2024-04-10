import React, {useState} from 'react';
import {Table} from "react-bootstrap";
import {Solution} from "../types";

function get(url: string) {
    console.log(`get from ${url}`);
    return fetch(url).then((res) => res.json());
}

export const Solutions: React.FC = () => {
    const [solutions, setSolutions] = useState([] as Solution[]);

    // const [bestSolutions, setBestSolutions]= useState({});

    function upd() {
        console.log("upd")
        get("http://localhost:8080/api/solutions")
            .then(data => {
                if (data != solutions) {
                    setSolutions(data);
                }
            });
        // get("http://localhost:8080/api/solution/best")
        //     .then(data => {
        //         if (data != bestSolutions) {
        //             setBestSolutions(data)
        //         }
        //     });
    }

    //
    // upd();

    setTimeout(upd, 5000);

    // setInterval(upd, 5000);

    // useEffect(() => {
    //     document.
    //     if (darkMode) {
    //         document.documentElement.removeAttribute("data-bs-theme");
    //     } else {
    //         document.documentElement.setAttribute("data-bs-theme", "dark");
    //     }
    // }, [problems]);

    // const root = ReactDOM.createRoot(
    //     document.getElementById('root')
    // );
    //
    // function tick() {
    //     const element = (
    //         <div>
    //             <h1>Hello, world!</h1>
    //             <h2>It is {new Date().toLocaleTimeString()}.</h2>
    //         </div>
    //     );
    //     root.render(element);
    // }
    //
    // setInterval(tick, 1000);

    return <Table striped bordered hover>
        <thead>
        <tr>
            <th>#</th>
            <th>Preview</th>
            <th>Instrs</th>
            <th>Musicns</th>
            <th>Attends</th>
            <th>Tastes</th>
            <th>Pillars</th>
            <th>Stage Size</th>
            <th>Score</th>
        </tr>
        </thead>
        <tbody>
        {solutions?.map(({id, problemId, score, tags}) => (
            <tr>
                <td>{problemId}</td>
                <td><img src={`http://localhost:8080/preview/${id}?imgSize=200`}/></td>
                <td>{score ?? 0}</td>
                <td>{tags ?? []}</td>
            </tr>
        ))}
        </tbody>
    </Table>
}