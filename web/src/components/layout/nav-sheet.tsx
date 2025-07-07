import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Menu, Shield } from "lucide-react";
import { NavMenu } from "./nav-menu";

export const NavigationSheet = () => {
  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button variant="outline" size="icon">
          <Menu />
        </Button>
      </SheetTrigger>
      <SheetContent>
        <div className="flex items-center gap-2">
          <Shield className="h-8 w-8 text-blue-600" />

          <div>
            <h1 className="text-2xl font-bold text-gray-900">OSS Health</h1>
            <p className="text-sm text-gray-500">
              Dependency Security & Health Monitoring
            </p>
          </div>
        </div>
        <NavMenu orientation="vertical" className="mt-12" />
      </SheetContent>
    </Sheet>
  );
};
