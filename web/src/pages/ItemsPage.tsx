import Header from "./Header";
import {ItemData, ItemRow} from "../library/item";
import {useEffect, useState} from "react";


export interface ItemsPageProps {
    // Items: ItemData[];
}
export function ItemsPage(props: ItemsPageProps) {
    const initialItems:ItemData[] = [];
    const [items, setItems] = useState(initialItems);

    useEffect(() => {
        fetch(`http://localhost:8080/api/items`)
            .then((resp) => resp.json())
            .then((actualData) => {
                setItems(actualData);
            })
            .catch((err) => {
                console.log(err.message);
            });
    }, []);


    const itemRows = items.map(item =>
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