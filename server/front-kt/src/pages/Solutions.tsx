import React, {useState} from 'react';
import {Table} from "react-bootstrap";
import {Solution} from "../types";
import {domain, get} from "./Problems";

export const Solutions: React.FC = () => {
    const [solutions, setSolutions] = useState([] as Solution[]);

    function upd() {
        get(`/api/solutions`)
            .then(data => {
                if (data !== solutions) {
                    setSolutions(data);
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
        {solutions?.map(({id, problemId, score, tags}) => (
            <tr>
                <td>{problemId}</td>
                <td><img src={`${domain}/preview/${id}?imgSize=200`} alt={`${id}`}/></td>
                <td>{score ?? 0}</td>
                <td>{tags ?? []}</td>
            </tr>
        ))}
        </tbody>
    </Table>
}
