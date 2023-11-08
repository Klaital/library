import './header.css'

export default function Header() {
    return (
        <header className="site-header">
            <div className="wrapper site-header_wrapper">
            <div className="header-logo">
                <h1><a className="header-tile header-title" tabIndex={0} href="/">AF Library</a></h1>
            </div>
            <div className="header-menu">
                <a className="header-tile" href="/add">Add</a>
                <a className="header-tile" href="/locations">By Location</a>
                <a className="header-tile" href="/items">All Items</a>
            </div>
            <div className="header-actions">
                <ul className="user-actions">
                    <li className="user-actions-item">
                        <a className="header-tile" href="/locations">Search</a>
                    </li>
                </ul>
            </div>
            </div>
        </header>
    )
}