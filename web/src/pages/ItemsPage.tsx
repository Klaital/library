import Header from "./Header";
import {ItemData, ItemRow} from "../library/item";


export interface ItemsPageProps {
    Items: ItemData[];
}
export function ItemsPage(props: ItemsPageProps) {
    const itemRows = props.Items.map(item =>
        <ItemRow
            key={item.ID}
            ID={item.ID}
            LocationID={item.LocationID}
            Code={item.Code} CodeSource={item.CodeSource}
            Title={item.Title}
        />
    )
    return (
        <>
        <Header />
        <h2>Items</h2>
        <table>
            <thead>
            <tr>
                <th>Item ID</th>
                <th>Title</th>
            </tr>
            </thead>
            <tbody>
            {itemRows}
            </tbody>
        </table>
        </>
    )
}