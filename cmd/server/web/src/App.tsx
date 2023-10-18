import {Button} from "./components/ui/button.tsx";
import {ThemeProvider} from "./components/theme-provider.tsx";
import DomainsPage from "./pages/domains/page.tsx";
import {Toaster} from "./components/ui/toaster.tsx";

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <Button variant="secondary">Hello</Button>
      <DomainsPage />
      <Toaster />
    </ThemeProvider>
  )
}

export default App
