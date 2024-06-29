import React, {useEffect, useState} from 'react';
import {HistoryResponse} from "../types";
import {Table} from "@mantine/core";

export function get(url: string) {
    return fetch(url).then((res) => res.json());
}

export function formatNum(num: string | number): string {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
}

export const HistoryPage: React.FC = () => {
    const [history, setHistory] = useState({} as HistoryResponse);

    useEffect(() => {
        fetch(`/api/history`)
            .then((res) => res.json())
            .then(data => {
                setHistory(data);
            });
    }, []);

    const hist = history.history ?? [];
    const content = history.content ?? {};

    return <Table>
        <Table.Thead>
            <Table.Tr>
                <Table.Th>createdAt</Table.Th>
                <Table.Th>request</Table.Th>
                <Table.Th>response</Table.Th>
            </Table.Tr>
        </Table.Thead>
        <Table.Tbody>{
            hist.map(value => (
                    <Table.Tr>
                        <Table.Td>{value.createdAt}</Table.Td>
                        <Table.Td>{content[value.request].content}</Table.Td>
                        <Table.Td>{content[value.response].content}</Table.Td>
                    </Table.Tr>
                )
            )
        }</Table.Tbody>
    </Table>
}
