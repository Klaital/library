import {ItemData, ItemRow} from "./item";

export interface LocationData {
    Name: string;
    ID: number;
    Notes: string;
    Items: ItemData[];
}

export interface LocationProps {
    Location: LocationData;
    TitleFilter: string
}
export function Location(props:LocationProps) {
    const itemRows = props.Location.Items.map(item =>
        (item.Title.toUpperCase().indexOf(props.TitleFilter.toUpperCase()) > -1)
            ? <ItemRow key={item.ID} ID={item.ID} LocationID={item.LocationID} Code={item.Code} CodeSource={item.CodeSource} Title={item.Title} />
            : null
    ).filter(n=>n);
    return (
        <div className="location-block">
            <h2>{props.Location.Name}</h2>
            <div className="location-items">
                <h3>Items here</h3>
                <table>
                    <thead>
                    <tr>
                        <th>ID</th>
                        <th>Name</th>
                    </tr>
                    </thead>
                    <tbody>
                    {itemRows}
                    </tbody>
                </table>
            </div>
        </div>
    )
}
