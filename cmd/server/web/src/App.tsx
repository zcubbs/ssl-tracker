import {ThemeProvider} from "./components/theme-provider.tsx";
import DomainsPage from "./pages/domains/page.tsx";
import {Toaster} from "./components/ui/toaster.tsx";
import Nav from "./components/nav/nav.tsx";
import {BrowserRouter} from "react-router-dom";

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <BrowserRouter>
        <Nav/>
        <DomainsPage/>
      </BrowserRouter>
      <Toaster/>
    </ThemeProvider>
  )
}

export default App
