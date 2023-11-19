import Header from "./Header";
import Select from "react-select";
import {FormEvent, useEffect, useMemo, useRef, useState} from "react";
import {LibraryApiClient} from "../library/client";
import {CodeTypeSelectOptions, LocationSelectOptions} from "./FormComponentHelpers";
import {LocationData} from "../library/location";

interface SelectCodeFormElements extends HTMLFormControlsCollection {
    codeInput: HTMLInputElement,
}
interface SelectCodeFormElement extends HTMLFormElement {
    readonly elements: SelectCodeFormElements
}

function SelectedItemsList(props: {Codes:{Type: string, Code: string}[]}) {
    const listItems = props.Codes.map(x => <li>{x.Type} / {x.Code}</li>)
    return <ul className="SelectedItemsList">{listItems}</ul>
}
export function MoveVolumePage() {
    const initialLocation: LocationSelectOptions[] = []
    const [ locations, setLocations ] = useState(initialLocation)
    const [ codeType, setCodeType ] = useState("upc")
    // const [ code, setCode ] = useState("")
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
    const initialSelectedCodes: { Type: string, Code: string }[] = []
    const [ selectedCodes, setSelectedCodes ] = useState(initialSelectedCodes)

    const api = useMemo(() => new LibraryApiClient(), []);
    const codeInputRef = useRef<HTMLInputElement | null>(null);

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

    function handleSelectCode(event: FormEvent<SelectCodeFormElement>) {
        event.preventDefault()
        let codeDupe: boolean = false

        for (const [k, v] of Object.entries(selectedCodes)) {
            if (v.Code === event.currentTarget.elements.codeInput.value) {
                codeDupe = true
                break
            }
        }

        if (!codeDupe) {
            console.log(`Adding code: ${event.currentTarget.elements.codeInput.value}`)
            selectedCodes.push({Type: codeType, Code: event.currentTarget.elements.codeInput.value})
        } else {
            console.log(`Skipping dupe: ${event.currentTarget.elements.codeInput.value}`)
        }
        setSelectedCodes(selectedCodes)
    }
    function handleClearList(event: FormEvent) {
        event.preventDefault()
        console.log(`Clearing selection set (was length ${selectedCodes.length})`)
        setSelectedCodes(initialSelectedCodes)
    }

    return <>
        <Header />
        <h1>Move items around</h1>
        <p>Select a Location and Code Type. Then scan the items to be moved.
            Codes not found in the database will be created. When you're
            satisfied with the data, hit Save to update in bulk.</p>


        <form onSubmit={handleSelectCode}>
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
                       ref={codeInputRef}
                    // value={code}
                    // onChange={e => setCode(e.target.value)}
                />
            </div>
            <input type="submit" value="Add to list" />
        </form>
        <form>
            <div className="form-item">
                <label htmlFor="locationSelect">Location</label>
                <Select id="locationSelect" options={locations} />
            </div>
        </form>

        <h2>Selected so far</h2>
        <input type="button" value="Clear list" onClick={handleClearList} />
        <SelectedItemsList Codes={selectedCodes} />
    </>
}