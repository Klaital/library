export interface Settings {
    LIBRARY_API_BASE: string,
}

const dev: Settings = {
    LIBRARY_API_BASE: 'http://localhost:8080'
}

const prod: Settings = {
    LIBRARY_API_BASE: 'https://library.klaital.com'
}

export const Config = dev;
