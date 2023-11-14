import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import reportWebVitals from './reportWebVitals';
import {createBrowserRouter, RouterProvider} from "react-router-dom";
import LocationsPage from "./pages/LocationsPage";
import {ItemsPage} from "./pages/ItemsPage";
import {AddVolumePage} from "./pages/AddVolumePage";

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

const router = createBrowserRouter([
    {
        path: "/",
        element: <LocationsPage />
    },
    {
        path: "/locations",
        element: <LocationsPage />
    },
    {
        path: "/items",
        element: <ItemsPage />
    },
    {
        path: "/add",
        element: <AddVolumePage />
    }
]);

root.render(
  <React.StrictMode>
      <RouterProvider router={router} />
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
