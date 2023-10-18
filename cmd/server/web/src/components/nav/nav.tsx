import {UserNav} from "../ui/user-nav.tsx";
import {Search} from "../search.tsx";
import SpaceSwitcher from "../space-switcher.tsx";
import {MainNav} from "./main-nav.tsx";

export default function Nav() {

  return (
    <div className="hidden flex-col md:flex">
      <div className="border-b">
        <div className="flex h-16 items-center px-4">
          <SpaceSwitcher />
          <MainNav className="mx-6" />
          <div className="ml-auto flex items-center space-x-4">
            <Search />
            <UserNav />
          </div>
        </div>
      </div>
    </div>
  )
}
