import { ref } from "vue";
import { TableData } from "./TableData";

export enum TabOptions {
    Table = 0,
    Query = 10,
}

// When TabOption is a Query the payload is a string
// when is a Table the payload is TableData
export interface TabData {
    Title: string,
    Type: TabOptions,
    Payload: TableData | string | null
}

export class TabStore {
    private tabs = ref([] as TabData[])
    public selectedTab = ref(0)

    public readonly currentIdx = () => this.selectedTab.value
    public readonly length = () => this.tabs.value.length
    public readonly list = () => this.tabs.value

    public readonly currentPayload = <T extends String | TableData>() => this.tabs.value[this.currentIdx()].Payload as T
    public readonly currentTabType = () => this.tabs.value[this.currentIdx()].Type

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
