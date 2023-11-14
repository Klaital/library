import {Config} from "../config/config";
import {LocationData} from "./location";
import {ItemData} from "./item";

export class LibraryApiClient {

    LookupCode(code: string, type: string): Promise<ItemData> {
        return fetch(Config.LIBRARY_API_BASE + '/api/code/' + type + '/' + code)
            .then((resp) => resp.json())
    }

    async FetchLocations(): Promise<LocationData[]> {
        const resp = await fetch(Config.LIBRARY_API_BASE + `/api/locations`);
        return await resp.json();
    }

    async CreateLocation(name: string, notes: string): Promise<LocationData> {
        const resp = await fetch(Config.LIBRARY_API_BASE + '/api/locations', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                Name: name,
                Notes: notes,
            })
        })
        return await resp.json()
    }
}
