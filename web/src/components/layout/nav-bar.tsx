import { Button } from "@/components/ui/button";
import { NavMenu } from "./nav-menu";
import { NavigationSheet } from "./nav-sheet";
import { Github, Shield, SunIcon } from "lucide-react";
import { useAuth } from "@/auth/AuthContext";
import { useEffect } from "react";
import { useGetUserQuery } from "@/generated/graphql";
import { toast } from "sonner";

export const Navbar = () => {
  const { jwt, loginWithGitHub } = useAuth()

  const { data, loading, error } = useGetUserQuery({
    skip: !jwt,
    notifyOnNetworkStatusChange: true
  })

  useEffect(() => {
    if (!error) {
      return
    }

    toast.error("Could not log in", {
      description: error.message,
    })

  }, [error])
  const isAuthenticated = !loading && jwt

  return (
    <div className="">
      <nav className="h-16 bg-background border-b">
        <div className="h-full flex items-center justify-between max-w-screen-xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-8">
            <div className="flex items-center gap-2">
              <Shield className="h-8 w-8 text-blue-600" />

              <div>
                <h1 className="text-2xl font-bold text-gray-900">OSS Health</h1>
                <p className="text-sm text-gray-500">
                  Dependency Security & Health Monitoring
                </p>
              </div>
            </div>

            <NavMenu className="hidden md:block" />
          </div>

          <div className="flex items-center gap-3">

            {isAuthenticated ? (
              <div className="flex items-center space-x-4">
                <span className="flex items-center text-sm text-gray-700">
                  <Github className="w-4 h-4 inline mr-2" />
                  {data?.username ?? ''}
                </span>
              </div>
            ) : (
              <Button onClick={loginWithGitHub} className="px-4 py-2">
                <Github className="w-4 h-4 inline mr-2" />
                Log in with GitHub
              </Button>
            )}
            <Button size="icon" variant="outline">
              <SunIcon />
            </Button>

            <div className="md:hidden">
              <NavigationSheet />
            </div>
          </div>
        </div>
      </nav>
    </div>
  );
};

