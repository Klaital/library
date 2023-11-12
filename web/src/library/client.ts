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
}
