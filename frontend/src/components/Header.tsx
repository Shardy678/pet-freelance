import { useContext } from 'react'
import { Link } from 'react-router-dom'
import { AuthContext } from '../context/AuthContext'

export default function Header() {
    const { token, logout } = useContext(AuthContext)

    return (
        <header className="bg-white shadow">
            <nav className="container mx-auto px-4 py-3 flex justify-between items-center">
                <div className="space-x-4">
                    <Link to="/services" className="hover:underline">Services</Link>
                    {token && (
                        <Link to="/bookings" className="hover:underline">
                            My Bookings
                        </Link>
                    )}
                </div>
                <div className="space-x-4">
                    <Link to="/profile" className='hover:underline'>Profile</Link>
                    {token ? (
                        <button
                            onClick={logout}
                            className="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600"
                        >
                            Logout
                        </button>
                    ) : (
                        <>
                            <Link to="/login" className="px-3 py-1 border rounded hover:bg-gray-100">
                                Login
                            </Link>
                            <Link to="/register" className="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700">
                                Register
                            </Link>
                        </>
                    )}
                </div>
            </nav>
        </header>
    )
}
