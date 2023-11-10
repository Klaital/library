import Header from "./Header";
import {FormEvent, useEffect, useState} from "react";
import {LocationData} from "../library/location";
import Select from 'react-select';

interface NewItemFormElements extends HTMLFormControlsCollection {
    codeInput: HTMLInputElement,
    codeTypeInput: HTMLInputElement,
    titleInput: HTMLInputElement,
    titleTranslatedInput: HTMLInputElement,
}
interface NewItemFormElement extends HTMLFormElement {
    readonly elements: NewItemFormElements
}

interface CreateLocationFormElements extends HTMLFormControlsCollection {
    locationNameInput: HTMLInputElement,
    locationNotesInput: HTMLInputElement,
}
interface CreateLocationFormElement extends HTMLFormElement {
    readonly elements: CreateLocationFormElements
}
export interface AddItemPageProps {

}

interface LocationSelectOptions {
    value: string,
    label: string,
}
export function AddItemPage(props: AddItemPageProps) {
    const initialLocations: LocationSelectOptions[] = [];
    const [locations, setLocations] = useState(initialLocations);
    function handleAddItemSubmit(event: FormEvent<NewItemFormElement>) {
        event.preventDefault();
        console.log(event.currentTarget.elements.titleInput.value);
    }

    function handleCreateLocation(event: FormEvent<CreateLocationFormElement>) {
        event.preventDefault();
        fetch('http://localhost:8080/api/locations', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            }, body: JSON.stringify({
                Name: event.currentTarget.elements.locationNameInput.value,
                Notes: event.currentTarget.elements.locationNotesInput.value,
            })
        })
            .then(res => res.json())
            .then(
                (result) => {
                    locations.push({
                        value: result.ID,
                        label: result.Name,
                    })
                    setLocations(locations)
                },
                (error) => {
                    console.log("Failed to create location")
                    console.log(error)
                }
            )
    }

    useEffect(() => {
        fetch(`http://localhost:8080/api/locations`)
            .then((resp) => resp.json())
            .then(
                (result) => {
                    const optionset = result.map(function(l:LocationData) {
                        return {
                            value: l.ID,
                            label: l.Name,
                        }
                    })
                    setLocations(optionset)
                },
                (error) => {
                    console.log("Failed to fetch locations")
                    console.log(error)
                }
            )
    }, [])
    return (
        <>
            <Header />
            <form onSubmit={handleAddItemSubmit}>
                <h2>Add An Item Here</h2>
                <div className="form-item">
                    <label htmlFor="locationSelect">Code</label>
                    <Select id="locationSelect" options={locations} />
                </div>
                <div className="form-item">
                    <label htmlFor="codeInput">Code</label>
                    <input id="codeInput" />
                </div>
                <div className="form-item">
                    <label htmlFor="codeTypeInput">Code Type</label>
                    <input id="codeTypeInput" />
                </div>
                <div className="form-item">
                    <label htmlFor="titleInput">Title</label>
                    <input id="titleInput" />
                </div>
                <div className="form-item">
                    <label htmlFor="titleTranslatedInput">Title (Translated)</label>
                    <input id="titleTranslatedInput" />
                </div>
                <input type="submit" value="Submit" />
            </form>

            <form onSubmit={handleCreateLocation}>
                <h2>Create a new Location here</h2>
                <div className="form-item">
                    <label htmlFor="locationNameInput">New Location Name: </label>
                    <input id="locationNameInput" />
                </div>
                <div className="form-item">
                    <label htmlFor="locationNotesInput">New Location Notes: </label>
                    <input id="locationNotesInput" />
                </div>
                <input type="submit" value="Submit" />
            </form>
            <h2>History</h2>

        </>
    )
}
