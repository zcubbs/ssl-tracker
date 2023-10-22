import {ThemeProvider} from "./components/theme-provider.tsx";
import DomainsPage from "./pages/domains/page.tsx";
import {Toaster} from "./components/ui/toaster.tsx";
import Nav from "./components/nav/nav.tsx";
import {Route, Routes, useNavigate} from "react-router-dom";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import {useMemo, useState} from "react";
import AuthenticationPage from "@/pages/authentication/page.tsx";

const queryClient = new QueryClient();

function App() {
    const [isSignedIn, setIsSignedIn] = useState<boolean>(false)
    const navigate = useNavigate();
    // check if user is signed in
    // if not, redirect to login page
    // if yes, render the app


    return (
        <ThemeProvider defaultTheme="dark" storageKey=" vite-ui-theme">
            <QueryClientProvider client={queryClient}>
                <Routes>
                    <Route path="/" element={<Content/>}/>
                    <Route path="/domains" element={<Content/>}/>
                    <Route path="/login" element={<AuthenticationPage/>}/>
                </Routes>
                <Toaster/>
            </QueryClientProvider>
        </ThemeProvider>
    )
}

function Content() {
    return (
        <>
            <Nav/>
            <DomainsPage/>
        </>
    )
}

export default App
