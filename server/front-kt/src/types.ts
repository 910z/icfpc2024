export type Problem = {
    id: number
    bestSolution: Solution | null
}



export type Solution = {
    id: number
    problemId: number
    score: number
    tags: string[]
}

export type HistoryResponse = {
    history: History[]
    content: { [_: string]: Content }
}

export type History = {
    uuid: string
    createdAt: string
    request: string
    response: string
}

export type Content = {
    id: string
    content: string
}
