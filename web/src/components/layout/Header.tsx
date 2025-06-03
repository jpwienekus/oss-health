import { Link } from "react-router-dom"
import { Button } from "@/components/ui/button"
import { Menu, Moon, Sun } from "lucide-react"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/dropdown-menu"

type HeaderProps = {
  isDarkMode: boolean,
  toggleDarkMode: () => void
}

export const Header = ({ isDarkMode, toggleDarkMode }: HeaderProps) => {
  return (
    <header className="border-b bg-background sticky top-0 z-50">
      <div className="container mx-auto px-2 py-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center">
            <Link to="/" className="text-2xl">
              OSS-Health
            </Link>
          </div>

          <nav className="hidden md:flex space-x-6 items-center">
            <Link to="/" className="text-foreground hover:text-foreground/80 font-medium">
              Projects
            </Link>
            <div className="flex items-center space-x-2">
              <Button variant="ghost" size="icon" onClick={toggleDarkMode} aria-label="Search">
                {isDarkMode ? (
                  <Sun className="h-5 w-5" />
                ) : (
                  <Moon className="h-5 w-5" />
                )}
              </Button>
            </div>
          </nav>
          <nav className="flex items-center md:hidden space-x-2">
            <Button variant="ghost" size="icon" onClick={toggleDarkMode} aria-label="Search">
              {isDarkMode ? (
                <Sun className="h-5 w-5" />
              ) : (
                <Moon className="h-5 w-5" />
              )}
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" aria-label="Menu">
                  <Menu className="h-5 w-5" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent>

                <DropdownMenuItem>
                  <Link to="/" className="text-foreground hover:text-foreground/80 font-medium" >
                    Repositories
                  </Link>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </nav>
        </div>
      </div>

    </header>
  )
}
