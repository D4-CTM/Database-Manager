import { TableData } from "./TableData";

export enum TabOptions {
    Table = 0,
    Query = 10,
}

export interface TabData {
    Title: string,
    Type: TabOptions,
    Payload: TableData | null
}
