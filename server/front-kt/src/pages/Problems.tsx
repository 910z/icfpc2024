import React, {useState} from 'react';
import {Table} from "react-bootstrap";

function get(url: string) {
    console.log(`get from ${url}`);
    return fetch(url).then((res) => res.json());
}

export const Problems: React.FC = () => {
    const [problems, setProblems] = useState([]);

    function upd() {
        console.log("upd")
        get("http://localhost:8080/api/problems")
            .then(data => setProblems(data));
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
        {problems?.map(({id, contentId}) => (
            <tr>
                <td>{id}</td>
                <td>{contentId}</td>
            </tr>
        ))}
        </tbody>
    </Table>
}