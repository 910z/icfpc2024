import React, {useState} from 'react';
import {Table} from "react-bootstrap";
import {Problem} from "../types";

export function get(url: string) {
    console.log(`get from ${url}`);
    return fetch(domain + url).then((res) => res.json());
}

export const domain = process.env.API_URL || 'http://localhost:8080'

export const Problems: React.FC = () => {
    const [problems, setProblems] = useState([] as Problem[]);

    // const [bestSolutions, setBestSolutions]= useState({});

    function upd() {
        console.log("upd")
        get(`/api/problems`)
            .then(data => {
                if (data !== problems) {
                    setProblems(data);
                }
            });
    }

    setTimeout(upd, 5000);

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
        {problems?.map(({id, bestSolution}) => (
            <tr>
                <td>{id}</td>
                <td>{
                    bestSolution != null
                        ? <img src={`${domain}/preview/${bestSolution.id}?imgSize=200`} alt={`${id}`}/>
                        : <p>Nope</p>
                }</td>
                <td>{bestSolution?.score ?? 0}</td>
                <td>{bestSolution?.tags ?? []}</td>
            </tr>
        ))}
        </tbody>
    </Table>
}
