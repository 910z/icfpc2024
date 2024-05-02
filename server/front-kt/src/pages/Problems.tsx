import React, {useEffect, useState} from 'react';
import {Table} from "react-bootstrap";
import {Problem} from "../types";

export function get(url: string) {
    return fetch(url).then((res) => res.json());
}

export function formatNum(num: string | number): string {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
}

export const Problems: React.FC = () => {
    const [problems, setProblems] = useState([] as Problem[]);

    useEffect(() => {
        fetch(`/api/problems`)
            .then((res) => res.json())
            .then(data => {
                setProblems(data);
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
        {problems?.map(({id, bestSolution}) => (
            <tr>
                <td>{id}</td>
                <td>{
                    bestSolution != null
                        ? <img src={`/preview/${bestSolution.id}`} alt={`${id}`} width="200" height="200"/>
                        : <p>Nope</p>
                }</td>
                <td><img src={`/tastes/${id}`} alt={`${id}`} width="200" height="200" /></td>
                <td>{formatNum(bestSolution?.score ?? 0)}</td>
                <td>{bestSolution?.tags ?? []}</td>
            </tr>
        ))}
        </tbody>
    </Table>
}
