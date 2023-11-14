import Header from "./Header";
import Select from "react-select";
import {FormEvent, useRef, useEffect, useMemo, useState} from "react";
import {LocationData} from "../library/location";
import {LibraryApiClient} from "../library/client";

interface LocationSelectOptions {
    value: number,
    label: string,
}

interface CodeTypeSelectOptions {
    value: string,
    label: string,
}

interface NewItemFormElements extends HTMLFormControlsCollection {
    titleInput: HTMLInputElement,
    titleTranslatedInput: HTMLInputElement,
}
interface NewItemFormElement extends HTMLFormElement {
    readonly elements: NewItemFormElements
}

interface NewCodeFormElements extends HTMLFormControlsCollection {
    codeInput: HTMLInputElement,
}
interface NewCodeFormElement extends HTMLFormElement {
    readonly elements: NewCodeFormElements
}

interface CreateLocationFormElements extends HTMLFormControlsCollection {
    locationNameInput: HTMLInputElement,
    locationNotesInput: HTMLInputElement,
}
interface CreateLocationFormElement extends HTMLFormElement {
    readonly elements: CreateLocationFormElements
}

export function AddVolumePage() {
    const initailLocation: LocationSelectOptions[] = []
    const [ locations, setLocations ] = useState(initailLocation)
    const [ codeType, setCodeType ] = useState("upc")
    const [ code, setCode ] = useState("")
    const [ title, setTitle ] = useState("")
    const [ titleTranslated, setTitleTranslated ] = useState("")
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

    const api = useMemo(() => new LibraryApiClient(), []);
    const titleInputRef = useRef<HTMLInputElement | null>(null);

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

    function handleAddVolumeSubmit(event: FormEvent<NewItemFormElement>) {
        event.preventDefault();
        console.log("form submitted: "+ event.currentTarget.elements.titleInput.value)
    }

    function handleCodeInput(event: FormEvent<NewCodeFormElement>) {
        event.preventDefault();
        console.log("code submitted: " + codeType + " / " + event.currentTarget.elements.codeInput.value);
        api.LookupCode(event.currentTarget.elements.codeInput.value, codeType)
            .then(
                (res) => {
                    setTitle(res.Title);
                    titleInputRef.current?.focus();
                })
            .catch((err) => {
                console.log("Failed code lookup: " + err);
            })
    }

    function handleCreateLocation(event: FormEvent<CreateLocationFormElement>) {
        event.preventDefault();
        api.CreateLocation(
            event.currentTarget.elements.locationNameInput.value,
            event.currentTarget.elements.locationNotesInput.value)
    }

    return <>
        <Header />
        <form onSubmit={handleCodeInput}>
            <div className="form-item">
                <label htmlFor="codeTypeInput">Code Type</label>
                <Select id="codeTypeInput"
                        options={codeTypeOptions}
                        onChange={handleCodeTypeChanges}
                />
            </div>
            <div className="form-item">
                <label htmlFor="codeInput">Code</label>
                <input id="codeInput"
                       value={code}
                       onChange={e => setCode(e.target.value)}
                />
            </div>
            <input type="submit" value="Lookup" />
        </form>
        <form onSubmit={handleAddVolumeSubmit}>
            <div className="form-item">
                <label htmlFor="locationSelect">Location</label>
                <Select id="locationSelect" options={locations} />
            </div>
            <div className="form-item">
                <label htmlFor="titleInput">Title</label>
                <input id="titleInput"
                       value={title}
                       ref={titleInputRef}
                       onChange={e => setTitle(e.target.value)}
                />
            </div>
            <div className="form-item">
                <label htmlFor="titleTranslatedInput">Title (Translated)</label>
                <input id="titleTranslatedInput"
                       value={titleTranslated}
                       onChange={e => setTitleTranslated(e.target.value)}
                />
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
    </>
}