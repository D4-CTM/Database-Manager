export interface TableData {
    Name: string
    ColumnNames: string[]
    Rows: any[][]
}

export interface ColumnsMetadata {
    Name: string
    DataType: string
    DefaultData: string
    Detail: string
    IsNullable: boolean
    IsIdentity: boolean
    OrdPosition: number
}

export interface ConstraintMetadata {
    ConstraintName: string
    ConstraintType: string
    ColumnName: string
    SearchCondition: string
    RefConstraint: string
    Position: number
}

export interface FunctionArgument {
    Name: string
    Position: number
    DataType: string
    InOut: string
    Length: number
    Precision: number
    Scale: number
    HasDefault: boolean
}
