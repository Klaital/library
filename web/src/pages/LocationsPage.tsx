import Header from "./Header";
import {useEffect, useState} from "react";
import {LocationData, Location} from "../library/location";

export default function LocationsPage() {
    const initialLocations:LocationData[] = []
    const [locations, setLocations] = useState(initialLocations);
    const [titleFilter, setTitleFilter] = useState("");
    useEffect(() => {
        fetch(`http://localhost:8080/api/locations`)
            .then((resp) => resp.json())
            .then((actualData) => {
                setLocations(actualData);
            })
            .catch((err) => {
                console.log(err.message);
            });
    }, []);

    const locationsSet = locations.map(loc =>
        <Location key={loc.ID}
                  Location={loc}
                  TitleFilter={titleFilter}
        />
    );

    return (
        <>
        <Header />
        <h2>Locations</h2>
        <div>
            <label htmlFor="filter-box">Filter by title: </label>
            <input onChange={e => setTitleFilter(e.target.value)} id="filter-box"/>
        </div>
        {locationsSet}
        </>
    )
}
