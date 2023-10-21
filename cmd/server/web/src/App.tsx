import {ThemeProvider} from "./components/theme-provider.tsx";
import DomainsPage from "./pages/domains/page.tsx";
import {Toaster} from "./components/ui/toaster.tsx";
import Nav from "./components/nav/nav.tsx";
import {BrowserRouter} from "react-router-dom";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";

const queryClient = new QueryClient();

function App() {
    return (
        <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
            <QueryClientProvider client={queryClient}>
                <BrowserRouter>
                    <Nav/>
                    <DomainsPage/>
                </BrowserRouter>
                <Toaster/>
            </QueryClientProvider>
        </ThemeProvider>
    )
}

export default App
