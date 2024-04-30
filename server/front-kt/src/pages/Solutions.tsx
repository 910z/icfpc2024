import React, {useEffect, useState} from 'react';
import {Table} from "react-bootstrap";
import {Solution} from "../types";
import {get} from "./Problems";

export const Solutions: React.FC = () => {
    const [solutions, setSolutions] = useState([] as Solution[]);

    useEffect(() => {
        fetch(`/api/solutions`)
            .then((res) => res.json())
            .then(data => {
                    setSolutions(data);
            });
    }, []);

    return <Table striped bordered hover>
        <thead>
        <tr>
            <th>#</th>
            <th>Preview</th>
            {/*<th>Instrs</th>*/}
            {/*<th>Musicns</th>*/}
            {/*<th>Attends</th>*/}
            {/*<th>Tastes</th>*/}
            {/*<th>Pillars</th>*/}
            {/*<th>Stage Size</th>*/}
            <th>Score</th>
            <th>Version</th>
        </tr>
        </thead>
        <tbody>
        {solutions?.map(({id, problemId, score, tags}) => (
            <tr>
                <td>{problemId}</td>
                <td><img src={`/preview/${id}?imgSize=200`} alt={`${id}`} width="200" height="200"/></td>
                <td>{score ?? 0}</td>
                <td>{tags ?? []}</td>
            </tr>
        ))}
        </tbody>
    </Table>
}
