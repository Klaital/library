import Header from "./Header";
import {FormEvent, useEffect, useMemo, useState} from "react";
import {LocationData} from "../library/location";
import Select from 'react-select';
import {LibraryApiClient} from "../library/client";

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
    value: number,
    label: string,
}

interface CodeTypeSelectOptions {
    value: string,
    label: string,
}
export function AddItemPage(props: AddItemPageProps) {
    const initialLocations: LocationSelectOptions[] = [];
    const [locations, setLocations] = useState(initialLocations);
    const [codeType, setCodeType] = useState('isbn')
    const codeTypeOptions = [
        {
            value: "upc",
            label: "UPC",
        },
        {
            value: "isbn",
            label: "ISBN",
        }
    ]

    const api = useMemo(() => new LibraryApiClient(), []);
    function handleAddItemSubmit(event: FormEvent<NewItemFormElement>) {
        event.preventDefault();
        console.log(event.currentTarget.elements.titleInput.value);
    }

    function handleCodeTypeChanges(selected?: CodeTypeSelectOptions | CodeTypeSelectOptions[] | null) {
        if (selected == null) {
            return
        }
        if (Array.isArray(selected)) {
            console.log("Unexpected array for Code Type")
            return
        }
        setCodeType(selected.value);
    }
    function handleCreateLocation(event: FormEvent<CreateLocationFormElement>) {
        event.preventDefault();
        fetch('https://library.klaital.com/api/locations', {
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
        api.FetchLocations().then(
            (resp) => {
                const locationsOptions = resp.map(function(l:LocationData) {
                    return {
                        value: l.ID,
                        label: l.Name,
                    }
                })
                setLocations(locationsOptions)
            })
            .catch((error) => {
                    console.log("Failed to fetch locations")
            })
    }, [api])
    return (
        <>
            <Header />
            <form onSubmit={handleAddItemSubmit}>
                <h2>Add An Item Here</h2>
                <div className="form-item">
                    <label htmlFor="locationSelect">Location</label>
                    <Select id="locationSelect" options={locations} />
                </div>
                <div className="form-item">
                    <label htmlFor="codeInput">Code</label>
                    <input id="codeInput" />
                </div>
                <div className="form-item">
                    <label htmlFor="codeTypeInput">Code Type</label>
                    <Select id="codeTypeInput"
                            options={codeTypeOptions}
                            defaultInputValue={codeType}
                            onChange={handleCodeTypeChanges}
                    />
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
