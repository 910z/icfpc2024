import React, {useEffect, useState} from 'react';
import {Table} from "react-bootstrap";
import {Solution} from "../types";
import {formatNum, get} from "./Problems";

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
            <th>Tastes</th>
            <th>Score</th>
            <th>Version</th>
        </tr>
        </thead>
        <tbody>
        {solutions?.map(({id, problemId, score, tags}) => (
            <tr>
                <td>{problemId}</td>
                <td><img src={`/preview/${id}`} alt={`${id}`} width="200" height="200"/></td>
                <td><img src={`/tastes/${problemId}`} alt={`${problemId}`} width="200" height="200" /></td>
                <td>{formatNum(score ?? 0)}</td>
                <td>{tags ?? []}</td>
            </tr>
        ))}
        </tbody>
    </Table>
}
