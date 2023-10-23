import {ThemeProvider} from "./components/theme-provider.tsx";
import DomainsPage from "./pages/domains/page.tsx";
import {Toaster} from "./components/ui/toaster.tsx";
import Nav from "./components/nav/nav.tsx";
import {Route, Routes} from "react-router-dom";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import LoginPage from "@/pages/login/page.tsx";
import RegisterPage from "@/pages/register/page.tsx";
import UsersPage from "@/pages/users/page.tsx";
import RequireAuth from "@/components/require-auth.tsx";
import NotFound from "@/pages/404/page.tsx";
import Unauthorized from "@/pages/unauthorized/page.tsx";

const queryClient = new QueryClient();

const ROLES = {
  'User': 2001,
  'Admin': 5150
}

function App() {
  // check if user is signed in
  // if not, redirect to login page
  // if yes, render the app


  return (
    <ThemeProvider defaultTheme="dark" storageKey=" vite-ui-theme">
      <QueryClientProvider client={queryClient}>
        <Routes>
          <Route path="/login" element={<LoginPage/>}/>
          <Route path="/register" element={<RegisterPage/>}/>
          <Route path="unauthorized" element={<Unauthorized />} />
          <Route element={<RequireAuth allowedRoles={[ROLES.User]}/>}>
            <Route path="/" element={<Content/>}/>
            <Route path="/domains" element={<Content/>}/>
          </Route>
          <Route element={<RequireAuth allowedRoles={[ROLES.Admin]}/>}>
            <Route path="/users" element={<UsersPage/>}/>
          </Route>
          <Route path="/logout" element={<h1>Logout</h1>}/>
          {/* catch all */}
          <Route path="*" element={<NotFound/>}/>
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
