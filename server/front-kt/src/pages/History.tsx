import React, {useEffect, useState} from 'react';
import {HistoryResponse} from "../types";
import {Divider, Paper, ScrollArea, SimpleGrid, Table} from "@mantine/core";
import {Badge} from "@mantine/core/lib";
import './History.css';

export function get(url: string) {
    return fetch(url).then((res) => res.json());
}

export function formatNum(num: string | number): string {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
}

export function BigText(text: string) {
    if (text.length > 32) {
        return <p>{text.substring(0, 32)}...</p>
    } else {
        return <p>{text}</p>
    }
}

export const HistoryPage: React.FC = () => {
    const [history, setHistory] = useState({} as HistoryResponse);
    const [select, setSelect] = useState("");

    useEffect(() => {
        fetch(`/api/history`)
            .then((res) => res.json())
            .then(data => {
                setHistory(data);
            });
    }, []);

    const hist = history.history ?? [];
    const content = history.content ?? {};

    const preview = hist.find((obj) => obj.uuid == select);

    return <SimpleGrid cols={{ sm: 1, md: 2 }}>
        <ScrollArea scrollbars="y">
        <Table striped highlightOnHover withTableBorder>
            <Table.Thead>
                <Table.Tr>
                    <Table.Th>createdAt</Table.Th>
                    <Table.Th>request</Table.Th>
                    <Table.Th>response</Table.Th>
                </Table.Tr>
            </Table.Thead>
            <Table.Tbody>{
                hist.map(value => (
                        <Table.Tr onClick={() => setSelect(value.uuid)}>
                            <Table.Td>{value.createdAt}</Table.Td>
                            <Table.Td>{BigText(content[value.request].content)}</Table.Td>
                            <Table.Td>{BigText(content[value.response].content)}</Table.Td>
                        </Table.Tr>
                    )
                )
            }</Table.Tbody>
        </Table>
        </ScrollArea>
        <ScrollArea scrollbars="y">
            {
                preview && <div>
                    {preview.createdAt}
                    <Divider my="md" />
                <Paper className="break" shadow="xs" withBorder p="xl">
                    {content[preview.request].content}
                </Paper>
                    <Divider my="md" />
                    <Paper className="break" shadow="xs" withBorder p="xl">
                        {content[preview.response].content}
                    </Paper>
                </div>
            }
        </ScrollArea>
    </SimpleGrid>

}
