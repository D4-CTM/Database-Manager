import { Tab } from "bootstrap";
import { ref } from "vue";

export enum TabOptions {
	Tables = "TABLE",
	Views = "VIEW",
	Index = "INDEX",
	Sequence = "SEQUENCE",
	Trigger = "TRIGGER",
	Functions = "FUNCTION",
	Procedure = "PROCEDURE",
	Package = "PACKAGE",
    Query = 'Query',
}

// When TabOption is a Query the payload is a string
// when is a Table the payload is TableData
export interface TabData {
    Title: string,
    Type: TabOptions,
    conName: string
    Payload: QueryPayload | null 
}

export interface QueryPayload {
    owner: string
    table: string
    objectType: string
}

export class TabStore {
    private tabs = ref([] as TabData[])
    public selectedTab = ref(0)

    public readonly currentIdx = () => this.selectedTab.value
    public readonly length = () => this.tabs.value.length
    public readonly list = () => this.tabs.value

    public readonly currentPayload = <T extends QueryPayload>() => this.tabs.value[this.currentIdx()].Payload as T
    public readonly currentConn = () => this.tabs.value[this.currentIdx()].conName
    public readonly currentTabType = () => this.tabs.value[this.currentIdx()].Type
    public readonly tabOptions = () => {
        switch (this.currentTabType()) {
            case TabOptions.Tables: return ['Data', 'Columns', 'Constraints','DDL']
            case TabOptions.Views: return ['Data', 'Columns','DDL']
            case TabOptions.Functions:
            case TabOptions.Procedure:
                return [ 'Arguments', 'DDL' ]
            case TabOptions.Package:
                return [ 'Body', 'DDL' ]
            case TabOptions.Index:
            case TabOptions.Sequence:
            case TabOptions.Trigger: return [ 'DDL' ]
            case TabOptions.Query: return [ '' ]
        }
    }

    add(tab: TabData) {
        let count = 0
        this.tabs.value.forEach(x => {
            if (x.Title.startsWith(tab.Title)) count++
        })

        if (count > 0)
            tab.Title += ` #${count}`

        this.tabs.value.push(tab)
    }

    removeAt(title: string) {
        const idx = this.tabs.value.findIndex(x => x.Title === title)
        if (idx === -1) return false;

        console.log(`idx: ${idx}, currentIdx: ${this.currentIdx()}, length: ${this.length()}`)
        if (idx === this.currentIdx()) {
            this.selectedTab.value = Math.max(0, idx - 1)
        } else if (this.currentIdx() >= this.length() - 1) {
            this.selectedTab.value = Math.max(0, this.length() - 2)
        }

        this.tabs.value.splice(idx, 1)
        return true
    }
}
