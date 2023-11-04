
export interface ItemData {
    ID: number;
    LocationID: number;
    Code: string;
    CodeSource: string;
    Title: string;
}
export function ItemRow(props:ItemData) {
    return (
        <tr>
            <td>
                {props.ID}
            </td>
            <td>
                {props.Title}
            </td>
        </tr>
    )
}