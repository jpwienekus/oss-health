import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuContent,
  NavigationMenuTrigger,
} from '@/components/ui/navigation-menu'
import { type NavigationMenuProps } from '@radix-ui/react-navigation-menu'
import { Link } from 'react-router-dom'

export const NavMenu = (props: NavigationMenuProps) => {
  const isAdmin = true

  return (
    <NavigationMenu {...props}>
      <NavigationMenuList className="gap-6 space-x-0 data-[orientation=vertical]:flex-col data-[orientation=vertical]:items-start">
        <NavigationMenuItem>
          <NavigationMenuLink asChild>
            <Link key="/" to="/">
              Home
            </Link>
          </NavigationMenuLink>
        </NavigationMenuItem>
        {isAdmin && (
          <NavigationMenuItem>
            <NavigationMenuTrigger>Admin</NavigationMenuTrigger>
            <NavigationMenuContent>
              <ul className="grid w-[200px] gap-4">
                <li>
                  <NavigationMenuLink asChild>
                    <Link key="/admin/repositories" to="/admin/repositories">
                      Repositories
                    </Link>
                  </NavigationMenuLink>
                  <NavigationMenuLink asChild>
                    <Link key="/admin/dependencies" to="/admin/dependencies">
                      Dependencies
                    </Link>
                  </NavigationMenuLink>
                </li>
              </ul>
            </NavigationMenuContent>
          </NavigationMenuItem>
        )}
      </NavigationMenuList>
    </NavigationMenu>
  )
}
