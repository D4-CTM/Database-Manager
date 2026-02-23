export interface DbCredential {
    Database: string
    Server: string
    Port: number
    User: string
    Password: string
    ShowAll: boolean

    conName: string
    isNew: boolean
}

export interface PutCredentialResult {
    OldName: string
    NewName: string
}
