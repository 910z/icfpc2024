import React, {useEffect, useState} from 'react';
import {HistoryResponse, Tokens} from "../types";
import {Divider, Paper, ScrollArea, SimpleGrid, Table, Tooltip} from "@mantine/core";
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

export function Token(tokens: Tokens, text: string) {
    if (tokens[text]) {
        const t = tokens[text];
        console.log(t);
        return <Tooltip multiline label={t} withArrow>
            <span>{text}&nbsp;</span>
        </Tooltip>
    } else {
        return <span>{text}&nbsp;</span>
    }
}

export function previewTokens(tokens: Tokens, text: string) {
    const list = text.split(" ");
    if (list.length == 1) {
        return tokens[text]
    } else {
        return <p className="break">
            {text.split(" ")
                .map((text) => Token(tokens, text))}
        </p>
    }
}

export function ShowTokens(tokens: Tokens, text: string) {
    const list = text.split(" ");
    if (list.length == 1) {
        return tokens[text]?.split("\n")?.map(line => <p className="break">{line}</p>)
    } else {
        return <p className="break">
            {text.split(" ")
                .map((text) => Token(tokens, text))}
        </p>
    }
}

export const HistoryPage: React.FC = () => {
    const [history, setHistory] = useState({} as HistoryResponse);
    const [select, setSelect] = useState("");
    // const [tokens, setTokens] = useState({} as Tokens);

    useEffect(() => {
        fetch(`/api/history`)
            .then((res) => res.json())
            .then(data => {
                setHistory(data);
            });
    }, []);

    const hist = history.history ?? [];
    const content = history.content ?? {};
    const tokens = history.tokens ?? {};

    const preview = hist.find((obj) => obj.uuid == select);
    // useEffect(() => {
    //     fetch(`/api/tokens?uuid=` + preview?.uuid)
    //         .then((res) => res.json())
    //         .then(data => {
    //             setTokens(data);
    //         });
    // }, [preview]);

    // <ScrollArea scrollbars="y">
    return <div>
        <SimpleGrid cols={{sm: 1, md: 2}}>
            <ScrollArea scrollbars="y" style={{height:"calc(100vh - 40px)"}}>
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
                                    <Table.Td>{value.createdAt.replace("T","\n").replace("Z","")}</Table.Td>
                                    <Table.Td className="over" style={{maxWidth: "200px"}}>
                                        {previewTokens(tokens, content[value.request].content)}
                                    </Table.Td>
                                    <Table.Td className="over">
                                        {previewTokens(tokens, content[value.response].content)}
                                    </Table.Td>
                                </Table.Tr>
                            )
                        )
                    }</Table.Tbody>
                </Table>
            </ScrollArea>
            <ScrollArea style={{height:"calc(100vh - 40px)"}}>
                {
                    preview && <div>
                        {preview.createdAt}
                        <Divider my="md"/>
                        {ShowTokens(tokens, content[preview.request].content)}
                        <Divider my="md"/>
                        {ShowTokens(tokens, content[preview.response].content)}
                    </div>
                }
            </ScrollArea>
        </SimpleGrid>
    </div>
}
